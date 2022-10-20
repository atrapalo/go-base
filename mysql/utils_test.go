package mysql_test

import (
	"github.com/atrapalo/go-base/mysql"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_timestamp_string_format(t *testing.T) {
	assert.Regexp(t, regexp.MustCompile(`\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}`), mysql.TimeStamp())
}

func Test_md5_generation(t *testing.T) {
	hash := mysql.GenerateMd5("whatever")

	assert.Regexp(t, regexp.MustCompile(`[0-9a-f]{32}`), hash)
	assert.Len(t, hash, 32)
}
