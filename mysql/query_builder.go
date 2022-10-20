package mysql

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
)

// The query types.
const (
	Select = iota
	Delete
	Update
	Insert
)

// The query flags.
const (
	HasSort   = "hasSort"
	IsJoin    = "isJoin"
	IsDefault = "isDefault"
	Inner     = "INNER"
	Left      = "LEFT"
	Right     = "RIGHT"
)

type (
	// FromSqlParts records table and alias
	FromSqlParts struct {
		table, alias string
	}

	// OrderBySqlParts records sort and order
	OrderBySqlParts struct {
		sort, order string
	}

	// JoinSqlParts records joinType, joinTable, joinAlias, joinCondition
	JoinSqlParts struct {
		joinType, joinTable, joinAlias, joinCondition string
	}

	// ValuesSqlParts records key, val
	ValuesSqlParts struct {
		key string
		val interface{}
	}

	// SetSqlParts records key, val
	SetSqlParts struct {
		key string
		val interface{}
	}

	//QueryBuilder defined a SQL query builder.
	QueryBuilder struct {
		firstResult, maxResults, queryType                                                 int
		flag, hasSort, sql, sqlPartsSelect, sqlPartsWhere, sqlPartsGroupBy, sqlPartsHaving string
		database                                                                           *sql.DB
		State                                                                              *sql.Stmt
		params                                                                             []interface{}
		sqlPartsFrom                                                                       []FromSqlParts
		sqlPartsOrderBy                                                                    []OrderBySqlParts
		sqlPartsValues                                                                     []ValuesSqlParts
		sqlPartsSet                                                                        []SetSqlParts
		sqlPartsJoin                                                                       []JoinSqlParts
	}
)

// NewQueryBuilder returns a newly initialized QueryBuilder that implements QueryBuilder
func NewQueryBuilder(database *sql.DB) *QueryBuilder {
	return &QueryBuilder{
		firstResult:     0,
		maxResults:      -1,
		queryType:       Select,
		database:        database,
		params:          []interface{}{},
		flag:            IsDefault,
		sql:             "",
		sqlPartsSet:     make([]SetSqlParts, 0),
		sqlPartsValues:  make([]ValuesSqlParts, 0),
		sqlPartsFrom:    make([]FromSqlParts, 0),
		sqlPartsOrderBy: make([]OrderBySqlParts, 0),
		sqlPartsJoin:    make([]JoinSqlParts, 0),
	}
}

// GetParams returns queryBuilder params
func (queryBuilder *QueryBuilder) GetParams() []interface{} {
	return queryBuilder.params
}

// Select returns QueryBuilder that Specifies an item that has to be returned to the query result.
func (queryBuilder *QueryBuilder) Select(value string) *QueryBuilder {
	queryBuilder.queryType = Select
	queryBuilder.sqlPartsSelect = value

	return queryBuilder
}

// From returns QueryBuilder that creates and adds a query root corresponding to the table identified by the
// given alias, forming a cartesian product with any existing query roots.
func (queryBuilder *QueryBuilder) From(table string, alias string) *QueryBuilder {
	queryBuilder.setFromWrap(table, alias)

	return queryBuilder
}

// Update returns QueryBuilder that turns the query being built into a bulk update query that ranges over
//a certain table
func (queryBuilder *QueryBuilder) Update(table string, alias string) *QueryBuilder {
	queryBuilder.queryType = Update
	queryBuilder.setFromWrap(table, alias)

	return queryBuilder
}

// Set returns QueryBuilder that sets a new value for a column in a bulk update query.
func (queryBuilder *QueryBuilder) Set(key string, val interface{}) *QueryBuilder {
	queryBuilder.sqlPartsSet = append(queryBuilder.sqlPartsSet, SetSqlParts{key: key, val: val})

	return queryBuilder
}

// Value returns QueryBuilder that sets a new value for a column in a bulk insert query.
func (queryBuilder *QueryBuilder) Value(key string, val interface{}) *QueryBuilder {
	queryBuilder.sqlPartsValues = append(queryBuilder.sqlPartsValues, ValuesSqlParts{key: key, val: val})

	return queryBuilder
}

// OrderBy returns QueryBuilder that specifies an ordering for the query results.
func (queryBuilder *QueryBuilder) OrderBy(sort string, order string) *QueryBuilder {

	queryBuilder.hasSort = HasSort
	if order == "" {
		order = "ASC"
	}

	queryBuilder.sqlPartsOrderBy = append(queryBuilder.sqlPartsOrderBy, OrderBySqlParts{sort, order})

	return queryBuilder
}

// GroupBy returns QueryBuilder that specifies a grouping over the results of the query.
func (queryBuilder *QueryBuilder) GroupBy(groupBy string) *QueryBuilder {
	if groupBy == "" {
		return queryBuilder
	}

	queryBuilder.sqlPartsGroupBy = groupBy

	return queryBuilder
}

// Having returns QueryBuilder that specifies a restriction over the groups of the query.
func (queryBuilder *QueryBuilder) Having(having string) *QueryBuilder {
	queryBuilder.sqlPartsHaving = having

	return queryBuilder
}

// SetFirstResult returns QueryBuilder that sets the position of the first result to retrieve.
func (queryBuilder *QueryBuilder) SetFirstResult(firstResult int) *QueryBuilder {
	queryBuilder.firstResult = firstResult

	return queryBuilder
}

// Where returns QueryBuilder that specifies one or more restrictions to the query result.
func (queryBuilder *QueryBuilder) Where(condition string) *QueryBuilder {
	queryBuilder.sqlPartsWhere = condition

	return queryBuilder
}

// Join returns QueryBuilder that creates and adds a join to the query.
func (queryBuilder *QueryBuilder) Join(join string, alias string, condition string) *QueryBuilder {
	return queryBuilder.InnerJoin(join, alias, condition)
}

// InnerJoin returns QueryBuilder that creates and adds a join to the query.
func (queryBuilder *QueryBuilder) InnerJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.flag = IsJoin
	queryBuilder.sqlPartsJoin = append(queryBuilder.sqlPartsJoin, JoinSqlParts{joinType: Inner, joinTable: join, joinAlias: alias, joinCondition: condition})

	return queryBuilder
}

// LeftJoin returns QueryBuilder that creates and adds a left join to the query.
func (queryBuilder *QueryBuilder) LeftJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.flag = IsJoin
	queryBuilder.sqlPartsJoin = append(queryBuilder.sqlPartsJoin, JoinSqlParts{joinType: Left, joinTable: join, joinAlias: alias, joinCondition: condition})

	return queryBuilder
}

// RightJoin returns QueryBuilder that creates and adds a right join to the query.
func (queryBuilder *QueryBuilder) RightJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.flag = IsJoin
	queryBuilder.sqlPartsJoin = append(queryBuilder.sqlPartsJoin, JoinSqlParts{joinType: Right, joinTable: join, joinAlias: alias, joinCondition: condition})

	return queryBuilder
}

// GetFirstResult gets the position of the first result the query object was set to retrieve.
func (queryBuilder *QueryBuilder) GetFirstResult() int {
	return queryBuilder.firstResult
}

// SetMaxResults sets the maximum number of results to retrieve.
func (queryBuilder *QueryBuilder) SetMaxResults(maxResults int) *QueryBuilder {
	queryBuilder.maxResults = maxResults

	return queryBuilder
}

// GetMaxResults gets the maximum number of results the query object was set to retrieve
func (queryBuilder *QueryBuilder) GetMaxResults() int {
	return queryBuilder.maxResults
}

// SetParam sets a query parameter for the query being constructed.
func (queryBuilder *QueryBuilder) SetParam(param interface{}) *QueryBuilder {
	queryBuilder.params = append(queryBuilder.params, param)

	return queryBuilder
}

// GetParameters gets all defined query parameters for the query being constructed indexed by parameter index or name.
func (queryBuilder *QueryBuilder) GetParameters() []interface{} {
	return queryBuilder.params
}

// GetSQL gets the complete SQL string formed by the current specifications of this QueryBuilder.
func (queryBuilder *QueryBuilder) GetSQL() string {
	sqlString := queryBuilder.sql

	if sqlString != "" {
		return sqlString
	}

	queryType := queryBuilder.queryType

	switch queryType {
	case Insert:
		sqlString = queryBuilder.getSQLForInsert()
	case Delete:
		sqlString = queryBuilder.getSQLForDelete()
	case Update:
		sqlString = queryBuilder.getSQLForUpdate()
	case Select:
		sqlString = queryBuilder.getSQLForSelect()
	default:
		sqlString = queryBuilder.getSQLForSelect()
	}
	queryBuilder.sql = sqlString

	return cleanUpMessySQL(sqlString)
}

// getSQLForUpdate returns an update string in SQL.
func (queryBuilder *QueryBuilder) getSQLForUpdate() string {
	sqlString := "UPDATE "

	table := ""
	for _, v := range queryBuilder.sqlPartsFrom {
		table = v.table + " " + v.alias
	}

	sortedKeys := make([]string, 0)

	paramsTemp := make([]interface{}, 0)

	sqlString += table + " SET "

	for _, v := range queryBuilder.sqlPartsSet {
		sortedKeys = append(sortedKeys, v.key)

		sqlString += v.key + " = ? ,"

		paramsTemp = append(paramsTemp, v.val)
	}

	for _, v := range queryBuilder.params {
		paramsTemp = append(paramsTemp, v)
	}

	queryBuilder.params = paramsTemp

	sqlString = sqlString[:len(sqlString)-1]

	if whereStr := queryBuilder.sqlPartsWhere; whereStr != "" {
		sqlString += " WHERE " + whereStr
	}

	return sqlString

}

// getSQLForJoins returns a join string in SQL.
func (queryBuilder *QueryBuilder) getSQLForJoins() string {
	sqlString := ""

	if queryBuilder.flag != IsJoin {
		return ""
	}

	for _, v := range queryBuilder.sqlPartsJoin {
		joinType := v.joinType
		joinTable := v.joinTable
		joinAlias := v.joinAlias
		joinCondition := v.joinCondition

		sqlString += " " + joinType + " JOIN " + joinTable + " " + joinAlias + " ON " + joinCondition

		return sqlString
	}

	return sqlString
}

// getFromClauses returns table or join sql string
func (queryBuilder *QueryBuilder) getFromClauses() string {
	tableSql := ""

	for _, v := range queryBuilder.sqlPartsFrom {
		tableSql = v.table + " " + v.alias
		return tableSql + queryBuilder.getSQLForJoins()
	}

	return tableSql
}

// getSQLForSelect returns a select string in SQL.
func (queryBuilder *QueryBuilder) getSQLForSelect() string {
	sqlString := "SELECT "

	if selectStr := queryBuilder.sqlPartsSelect; selectStr != "" {
		sqlString += selectStr
	}

	sqlString += " FROM " + queryBuilder.getFromClauses()

	if whereStr := queryBuilder.sqlPartsWhere; whereStr != "" {
		sqlString += " WHERE " + whereStr
	}

	if groupByStr := queryBuilder.sqlPartsGroupBy; groupByStr != "" {
		sqlString += " GROUP BY " + groupByStr
	}

	if havingStr := queryBuilder.sqlPartsHaving; havingStr != "" {
		sqlString += " HAVING " + havingStr
	}

	if queryBuilder.hasSort == HasSort {
		sqlString += " ORDER BY "

		for _, v := range queryBuilder.sqlPartsOrderBy {
			sqlString += v.sort + " " + v.order + ", "
		}
		sqlString = sqlString[:len(sqlString)-2]
	}

	if queryBuilder.isLimitQuery() {
		sqlString += " LIMIT " + strconv.Itoa(queryBuilder.firstResult) + "," + strconv.Itoa(queryBuilder.maxResults)
	}

	return sqlString
}

// getSQLForDelete returns an delete string in SQL.
func (queryBuilder *QueryBuilder) getSQLForDelete() string {
	sqlString := "DELETE "

	for _, v := range queryBuilder.sqlPartsFrom {
		sqlString += " FROM " + v.table
		if whereStr := queryBuilder.sqlPartsWhere; whereStr != "" {
			sqlString += " WHERE " + whereStr
		}

		return sqlString
	}

	return sqlString
}

// getSQLForInsert returns an insert string in SQL.
func (queryBuilder *QueryBuilder) getSQLForInsert() string {
	sqlString := "INSERT INTO" + " "

	for _, v := range queryBuilder.sqlPartsFrom {
		tableSql := v.table
		sqlString += tableSql + " ("

		sortedKeys := make([]string, 0)

		values := ""

		params := make([]interface{}, 0)

		for _, v := range queryBuilder.sqlPartsValues {
			sortedKeys = append(sortedKeys, v.key)

			sqlString += v.key + ", "
			values += "?, "

			params = append(params, v.val)

		}

		queryBuilder.params = params

		sqlString = sqlString[:len(sqlString)-2]
		values = values[:len(values)-2]
		sqlString += ") VALUES (" + values + ")"

		return sqlString
	}

	return sqlString
}

// isLimitQuery returns is a limited Query
func (queryBuilder *QueryBuilder) isLimitQuery() bool {
	if queryBuilder.maxResults == -1 {
		return false
	}
	if queryBuilder.maxResults > 0 && queryBuilder.firstResult >= 0 {
		return true
	}

	return false
}

// ExecuteQuery executes a query that returns rows
func (queryBuilder *QueryBuilder) ExecuteQuery(query string) (*sql.Rows, error) {
	if queryBuilder.params != nil {
		rows, err := queryBuilder.database.Query(query, queryBuilder.params...)

		return rows, err
	}

	rows, err := queryBuilder.database.Query(query, nil)

	return rows, err
}

// ExecuteQueryAndGetRowsMap executes a query that returns rows map
func (queryBuilder *QueryBuilder) ExecuteQueryAndGetRowsMap(query string) (map[int]map[string]Field, error) {
	if queryBuilder.params != nil {
		rows, err := queryBuilder.database.Query(query, queryBuilder.params...)
		if err != nil {
			return nil, err
		}

		return getRowsMap(rows), err
	}

	rows, err := queryBuilder.database.Query(query, nil)
	if err != nil {
		return nil, err
	}

	return getRowsMap(rows), err
}

// getRowsMap returns rows map
func getRowsMap(rows *sql.Rows) map[int]map[string]Field {
	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	values := make([]interface{}, len(columns))
	columnPointers := make([]interface{}, len(columns))

	resultId := 0
	result := map[int]map[string]Field{}

	for rows.Next() {
		for i := range columns {
			columnPointers[i] = &values[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			panic(err)
		}

		record := map[string]Field{}

		for i, col := range columns {
			val := values[i]
			record[col] = NewField(val, columnTypes[i].DatabaseTypeName())
		}

		result[resultId] = record
		resultId++
	}

	return result
}

// Query executes a query that returns rows
func (queryBuilder *QueryBuilder) Query() (*sql.Rows, error) {
	if queryBuilder.queryType == Select {
		return queryBuilder.ExecuteQuery(queryBuilder.GetSQL())
	}

	return nil, nil
}

// QueryAssoc executes a query that returns rows map
func (queryBuilder *QueryBuilder) QueryAssoc() (map[int]map[string]Field, error) {
	if queryBuilder.queryType == Select {
		return queryBuilder.ExecuteQueryAndGetRowsMap(queryBuilder.GetSQL())
	}

	return nil, nil
}

// prepareAndExecute creates a prepared statement for later queries or executions.
func (queryBuilder *QueryBuilder) prepareAndExecute() sql.Result {
	stmt, err := queryBuilder.database.Prepare(queryBuilder.GetSQL())
	if err != nil {
		panic(err)
	}
	queryBuilder.State = stmt
	res, err := stmt.Exec(queryBuilder.params...)
	if err != nil {
		panic(err)
	}

	return res
}

// PrepareAndExecute creates a prepared statement for later queries or executions.
func (queryBuilder *QueryBuilder) PrepareAndExecute() (int64, error) {
	if queryBuilder.queryType == Insert {
		res := queryBuilder.prepareAndExecute()
		return res.LastInsertId()
	}

	if queryBuilder.queryType == Delete {
		res := queryBuilder.prepareAndExecute()
		return res.RowsAffected()
	}

	if queryBuilder.queryType == Update {
		res := queryBuilder.prepareAndExecute()
		return res.RowsAffected()
	}

	return -1, nil
}

// Insert turns the query being built into an insert query that inserts into
func (queryBuilder *QueryBuilder) Insert(table string) *QueryBuilder {
	queryBuilder.queryType = Insert
	queryBuilder.setFromWrap(table, "")

	return queryBuilder
}

// Delete turns the query being built into a bulk delete query that ranges over
func (queryBuilder *QueryBuilder) Delete(table string) *QueryBuilder {
	queryBuilder.queryType = Delete
	queryBuilder.setFromWrap(table, "")

	return queryBuilder
}

// setFromWrap wraps sqlParts `from`
func (queryBuilder *QueryBuilder) setFromWrap(table string, alias string) {
	queryBuilder.sqlPartsFrom = append(queryBuilder.sqlPartsFrom, FromSqlParts{table, alias})
}

func cleanUpMessySQL(sql string) string {
	space := regexp.MustCompile(`\s+`)
	sql = space.ReplaceAllString(sql, " ")

	return strings.Trim(sql, " ")
}
