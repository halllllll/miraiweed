package scrape_eduboard

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/halllllll/miraiweed/ready"
)

var schoolListTableSchoolCode string = `.line tr td:nth-child(3)`
var schoolDistination string = `.line tr td:nth-child(8)`
var schoolDistinationQuerySelector func(int) string = func(num int) string {
	return fmt.Sprintf(".line > tr:nth-child(%d) > td:nth-child(8) > input:nth-child(1)", num)
}

type EduBoardScrape struct {
	P              *ready.Put
	schoolcodeList []string
}

type EduBoardScraper interface {
	NavigatingEduBoardLogin(string) chromedp.Tasks
	NavigatingEduBoardSchoolList() chromedp.Tasks
	GetSchoolCodeList() []string
	LoiterAllSchool() chromedp.ActionFunc
}

func NewEduBoardScrape(p *ready.Put) EduBoardScraper {
	return &EduBoardScrape{P: p}
}

func (eb *EduBoardScrape) GetSchoolCodeList() []string {
	return eb.schoolcodeList
}

func (eb *EduBoardScrape) NavigatingEduBoardLogin(url string) chromedp.Tasks {
	var name string
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
		chromedp.Text(".name", &name, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			eb.P.StdLog.Printf("Login sucessed - %s\n", name)
			return nil
		}),
	}
}

// TODO
func (eb *EduBoardScrape) NavigatingEduBoardSchoolList() chromedp.Tasks {
	var schoolcodeNodeList []*cdp.Node
	// var schoolDistinationList []*cdp.Node
	return chromedp.Tasks{
		// school code (不要？)
		chromedp.Nodes(schoolListTableSchoolCode, &schoolcodeNodeList, chromedp.NodeVisible, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, node := range schoolcodeNodeList {
				children := node.Children
				// nest
				for _, child := range children {
					v := child.NodeValue
					eb.schoolcodeList = append(eb.schoolcodeList, v)

				}
			}
			return nil
		}),
		// Edit Button Node info
	}
}

// TODO
func (eb *EduBoardScrape) LoiterAllSchool() chromedp.ActionFunc {
	return func(ctx context.Context) error {

		return nil
	}
}
