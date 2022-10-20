package mysql_test

import (
	"fmt"
	"github.com/atrapalo/go-base/mysql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func Test_select(t *testing.T) {
	db, mock, err1 := sqlmock.New()
	assert.Nilf(t, err1, "an error '%s' was not expected when opening a mysql database connection", err1)

	queryBuilder := mysql.NewQueryBuilder(db)
	sql := queryBuilder.Select("id, title").From("posts", "").GetSQL()
	expectedSql := "SELECT id, title FROM posts"
	assert.Equalf(t, expectedSql, sql, "returned unexpected sql: got %s want %s", sql, expectedSql)

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	mock.ExpectQuery(sql).WillReturnRows(rows)
	_, err2 := queryBuilder.Query()
	assert.Nilf(t, err2, "error '%s' was not expected, while SELECT a row", err2)

	err3 := mock.ExpectationsWereMet()
	assert.Nilf(t, err3, "there were unfulfilled expectations: %s", err3)
}

func Test_inner_join(t *testing.T) {
	joinTestCommonParts(t, "INNER")
}

func Test_left_join(t *testing.T) {
	joinTestCommonParts(t, "LEFT")
}

func Test_right_join(t *testing.T) {
	joinTestCommonParts(t, "RIGHT")
}

func joinTestCommonParts(t *testing.T, joinFlag string) {
	db, mock, err1 := sqlmock.New()
	assert.Nilf(t, err1, "an error '%s' was not expected when opening a mysql database connection", err1)

	queryBuilder := mysql.NewQueryBuilder(db)
	sql, expectedSql := "", ""
	switch joinFlag {
	case "INNER":
		sql = queryBuilder.
			Select("p.id, p.title").
			From("posts", "p").
			SetFirstResult(0).
			SetMaxResults(3).
			InnerJoin("user", "u", "u.uid = p.uid").
			GetSQL()
		expectedSql = "SELECT p.id, p.title FROM posts p INNER JOIN user u ON u.uid = p.uid LIMIT 0,3"
		break
	case "LEFT":
		sql = queryBuilder.
			Select("p.id, p.title").
			From("posts", "p").
			SetFirstResult(0).
			SetMaxResults(3).
			LeftJoin("user", "u", "u.uid = p.uid").
			GetSQL()
		expectedSql = "SELECT p.id, p.title FROM posts p LEFT JOIN user u ON u.uid = p.uid LIMIT 0,3"
		break
	case "RIGHT":
		sql = queryBuilder.
			Select("p.id, p.title").
			From("posts", "p").
			SetFirstResult(0).
			SetMaxResults(3).
			RightJoin("user", "u", "u.uid = p.uid").
			GetSQL()
		expectedSql = "SELECT p.id, p.title FROM posts p RIGHT JOIN user u ON u.uid = p.uid LIMIT 0,3"
		break
	}

	assert.Equalf(t, expectedSql, sql, "returned unexpected sql: got %v want %v", sql, expectedSql)

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")
	mock.ExpectQuery(sql).WillReturnRows(rows)
	_, err2 := queryBuilder.Query()
	assert.Nilf(t, err2, "error '%s' was not expected, while SELECT a row", err2)

	err3 := mock.ExpectationsWereMet()
	assert.Nilf(t, err3, "there were unfulfilled expectations: %s", err3)
}

func Test_get_sql(t *testing.T) {
	db, _, err1 := sqlmock.New()
	assert.Nilf(t, err1, "an error '%s' was not expected when opening a mysql database connection", err1)

	queryBuilder := mysql.NewQueryBuilder(db)
	sql := queryBuilder.
		Select("uid, username, created, textVal, price, name").
		From("userinfo", "").
		Where("username = ? AND department = ?").
		SetParam("john").
		SetParam("DT").
		SetFirstResult(0).
		SetMaxResults(3).
		OrderBy("uid", "DESC").
		GetSQL()

	expectedSql := "SELECT uid, username, created, textVal, price, name FROM userinfo WHERE username = ? AND department = ? ORDER BY uid DESC LIMIT 0,3"
	assert.Equalf(t, expectedSql, sql, "returned unexpected sql: got %v want %v", sql, expectedSql)
	assert.Equal(t, []interface{}{"john", "DT"}, queryBuilder.GetParameters())
	assert.Equal(t, 3, queryBuilder.GetMaxResults())
	assert.Equal(t, 0, queryBuilder.GetFirstResult())
	fmt.Println(queryBuilder.GetParameters())

	queryBuilder2 := mysql.NewQueryBuilder(db)
	sql2 := queryBuilder2.
		Select("u.uid, u.username, p.address, count(*) as num").
		From("userinfo", "u").
		SetFirstResult(0).
		SetMaxResults(3).
		RightJoin("profile", "p", "u.uid = p.uid").
		Having("num > 1").
		GetSQL()

	expectedSql2 := "SELECT u.uid, u.username, p.address, count(*) as num FROM userinfo u RIGHT JOIN profile p ON u.uid = p.uid HAVING num > 1 LIMIT 0,3"
	assert.Equalf(t, expectedSql2, sql2, "returned unexpected sql: got %v want %v", sql2, expectedSql2)
}

func Test_delete(t *testing.T) {
	db, _, err1 := sqlmock.New()
	assert.Nilf(t, err1, "an error '%s' was not expected when opening a mysql database connection", err1)

	queryBuilder := mysql.NewQueryBuilder(db)
	sql := queryBuilder.Delete("userinfo").Where("uid = ?").SetParam(7).GetSQL()
	expectedSql := "DELETE FROM userinfo WHERE uid = ?"

	assert.Equalf(t, expectedSql, sql, "returned unexpected sql: got %v want %v", sql, expectedSql)
}

func Test_update(t *testing.T) {
	db, _, err1 := sqlmock.New()
	assert.Nilf(t, err1, "an error '%s' was not expected when opening a mysql database connection", err1)

	queryBuilder := mysql.NewQueryBuilder(db)
	sql := queryBuilder.
		Update("userinfo", "u").
		Set("u.username", "john2").
		Where("u.uid = ?").
		SetParam(1).
		GetSQL()

	expectedSql := "UPDATE userinfo u SET u.username = ? WHERE u.uid = ?"
	assert.Equalf(t, expectedSql, sql, "returned unexpected sql: got %v want %v", sql, expectedSql)
	assert.Equal(t, []interface{}{"john2", 1}, queryBuilder.GetParameters())
}
