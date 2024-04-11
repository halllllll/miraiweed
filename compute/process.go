package compute

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/chromedp/chromedp"
	"github.com/halllllll/miraiweed/ready"
	"github.com/halllllll/miraiweed/scraping"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/semaphore"
)

func Procces(paths *ready.PATHs, urls *ready.URLs, P *ready.Put, bulk int) {
	// main process
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
						P.ErrLog.Println(err)
					}
					defer sm.Release(1)

					defer P.StdLog.Printf("%s DONE.\n", rec.Name)
					defer wg.Done()

					// goahead
					opts := append(chromedp.DefaultExecAllocatorOptions[:],
						chromedp.Flag("headless", true),
					)

					tasks := chromedp.Tasks{
						scraping.GetScrapeCookies(urls.Base),
						scraping.LoginTasks(urls.Login, rec.Name, rec.ID, rec.PW, P),
						scraping.NavigateStudentsTasks(urls.StudentsSearch, rec.Name, P),
						scraping.DownloadStudentsTask(filepath.Join(paths.StudentFolder(), rec.Name), rec.Name, P),
						scraping.NavigateTeachersTasks(urls.TeacherSearch, rec.Name, P),
						scraping.DownloadTeachersTask(filepath.Join(paths.TeacherFolder(), rec.Name), rec.Name, P),
					}

					operation := func() error {
						allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
						allocCtx, cancel = context.WithTimeout(allocCtx, 30*time.Second)

						// start
						ctx, cancel := chromedp.NewContext(
							allocCtx,
							chromedp.WithLogf(P.StdLog.Printf),
						)
						defer cancel()
						err := chromedp.Run(ctx, tasks)
						return err
					}
					b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 2)

					if err := backoff.Retry(operation, b); err != nil {
						err := fmt.Errorf("chromedp error occured during【%s】around ", rec.Name)
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

func AllForOneSheet(path string, h []string, targetSheetName string, P *ready.Put) error {
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
	_h := make([]interface{}, len(h))
	for i, v := range h {
		_h[i] = v
	}
	if err = sw.SetRow(header, _h); err != nil {
		return err
	}
	allElsx := filepath.Join(path, "all.xlsx")
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "walkdir failed")
		}
		if d.IsDir() || path == allElsx || filepath.Ext(path) != ".xlsx" {
			return nil
		}

		dir, _ := filepath.Split(path)
		P.StdLog.Printf("excelize - %s", filepath.Base(dir))
		targetxlsx, err := excelize.OpenFile(path)
		if err != nil {
			return fmt.Errorf("failed target xlsx at %s - %w", path, err)
		}
		rows, err := targetxlsx.Rows(targetSheetName)
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
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := sw.Flush(); err != nil {
		return fmt.Errorf("failed sream writer flush - %w", err)
	}
	if err := xlsx.SaveAs(allElsx); err != nil {
		return fmt.Errorf("faield to save xlsx - %w", err)
	}
	P.StdLog.Printf("save xlsx at %s\n", allElsx)
	return nil
}
