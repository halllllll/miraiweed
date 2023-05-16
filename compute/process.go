package compute

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hallllll/miraiweed/ready"
	"github.com/hallllll/miraiweed/scrape"
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
						scrape.GetScrapeCookies(urls.Base),
						scrape.LoginTasks(urls.Login, rec.Name, rec.ID, rec.PW, P),
						scrape.NavigateStudentsTasks(urls.StudentsSearch, rec.Name, P),
						scrape.DownloadStudentsTask(filepath.Join(paths.StudentsData, rec.Name), rec.Name, P),
						scrape.NavigateTeachersTasks(urls.TeacherSearch, rec.Name, P),
						scrape.DownloadTeachersTask(filepath.Join(paths.TeachersData, rec.Name), rec.Name, P),
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

func AllForOneSheet() {

}
