package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/chromedp/chromedp"
	"github.com/hallllll/miraiweed/ready"
	"github.com/hallllll/miraiweed/scrape"
)

var (
	err  error
	urls *ready.URLs
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
	hello()

	// init save directory and csv data info
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	downloadPath := filepath.Join(currentDir, "data", strings.ReplaceAll(time.Now().Format("2006_01_02_150405.000"), ".", "_"))
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

	// main process
	records, errChan := ready.ReadCsv(csvFilePath)
	var wg sync.WaitGroup
	sm := semaphore.NewWeighted(int64(5))
	for {
		select {
		case record, ok := <-records:
			if !ok {
				records = nil
			} else {

				wg.Add(1)
				go func(rec ready.Record) {
					if err = sm.Acquire(context.Background(), 1); err != nil {
						log.Println(err)
					}
					defer sm.Release(1)

					defer fmt.Printf("%s DONE.\n", rec.Name)
					defer wg.Done()

					// goahead

					opts := append(chromedp.DefaultExecAllocatorOptions[:],
						chromedp.Flag("headless", true),
						// chromedp.Flag("download_default_directory", downloadPath), 動くか不明
					)

					allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
					allocCtx2, cancel := context.WithTimeout(allocCtx, 60*time.Second)

					// start
					ctx, cancel := chromedp.NewContext(
						allocCtx2,
						chromedp.WithLogf(log.Printf),
					)
					defer cancel()

					tasks := chromedp.Tasks{
						scrape.GetScrapeCookies(urls.Base),
						scrape.LoginTasks(urls.Login, rec.ID, rec.PW),
						scrape.NavigateTasks(urls.ChildSearch),
						scrape.DownloadTask(filepath.Join(downloadPath, rec.Name)),
					}

					err = chromedp.Run(ctx, tasks)
					if err != nil {
						log.Println(err)
					}

				}(record)
				/*
					fmt.Printf("TARGET: %s\n", record.Name)
					opts := append(chromedp.DefaultExecAllocatorOptions[:],
						chromedp.Flag("headless", true),
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
						scrape.LoginTasks(login_url, record.ID, record.PW),
						scrape.NavigateTasks(child_search_url),
						scrape.DownloadTask(filepath.Join(downloadPath, record.Name)),
					}

					err = chromedp.Run(ctx, tasks)
					if err != nil {
						log.Println(err)
					}
				*/
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
	wg.Wait()
	ready.PromptAndRead("Byeby ﾉｼ")
}
