package rajasms

import (
	"strings"
	"time"
)

type waktu time.Time

const (
	mySQLDate     = "2006-01-02"
	mySQLDateTime = "2006-01-02 15:04:05"
)

func (r *waktu) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	tWrong, err := time.Parse(mySQLDate, s)
	if err != nil {
		return err
	}

	t := tWrong.Add(time.Hour * -7)
	wib, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	*r = waktu(t.In(wib))
	return nil
}

func (r waktu) MarshalJSON() ([]byte, error) {
	t := time.Time(r)
	s := `"` + t.Format(mySQLDateTime) + `"`
	return []byte(s), nil
}
