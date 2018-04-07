package logger

import (
	"io"
	"log"
)

var (
	Err *log.Logger
	Out *log.Logger
)

func Init(err, std io.Writer) {
	Err = log.New(err, "", log.Llongfile)
	Out = log.New(std, "", 0)
}
