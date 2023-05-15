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
	// miraiseed instance number prompt
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

	urls = ready.NewUrls()
	urls.PrepareUrl(miraiseedX)
}

func main() {
	P = ready.NewPut()
	P.LoggingSetting("miraiweed.log")
	// init save directory and csv data info

	paths, err := ready.NewPATHs()
	if err != nil {
		P.ErrLog.Fatal(err)
	}
	currentDir := paths.Cd
	downloadPath := paths.DLStorage
	csvFilePath := filepath.Join(currentDir, ready.LoginCsvFileName)

	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		fmt.Println("The CSV file does not exist. Creating a new template...")
		if err := ready.CreateCsvTemplate(csvFilePath); err != nil {
			P.StdLog.Fatalf("Failed to create the CSV template: %s", err)
		}
		fmt.Println("Template created. Please fill it with data and run the program again.")
		P.InfoLog.Println("Create CSV Template.")
		return
	}
	hello()

	if err = os.MkdirAll(downloadPath, 0755); err != nil {
		P.StdLog.Fatal(err)
	}
	Procces(paths)
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
	ready.PromptAndRead("Byebye ﾉｼ")
}
