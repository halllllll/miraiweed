package scraping

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/halllllll/miraiweed/ready"
)

var cookies []*network.Cookie

// Handle the modal that appears during the fiscal year update period
var UntilNotDoneAnnuralUpdateCssSelector string = `#ui-id-15 > div:nth-child(1) > div:nth-child(2) > span:nth-child(2) > input:nth-child(1)`

func DownloadStudentsTask(filePath, login_name string, p *ready.Put) chromedp.Tasks {
	p.StdLog.Printf("%s Students Download Challenge\n", login_name)

	return chromedp.Tasks{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(filePath),
		chromedp.WaitVisible("#downloadExcel", chromedp.ByID),
		chromedp.WaitEnabled("#downloadExcel", chromedp.ByID),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click("#downloadExcel", chromedp.ByID),
		// anually update modal
		chromedp.ActionFunc(func(ctx context.Context) error {
			// checking modal view
			var stillNONGoingAnuallyUpdate bool
			err := chromedp.Evaluate(fmt.Sprintf(`%s.offsetParent !== null`, fmt.Sprintf("document.querySelector('%s')", UntilNotDoneAnnuralUpdateCssSelector)), &stillNONGoingAnuallyUpdate).Do(ctx)
			if err != nil {
				return err
			}
			if stillNONGoingAnuallyUpdate {
				stillNOUpdateingStatus(filePath).Do(ctx)
			}
			normalDownloadStatus(filePath).Do(ctx)

			return nil
		}),
	}
}

func GetScrapeCookies(base_url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookieParams := network.GetCookies().WithUrls([]string{base_url})
			var err error
			cookies, err = cookieParams.Do(ctx)
			return err
		}),
	}
}

func SetScrapeCookies() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, cookie := range cookies {
				network.SetCookie(cookie.Name, cookie.Value)
			}
			return nil
		}),
	}
}

func LoginTasks(login_url, login_name, login_id, login_pw string, p *ready.Put) chromedp.Tasks {
	p.StdLog.Printf("%s Login Challenge\n", login_name)
	return chromedp.Tasks{
		chromedp.Navigate(login_url),
		chromedp.WaitNotVisible(`#loading`, chromedp.ByID),
		chromedp.WaitVisible(`input[name="number"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="number"]`, login_id, chromedp.ByQuery),
		chromedp.WaitVisible(`input[name="pass"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="pass"]`, login_pw, chromedp.ByQuery),
		chromedp.Click(`input[name="inputLogin"]`, chromedp.ByQuery),
		chromedp.Sleep(1 * time.Second),
	}
}

func NavigateStudentsTasks(student_search_url, login_name string, p *ready.Put) chromedp.Tasks {
	p.StdLog.Printf("%s Students Loitering...\n", login_name)
	return chromedp.Tasks{
		chromedp.Navigate(student_search_url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#ifGradeId`, chromedp.ByID),
		chromedp.SetValue(`#ifGradeId`, "0", chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#ifClassId`, chromedp.ByID),
		chromedp.SetValue(`#ifClassId`, "0", chromedp.ByID),
		chromedp.Click(".searchButton", chromedp.ByQuery),
	}
}

func NavigateTeachersTasks(teacher_search_url, login_name string, p *ready.Put) chromedp.Tasks {
	p.StdLog.Printf("%s Teacher Loitering...\n", login_name)
	return chromedp.Tasks{
		chromedp.Navigate(teacher_search_url),
		chromedp.Sleep(1 * time.Second),
	}
}

func DownloadTeachersTask(filePath, login_name string, p *ready.Put) chromedp.Tasks {
	p.StdLog.Printf("%s Teachers Download Challenge\n", login_name)

	return chromedp.Tasks{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(filePath),
		chromedp.WaitVisible("#downloadExcel", chromedp.ByID),
		chromedp.WaitEnabled("#downloadExcel", chromedp.ByID),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click("#downloadExcel", chromedp.ByID),
		chromedp.WaitVisible("#download", chromedp.ByID),
		chromedp.Click("#download", chromedp.ByID),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Printf("teacher excel downloaded : %s\n", filepath.Base(filePath))
			return nil
		}),
	}
}

func normalDownloadStatus(filePath string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible("#f30501 > div:nth-child(8)", chromedp.ByQuery),
		chromedp.WaitVisible("#download", chromedp.ByID),
		chromedp.Click("#download", chromedp.ByID),
		chromedp.Sleep(4 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Printf("students excel downloaded: %s\n", filepath.Base(filePath))
			return nil
		}),
	}
}

func stillNOUpdateingStatus(filePath string) chromedp.Tasks {
	// cilck button inner modal
	return chromedp.Tasks{
		chromedp.WaitVisible(UntilNotDoneAnnuralUpdateCssSelector, chromedp.ByQuery),
		chromedp.WaitEnabled(UntilNotDoneAnnuralUpdateCssSelector, chromedp.ByQuery),
		chromedp.Click(UntilNotDoneAnnuralUpdateCssSelector, chromedp.ByQuery),
		chromedp.Sleep(4 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Printf("excel found (but not anual update): %s\n", filepath.Base(filePath))
			return nil
		}),
	}
}
