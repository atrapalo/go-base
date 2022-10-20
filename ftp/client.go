package ftp

import (
	"github.com/webguerilla/ftps"
)

type ClientFtpGuerrilla struct {
	*ftps.FTPS
}

func New(allowInsecure bool, debug bool) ClientFtpGuerrilla {
	ftp := new(ftps.FTPS)
	ftp.TLSConfig.InsecureSkipVerify = allowInsecure
	ftp.Debug = debug

	return ClientFtpGuerrilla{
		FTPS: ftp,
	}
}
