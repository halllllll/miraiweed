package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/halllllll/miraiweed/compute"
	"github.com/halllllll/miraiweed/ready"
)

var (
	err                      error
	urls                     *ready.URLs
	paths                    *ready.PATHs
	P                        *ready.Put
	default_miraiseedx       int      = 7
	default_concuarrency_num int      = 10
	studentHeader            []string = []string{"学校名", "ID", "学年", "クラス", "出席番号", "氏名", "ふりがな", "ここには入力しないでください", "備考", "パスワード", "ここには入力しないでください", "", "", "ユーザーコネクトID", "まなびポケット共通ID", "G Suite SSO連携メールアドレス", "Azure AD SSO連携メールアドレス", "C4th共通ユーザーID", "エラー内容"}
	studentSheetName         string   = "子ども情報"
	teacherHeader            []string = []string{"学校名", "ID", "氏名", "ふりがな", "所属学年", "担任クラス", "担当教科", "授業を受け持つクラス", "備考", "パスワード", "ユーザーコネクトID", "ユーザーID（任意設定）", "まなびポケット共通ID", "先生カルテ閲覧権限 (Evit先生アンケート)", "G Suite SSO連携メールアドレス", "Azure AD SSO連携メールアドレス", "C4th共通ユーザーID", "エラー内容"}
	teacherSheetName         string   = "先生情報"

	educationBoardKey string = "giga"
)

type XlsxInfo struct {
	header    []string
	sheetName string
	path      string
}

type ScrapeConfig struct {
	IsEduboardMode bool
	EduBoard       EduBoardInfo
	Bulk           int
	MiraiseedX     string
}

type EduBoardInfo struct {
	Id string
	Pw string
}

func hello() ScrapeConfig {
	// This Program Requires some neccesary files.
	// And, there are some steps to assumed, prerequisite settings

	// 0.0 - Prepare struct will return
	helloResp := &ScrapeConfig{}

	// 0.1 -  For stdout and logs
	P = ready.NewPut()
	P.LoggingSetting("miraiweed.log")

	// 0.2 - prepare almost all of paths and naming roles for this project
	paths, err = ready.NewPATHs()
	if err != nil {
		P.ErrLog.Fatal(err)
	}

	// 1.0 - Education Board Mode
	var isEduBoard bool
	// INDEV
	/*

		yesIamEdu, err := ready.PromptAndRead(fmt.Sprintf("Type '%s' if you want to use Education Board Mode (required ID/PW): ", educationBoardKey))
		if err != nil {
			log.Fatal(err)
		}
		if strings.ToLower(yesIamEdu) == strings.ToLower(educationBoardKey) {
			isEduBoard = true
		}
	*/

	// 2.0 -  To login, ofcourse there are User ID/PW or going along with such info. Or, when you have selected EduBoard mode, required EduBoard ID/PW

	if isEduBoard {
		for {
			eduBoardId, err := ready.PromptAndRead(fmt.Sprint("Education Board Account ID: "))
			if err != nil {
				log.Fatal(err)
			}
			if eduBoardId == "" {
				fmt.Println("Can't Accept Empty Line")
				continue
			}
			fmt.Println("thank you 😘")
			helloResp.EduBoard.Id = eduBoardId
			break
		}
		for {
			eduBoardPw, err := ready.PromptAndRead(fmt.Sprint("Education Board Account Password: "))
			if err != nil {
				log.Fatal(err)
			}
			if eduBoardPw == "" {
				fmt.Println("Can't Accept Empty Line")
				continue
			}
			fmt.Println("thank you 🥰")
			helloResp.EduBoard.Pw = eduBoardPw
			break
		}
		helloResp.IsEduboardMode = true
		P.InfoLog.Println("-- EduBoard Mode --")
	} else {
		helloResp.IsEduboardMode = false
		P.InfoLog.Println("-- Normal Mode --")
		if _, err := os.Stat(paths.LoginInfo); os.IsNotExist(err) {
			fmt.Println("The CSV file does not exist. Creating a new template...")
			if err = ready.CreateCsvTemplate(paths.LoginInfo); err != nil {
				P.StdLog.Fatalf("Failed to create the CSV template: %s", err)
			}
			fmt.Println("Template created. Please fill it with data and run the program again.")
			P.InfoLog.Println("Create CSV Template.")
			os.Exit(1) // over
		}
	}

	// 3.0 -  miraiseed instance number prompt. cuz miraiseed serving some url-s for bunch of local goverments.
	for {
		_miraiseedX, err := ready.PromptAndRead(fmt.Sprintf("enter miraiseed[X] (default=%d): ", default_miraiseedx))
		if err != nil {
			log.Fatal(err)
		}
		if _miraiseedX == "" {
			_miraiseedX = strconv.Itoa(default_miraiseedx)
		}
		x, err := strconv.Atoi(_miraiseedX)
		if err != nil {
			log.Fatal(err)
		}
		if x <= 0 || 10 <= x {
			fmt.Println("we can only accept number between 1 from 9.")
		} else {
			helloResp.MiraiseedX = _miraiseedX
			break
		}
	}

	// 4.0 - confirm concurrency limit(for semaphore). in general GIGA School Management Organizations are NOT experts in their field, luck of IT knowledge and development skills. therefore they are forced to use cheap and low-spec business PCs. use of `answer` number as the default value is a consideration for such an environment and is not intended to be otherwize.
	for {
		concurrency_num, err := ready.PromptAndRead(fmt.Sprintf("Concarrency Limit (default=%d):", default_concuarrency_num))
		if err != nil {
			log.Fatal(err)
		}
		if concurrency_num == "" {
			concurrency_num = strconv.Itoa(default_concuarrency_num)
		}
		_bulk, err := strconv.Atoi(concurrency_num)
		if err != nil {
			log.Fatal(err)
		}
		if _bulk <= 0 {
			fmt.Println("we can only accept number upper 1.")
		} else {
			helloResp.Bulk = _bulk
			break
		}
	}
	P.InfoLog.Printf("miraiseedX = %s\n", helloResp.MiraiseedX)
	P.InfoLog.Printf("bulk = %d\n", helloResp.Bulk)
	// OK, so, some settings under here

	// main urls used in this project.
	urls = ready.NewUrls()
	urls.PrepareUrl(helloResp.MiraiseedX)

	return *helloResp
}

func main() {

	resp := hello()
	if resp.IsEduboardMode {
		p := &ready.LoginRecord{
			ID: resp.EduBoard.Id,
			PW: resp.EduBoard.Pw,
		}
		compute.EduBoardProcess(paths, urls, *p, P)
	} else {
		compute.Procces(paths, urls, P, resp.Bulk)
	}

	var wg sync.WaitGroup
	info := []XlsxInfo{
		{sheetName: studentSheetName, header: studentHeader, path: paths.StudentFolder()},
		{sheetName: teacherSheetName, header: teacherHeader, path: paths.TeacherFolder()},
	}
	for _, v := range info {
		wg.Add(1)
		go func(_v XlsxInfo) {
			if err = compute.AllForOneSheet(_v.path, _v.header, _v.sheetName, P); err != nil {
				P.ErrLog.Println(err)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	ready.PromptAndRead("Byebye ﾉｼ")
}
