package ready

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	dlFolderName       string = "data"
	studentsFolderName string = "miraiseed_students"
	teachersFolderName string = "miraiseed_teachers"
)

// save data directory architecture
type PATHs struct {
	Cd           string // ./miraiweed
	dlBase       string // ./miraiweed/dlBase
	dlStorage    string // ./miraiweed/dlBase/dlStorage (a.k.a per run data container)
	StudentsData string // ./miraiweed/dlBase/dlStorage/Students
	TeachersData string // ./miraiweed/dlBase/dlStorage/Teachers
	LoginInfo    string
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
	paths.dlStorage = filepath.Join(paths.Cd, paths.dlBase, strings.ReplaceAll(time.Now().Format("2006_01_02_150405.000"), ".", "_"))
	paths.StudentsData = filepath.Join(paths.dlStorage, studentsFolderName)
	paths.TeachersData = filepath.Join(paths.dlStorage, teachersFolderName)
	paths.LoginInfo = filepath.Join(paths.Cd, LoginCsvFileName)

	// create (if nothing) dl folder
	if err = os.MkdirAll(paths.dlStorage, 0755); err != nil {
		return nil, err
	}

	return paths, nil
}
