package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hallllll/miraiweed/compute"
	"github.com/hallllll/miraiweed/ready"
)

var (
	err   error
	bulk  int
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
		miraiseedX, err = ready.PromptAndRead("enter miraiseed[X](default=7): ")
		if err != nil {
			log.Fatal(err)
		}
		if miraiseedX == "" {
			miraiseedX = "7"
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

	// confirm concurrency number(for semaphore)
	for {
		answer, err := ready.PromptAndRead("Concarrency Limit(default=5):")
		if err != nil {
			log.Fatal(err)
		}
		if answer == "" {
			answer = "5"
		}
		bulk, err = strconv.Atoi(answer)
		if err != nil {
			log.Fatal(err)
		}
		if bulk <= 0 {
			fmt.Println("we can only accept number upper 1.")
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

	compute.Procces(paths, urls, P, bulk)
	ready.PromptAndRead("Byebye ﾉｼ")
}
