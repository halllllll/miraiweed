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
	cd           string // ./miraiweed
	dlBase       string // ./miraiweed/dlBase
	dlStorage    string // ./miraiweed/dlBase/dlStorage (a.k.a per run data container)
	studentsData string // ./miraiweed/dlBase/dlStorage/Students
	teachersData string // ./miraiweed/dlBase/dlStorage/Teachers
	LoginInfo    string
}

func NewPATHs() (*PATHs, error) {
	paths := new(PATHs)
	// Cd
	pathsCd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	paths.cd = pathsCd
	// dlBase
	paths.dlBase = dlFolderName

	// DLStorage(Name Rule: yyyy_MM_dd_hhmmssSSS)
	paths.dlStorage = strings.ReplaceAll(time.Now().Format("2006_01_02_150405.000"), ".", "_")
	paths.studentsData = studentsFolderName
	paths.teachersData = teachersFolderName
	paths.LoginInfo = LoginCsvFileName

	// create (if nothing) dl folder
	if err = os.MkdirAll(paths.Storage(), 0755); err != nil {
		return nil, err
	}

	return paths, nil
}

func (paths *PATHs) Base() string {
	return filepath.Join(paths.cd, paths.dlBase)
}

func (paths *PATHs) Storage() string {
	return filepath.Join(paths.cd, paths.dlBase, paths.dlStorage)
}

func (paths *PATHs) StudentFolder() string {
	return filepath.Join(paths.cd, paths.dlBase, paths.dlStorage, paths.studentsData)
}

func (paths *PATHs) TeacherFolder() string {
	return filepath.Join(paths.cd, paths.dlBase, paths.dlStorage, paths.teachersData)
}

// overload path

func (paths *PATHs) ChangeBase(newName string) (string, error) {
	if _, err := os.Stat(filepath.Join(paths.cd, newName)); os.IsNotExist(err) {
		return "", err
	}
	paths.dlBase = newName
	return paths.Base(), nil
}

func (paths *PATHs) ChangeStorage(newName string) (string, error) {
	if _, err := os.Stat(filepath.Join(paths.cd, paths.dlBase, newName)); os.IsNotExist(err) {
		return "", err
	}
	paths.dlStorage = newName
	return paths.Storage(), nil
}

func (paths *PATHs) ChangeStudentsFolder(newName string) (string, error) {
	if _, err := os.Stat(filepath.Join(paths.cd, paths.dlBase, newName)); os.IsNotExist(err) {
		return "", err
	}
	paths.studentsData = newName
	return paths.StudentFolder(), nil
}

func (paths *PATHs) ChangeTeachersFolder(newName string) (string, error) {
	if _, err := os.Stat(filepath.Join(paths.cd, paths.dlBase, newName)); os.IsNotExist(err) {
		return "", err
	}
	paths.teachersData = newName
	return paths.StudentFolder(), nil
}
