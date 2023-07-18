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
	bulk                     int
	urls                     *ready.URLs
	paths                    *ready.PATHs
	P                        *ready.Put
	default_miraiseedx       int      = 7
	default_concuarrency_num int      = 10
	studentHeader            []string = []string{"学校名", "ID", "学年", "クラス", "出席番号", "氏名", "ふりがな", "ここには入力しないでください", "備考", "パスワード", "ここには入力しないでください", "", "", "ユーザーコネクトID", "まなびポケット共通ID", "G Suite SSO連携メールアドレス", "Azure AD SSO連携メールアドレス", "C4th共通ユーザーID", "エラー内容"}
	studentSheetName         string   = "子ども情報"
	teacherHeader            []string = []string{"学校名", "ID", "氏名", "ふりがな", "所属学年", "担任クラス", "担当教科", "授業を受け持つクラス", "備考", "パスワード", "ユーザーコネクトID", "ユーザーID（任意設定）", "まなびポケット共通ID", "先生カルテ閲覧権限 (Evit先生アンケート)", "G Suite SSO連携メールアドレス", "Azure AD SSO連携メールアドレス", "C4th共通ユーザーID", "エラー内容"}
	teacherSheetName         string   = "先生情報"
)

type XlsxInfo struct {
	header    []string
	sheetName string
	path      string
}

func hello() {
	// This Program Requires some neccesary files.
	// And, there are some steps to assumed, prerequisite settings

	// 0.1 -  For stdout and logs
	P = ready.NewPut()
	P.LoggingSetting("miraiweed.log")

	// 0.2 - prepare almost all of paths and naming roles for this project
	paths, err = ready.NewPATHs()
	if err != nil {
		P.ErrLog.Fatal(err)
	}

	// 1.0 -  To login, ofcourse there are User ID/PW or going along with such info.
	if _, err := os.Stat(paths.LoginInfo); os.IsNotExist(err) {
		fmt.Println("The CSV file does not exist. Creating a new template...")
		if err = ready.CreateCsvTemplate(paths.LoginInfo); err != nil {
			P.StdLog.Fatalf("Failed to create the CSV template: %s", err)
		}
		fmt.Println("Template created. Please fill it with data and run the program again.")
		P.InfoLog.Println("Create CSV Template.")
		os.Exit(1) // over
	}

	// 2.0 -  miraiseed instance number prompt. cuz miraiseed serving some url-s for bunch of local goverments.
	var miraiseedX string
	for {
		miraiseedX, err = ready.PromptAndRead(fmt.Sprintf("enter miraiseed[X] (default=%d): ", default_miraiseedx))
		if err != nil {
			log.Fatal(err)
		}
		if miraiseedX == "" {
			miraiseedX = strconv.Itoa(default_miraiseedx)
		}
		x, err := strconv.Atoi(miraiseedX)
		if err != nil {
			log.Fatal(err)
		}
		if x <= 0 || 10 <= x {
			fmt.Println("we can only accept number between 1 from 9.")
		} else {
			break
		}
	}

	// 3.0 - confirm concurrency limit(for semaphore). in general GIGA School Management Organizations are NOT experts in their field, luck of IT knowledge and development skills. therefore they are forced to use cheap and low-spec business PCs. use of `answer` number as the default value is a consideration for such an environment and is not intended to be otherwize.
	for {
		concurrency_num, err := ready.PromptAndRead(fmt.Sprintf("Concarrency Limit (default=%d):", default_concuarrency_num))
		if err != nil {
			log.Fatal(err)
		}
		if concurrency_num == "" {
			concurrency_num = strconv.Itoa(default_concuarrency_num)
		}
		bulk, err = strconv.Atoi(concurrency_num)
		if err != nil {
			log.Fatal(err)
		}
		if bulk <= 0 {
			fmt.Println("we can only accept number upper 1.")
		} else {
			break
		}
	}
	P.InfoLog.Printf("miraiseedX = %s\n", miraiseedX)
	P.InfoLog.Printf("bulk = %d\n", bulk)
	// OK, so, some settings under here

	// main urls used in this project.
	urls = ready.NewUrls()
	urls.PrepareUrl(miraiseedX)

}

func main() {

	hello()

	compute.Procces(paths, urls, P, bulk)

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
