package ready

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spkg/bom"
	"golang.org/x/sync/semaphore"
)

var LoginCsvFileName string = "info.csv"

type LoginRecord struct {
	Name string `csv:"Name"`
	ID   string `csv:"ID"`
	PW   string `csv:"PW"`
}

func CreateCsvTemplate(csvfilepath string) error {
	f, err := os.Create(csvfilepath)
	if err != nil {
		return err
	}
	defer f.Close()
	// csvContents := []LoginRecord{} // ヘッダーのみ

	csvContents := []LoginRecord{ // 中身も入れる場合
		{Name: "awesome_name", ID: "great_id", PW: "powerfull_pw"},
		{Name: "nice_name", ID: "good_id", PW: "special_pw"},
		{Name: "perfect_name", ID: "royal_id", PW: "complicated_pw"},
	}
	if err := gocsv.MarshalFile(&csvContents, f); err != nil {
		return err
	}
	return nil
}

func ReadCsv(csvfilepath string) (chan LoginRecord, chan error) {
	f, err := os.Open(csvfilepath)
	if err != nil {
		errChan := make(chan error, 1)
		errChan <- err
		close(errChan)
		return nil, errChan
	}
	c := make(chan LoginRecord)
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
			c <- LoginRecord{
				Name: record[0],
				ID:   record[1],
				PW:   record[2],
			}
		}
	}()

	return c, errChan
}
