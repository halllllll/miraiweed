package ready

import (
	"io"
	"log"
	"os"
)

type Put struct {
	PutWriter io.Writer
	StdLog    *log.Logger
	ErrLog    *log.Logger
	InfoLog   *log.Logger
	Panic     *log.Logger
}

func NewPut() *Put {
	return &Put{}
}

func (p *Put) LoggingSetting(fileName string) *io.Writer {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	p.PutWriter = io.MultiWriter(f, os.Stderr)
	p.StdLog = log.New(f, "[Std] ", log.Ldate|log.Ltime)
	p.ErrLog = log.New(f, "[Error] ", log.Ldate|log.Ltime)
	p.InfoLog = log.New(f, "[Info] ", log.Ldate|log.Ltime)
	p.Panic = log.New(f, "[PANIC] ", log.Ldate|log.Ltime)
	p.StdLog.SetOutput(p.PutWriter)
	p.ErrLog.SetOutput(p.PutWriter)
	p.InfoLog.SetOutput(p.PutWriter)
	p.Panic.SetOutput(p.PutWriter)
	return &p.PutWriter
}
