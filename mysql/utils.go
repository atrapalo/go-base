package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

const dateTimeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"

func TimeStamp() string {
	return time.Now().Format(dateTimeFormat)
}

func GenerateMd5(value string) string {
	hash := md5.Sum([]byte(value))

	return hex.EncodeToString(hash[:])
}
