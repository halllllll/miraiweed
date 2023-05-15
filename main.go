package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/chromedp/chromedp"
	"github.com/hallllll/miraiweed/ready"
	"github.com/hallllll/miraiweed/scrape"
)

var (
	err   error
	urls  *ready.URLs
	P     *ready.Put
	paths *ready.PATHs
)

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
		return // over
	}

	// 2.0 -  miraiseed instance number prompt. cuz miraiseed serving some url-s for bunch of local goverments.
	var miraiseedX string
	for {
		miraiseedX, err = ready.PromptAndRead("enter miraiseedX(1~9): ")
		if err != nil {
			log.Fatal(err)
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

	// OK, so, some settings under here

	// main urls used in this project.
	urls = ready.NewUrls()
	urls.PrepareUrl(miraiseedX)

}

func main() {

	hello()

	if err = os.MkdirAll(paths.DLStorage, 0755); err != nil {
		P.StdLog.Fatal(err)
	}
	Procces(paths)
	ready.PromptAndRead("Byebye ﾉｼ")

}

func Procces(paths *ready.PATHs) {
	// main process
	records, errChan := ready.ReadCsv(paths.LoginInfo)
	var wg sync.WaitGroup
	sm := semaphore.NewWeighted(int64(5))
	for {
		select {
		case record, ok := <-records:
			if !ok {
				records = nil
			} else {
				wg.Add(1)
				go func(rec ready.LoginRecord) {
					if err = sm.Acquire(context.Background(), 1); err != nil {
						P.ErrLog.Fatalln(err)
					}
					defer sm.Release(1)

					defer P.StdLog.Printf("%s DONE.\n", rec.Name)
					defer wg.Done()

					// goahead

					opts := append(chromedp.DefaultExecAllocatorOptions[:],
						chromedp.Flag("headless", false),
						// chromedp.Flag("download_default_directory", downloadPath), 動くか不明
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
						// scrape.NavigateStudentsTasks(urls.StudentsSearch, rec.Name, p),
						scrape.NavigateTeachersTasks(urls.TeacherSearch, rec.Name, P),
						// scrape.DownloadStudentsTask(filepath.Join(downloadPath, rec.Name), rec.Name, p),
						scrape.DownloadTeachersTask(filepath.Join(paths.DLStorage, rec.Name), rec.Name, P),
					}

					err = chromedp.Run(ctx, tasks)
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
