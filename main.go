package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/hallllll/miraiweed/ready"
	"github.com/hallllll/miraiweed/scrape"
)

var (
	err                     error
	base_url                string = "https://miraiseed7.benesse.ne.jp"
	login_url               string = "https://miraiseed7.benesse.ne.jp/seed/vw020101/displayLogin/1"
	service_url             string = "https://miraiseed7.benesse.ne.jp/seed/vw030101/displaySchoolAdminMenu"
	child_search_url        string = "https://miraiseed7.benesse.ne.jp/seed/vw030501/displaySearchChildInfo"
	child_search_refles_url string = "https://miraiseed7.benesse.ne.jp/seed/vw030501/refresh"
	search_url              string = "https://miraiseed7.benesse.ne.jp/seed/vw030501/search"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	downloadPath := filepath.Join(currentDir, "data", time.Now().Format("2006_01_02_150405000"))
	if err = os.MkdirAll(downloadPath, 0755); err != nil {
		log.Fatal(err)
	}
	csvFilePath := filepath.Join(currentDir, "info.csv")

	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		fmt.Println("The CSV file does not exist. Creating a new template...")
		if err := ready.CreateCsvTemplate(csvFilePath); err != nil {
			log.Fatalf("Failed to create the CSV template: %s", err)
		}
		fmt.Println("Template created. Please fill it with data and run the program again.")
		return
	}

	records, errChan := ready.ReadCsv(csvFilePath)
	var wg sync.WaitGroup
	for {
		select {
		case record, ok := <-records:
			if !ok {
				records = nil
			} else {
				wg.Add(1)
				go func(rec ready.Record) {
					defer wg.Done()
					fmt.Printf("Name: %s, ID: %s, PW: %s\n", rec.Name, rec.ID, rec.PW)

					// goahead

					opts := append(chromedp.DefaultExecAllocatorOptions[:],
						chromedp.Flag("headless", false),
						// chromedp.Flag("download_default_directory", downloadPath), 動くか不明
					)

					allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

					// start
					ctx, cancel := chromedp.NewContext(
						allocCtx,
						chromedp.WithLogf(log.Printf),
					)
					defer cancel()

					tasks := chromedp.Tasks{
						scrape.GetScrapeCookies(base_url),
						scrape.LoginTasks(login_url, rec.ID, rec.PW),
						scrape.NavigateTasks(child_search_url),
						scrape.DownloadTask(filepath.Join(downloadPath, rec.Name)),
					}

					err = chromedp.Run(ctx, tasks)
					if err != nil {
						log.Fatal(err)
					}

				}(record)
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				log.Fatalf("Error reading CSV: %v\n", err)
			}
		}

		if records == nil && errChan == nil {
			break
		}
	}

	// if err := ioutil.WriteFile("screenshot.png", screenshot, 0644); err != nil {
	// 	log.Fatal(err)
	// }
	// time.Sleep(1000 * time.Second)
}
