package ready

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	dlFolderName string = "data"
)

type PATHs struct {
	Cd        string
	dlBase    string
	DLStorage string
	LoginInfo string
}

func NewPATHs() (*PATHs, error) {
	paths := new(PATHs)
	// Cd
	pathsCd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	paths.Cd = pathsCd
	// dlBase
	paths.dlBase = dlFolderName

	// DLStorage(Name Rule: yyyy_MM_dd_hhmmssSSS)
	paths.DLStorage = filepath.Join(paths.Cd, dlFolderName, strings.ReplaceAll(time.Now().Format("2006_01_02_150405.000"), ".", "_"))
	paths.LoginInfo = filepath.Join(paths.Cd, LoginCsvFileName)

	return paths, nil
}
