package ready

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spkg/bom"
)

type Record struct {
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
	// gocsvだとうまく逐次処理が書けなかった（別ファイルからだと難しいのか？）
	go func() {
		defer f.Close()
		defer close(c)
		defer close(errChan)

		reader := csv.NewReader(bom.NewReader(f))
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
