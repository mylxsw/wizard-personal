package utils

import (
	"github.com/mylxsw/asteria/log"
	"io"
)

type LogConverter struct{}

func NewLogConverter() io.Writer {
	return &LogConverter{}
}

func (pw LogConverter) Write(p []byte) (n int, err error) {
	log.Debugf("%s", string(p)[len("2020/05/13 00:10:09 "):])
	return len(p), nil
}
