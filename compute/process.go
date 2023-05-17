package compute

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hallllll/miraiweed/ready"
	"github.com/hallllll/miraiweed/scraping"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/semaphore"
)

func Procces(paths *ready.PATHs, urls *ready.URLs, P *ready.Put, bulk int) {
	// main process
	// とりあえず実行はするがエラーを返さないのはなんかアレだな
	records, errChan := ready.ReadCsv(paths.LoginInfo)
	var wg sync.WaitGroup
	sm := semaphore.NewWeighted(int64(bulk))
	for {
		select {
		case record, ok := <-records:
			if !ok {
				records = nil
			} else {
				wg.Add(1)
				go func(rec ready.LoginRecord) {
					if err := sm.Acquire(context.Background(), 1); err != nil {
						P.ErrLog.Fatalln(err)
					}
					defer sm.Release(1)

					defer P.StdLog.Printf("%s DONE.\n", rec.Name)
					defer wg.Done()

					// goahead
					opts := append(chromedp.DefaultExecAllocatorOptions[:],
						chromedp.Flag("headless", true),
					)

					allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
					allocCtx2, cancel := context.WithTimeout(allocCtx, 60*time.Second)

					// start
					ctx, cancel := chromedp.NewContext(
						allocCtx2,
						chromedp.WithLogf(P.StdLog.Printf),
					)
					defer cancel()

					tasks := chromedp.Tasks{
						scraping.GetScrapeCookies(urls.Base),
						scraping.LoginTasks(urls.Login, rec.Name, rec.ID, rec.PW, P),
						scraping.NavigateStudentsTasks(urls.StudentsSearch, rec.Name, P),
						scraping.DownloadStudentsTask(filepath.Join(paths.StudentFolder(), rec.Name), rec.Name, P),
						scraping.NavigateTeachersTasks(urls.TeacherSearch, rec.Name, P),
						scraping.DownloadTeachersTask(filepath.Join(paths.TeacherFolder(), rec.Name), rec.Name, P),
					}

					err := chromedp.Run(ctx, tasks)
					if err != nil {
						P.ErrLog.Println(err)
					}

				}(record)
			}

		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				P.ErrLog.Fatalf("Error reading CSV: %v\n", err)
			}
		}

		if records == nil && errChan == nil {
			break
		}
	}
	wg.Wait()
	time.Sleep(2 * time.Second)
}

func AllForOneSheet(paths *ready.PATHs) error {

	xlsx := excelize.NewFile()
	sw, err := xlsx.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	xlsxRowCount := 1
	header, err := excelize.CoordinatesToCellName(1, xlsxRowCount)
	if err != nil {
		return err
	}
	h := []interface{}{"学校名", "ID", "学年", "クラス", "出席番号", "氏名", "ふりがな", "ここには入力しないでください", "備考", "パスワード", "ここには入力しないでください", "ユーザーコネクトID", "まなびポケット共通ID", "G Suite", "SSO連携メールアドレス", "Azure AD", "SSO連携メールアドレス", "C4th共通ユーザーID", "エラー内容"}
	if err = sw.SetRow(header, h); err != nil {
		return err
	}
	fmt.Println(paths.StudentFolder())
	// // ???
	// filepath.WalkDir(paths.StudentFolder(), func(path string, d fs.DirEntry, err error) error {
	// 	fmt.Printf("とれてるよね？ %s\n", path)
	// 	return nil
	// })

	err = filepath.WalkDir(paths.StudentFolder(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "walkdir failed")
		}
		if d.IsDir() {
			return nil
		}

		dir, _ := filepath.Split(path)
		targetxlsx, err := excelize.OpenFile(path)
		if err != nil {
			log.Fatal(err)
			return err
		}
		rows, err := targetxlsx.Rows("子ども情報")
		for i := 0; rows.Next(); i++ {
			row, err := rows.Columns()
			if err != nil {
				return err
			}
			if i <= 2 {
				continue
			}
			xlsxRowCount++
			cell, err := excelize.CoordinatesToCellName(1, xlsxRowCount)
			if err != nil {
				xlsxRowCount--
				return err
			}
			if len(row) >= 1 && row[1] != "" {
				row[0] = filepath.Base(dir)
				val := make([]interface{}, len(row))
				for i, v := range row {
					val[i] = v
				}
				sw.SetRow(cell, val)
			} else {
				xlsxRowCount--
			}
		}
		return err
	})
	if err := sw.Flush(); err != nil {
		return err
	}
	if err := xlsx.SaveAs(filepath.Join(paths.StudentFolder(), "all.xlsx")); err != nil {
		return err
	}
	return nil
}
