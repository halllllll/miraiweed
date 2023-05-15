package scrape

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ScrapeMenu struct {
}

func DownloadTask(filePath string) chromedp.Tasks {
	return chromedp.Tasks{
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).WithDownloadPath(filePath),
		chromedp.WaitVisible("#downloadExcel", chromedp.ByID),
		chromedp.WaitEnabled("#downloadExcel", chromedp.ByID),
		chromedp.Sleep(4 * time.Second),
		chromedp.Click("#downloadExcel", chromedp.ByID),
		chromedp.WaitVisible("#f30501 > div:nth-child(8)", chromedp.ByQuery),
		chromedp.WaitVisible("#download", chromedp.ByID),
		chromedp.Click("#download", chromedp.ByID),
		chromedp.Sleep(4 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Printf("downloaded at %s\n", filePath)
			return nil
		}),
	}
}

var cookies []*network.Cookie

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

func LoginTasks(login_url, login_id, login_pw string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(login_url),
		chromedp.WaitNotVisible("#loading", chromedp.ByID),
		chromedp.WaitVisible(`input[name="number"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="number"]`, login_id, chromedp.ByQuery),
		chromedp.WaitVisible(`input[name="pass"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="pass"]`, login_pw, chromedp.ByQuery),
		chromedp.Click(`input[name="inputLogin"]`, chromedp.ByQuery),
		chromedp.Sleep(1 * time.Second),
	}
}

func NavigateTasks(child_search_url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(child_search_url),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#ifGradeId`, chromedp.ByID),
		chromedp.SetValue(`#ifGradeId`, "0", chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.WaitVisible(`#ifClassId`, chromedp.ByID),
		chromedp.SetValue(`#ifClassId`, "0", chromedp.ByID),
		chromedp.Click(".searchButton", chromedp.ByQuery),
	}
}
