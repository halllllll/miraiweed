package ready

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/sync/semaphore"

	"github.com/gocarina/gocsv"
	"github.com/spkg/bom"
)

type Record struct {
	Name string `csv:"Name"`
	ID   string `csv:"ID"`
	PW   string `csv:"PW"`
}

type URLs struct {
	Num                    string
	Base                   string
	Login                  string
	Service                string
	DisplaySchoolAdminMenu string
	StudentsSearch         string
	StudentsSearchReflesh  string
	TeacherSearch          string
	Search                 string
}

func NewUrls() *URLs {
	return &URLs{}
}

func (u *URLs) PrepareUrl(num string) {
	u.Num = num
	u.Base = fmt.Sprintf("https://miraiseed%s.benesse.ne.jp", u.Num)
	u.Login = fmt.Sprintf("%s/seed/vw020101/displayLogin/1", u.Base)
	u.Service = fmt.Sprintf("%s/seed/vw030101/displaySchoolAdminMenu", u.Base)
	u.StudentsSearch = fmt.Sprintf("%s/seed/vw030501/displaySearchChildInfo", u.Base)
	u.StudentsSearchReflesh = fmt.Sprintf("%s/seed/vw030501/refresh", u.Base)
	u.Search = fmt.Sprintf("%s/seed/vw030501/search", u.Base)
	u.TeacherSearch = fmt.Sprintf("%s/seed/vw030401/", u.Base)
}

func CreateCsvTemplate(csvfilepath string) error {
	f, err := os.Create(csvfilepath)
	if err != nil {
		return err
	}
	defer f.Close()
	// csvContents := []Record{} // ヘッダーのみ
	// 中身も入れる場合
	csvContents := []Record{
		{Name: "awesome_name", ID: "great_id", PW: "powerfull_pw"},
		{Name: "nice_name", ID: "good_id", PW: "special_pw"},
		{Name: "perfect_name", ID: "royal_id", PW: "complicated_pw"},
	}
	if err := gocsv.MarshalFile(&csvContents, f); err != nil {
		return err
	}
	return nil
}

func ReadCsv(csvfilepath string) (chan Record, chan error) {
	f, err := os.Open(csvfilepath)
	if err != nil {
		errChan := make(chan error, 1)
		errChan <- err
		close(errChan)
		return nil, errChan
	}
	c := make(chan Record)
	errChan := make(chan error, 1)
	// gocsvのMarshalToChanだとうまく逐次処理が書けなかった（別ファイルからだと難しいのか？）

	sm := semaphore.NewWeighted(int64(4))

	go func() {
		defer f.Close()
		defer close(c)
		defer close(errChan)
		if err = sm.Acquire(context.Background(), 1); err != nil {
			log.Println(err)
		}
		defer sm.Release(1)

		reader := csv.NewReader(bom.NewReader(f))
		// headerを捨てる
		_, _ = reader.Read()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- err
				return
			}
			if len(record) != 3 {
				errChan <- fmt.Errorf("invalid record: %v", record)
				return
			}
			c <- Record{
				Name: record[0],
				ID:   record[1],
				PW:   record[2],
			}
		}
	}()

	return c, errChan
}

func PromptAndRead(message string) (string, error) {
	fmt.Fprint(os.Stdout, message)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

type Put struct {
	PutWriter io.Writer
	StdLog    *log.Logger
	ErrLog    *log.Logger
	InfoLog   *log.Logger
	Panic     *log.Logger
}

func NewPut() *Put {
	return &Put{}
}

func (p *Put) LoggingSetting(fileName string) *io.Writer {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	p.PutWriter = io.MultiWriter(f, os.Stderr)
	p.StdLog = log.New(f, "[Std] ", log.Ldate|log.Ltime)
	p.ErrLog = log.New(f, "[Error] ", log.Ldate|log.Ltime)
	p.InfoLog = log.New(f, "[Info] ", log.Ldate|log.Ltime)
	p.Panic = log.New(f, "[PANIC] ", log.Ldate|log.Ltime)
	p.StdLog.SetOutput(p.PutWriter)
	p.ErrLog.SetOutput(p.PutWriter)
	p.InfoLog.SetOutput(p.PutWriter)
	p.Panic.SetOutput(p.PutWriter)
	return &p.PutWriter
}
