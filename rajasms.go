package rajasms

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type waktu time.Time
type cacah uint
type resi string

type apiResponse struct {
	Code    uint8  `json:"globalstatus"`
	Message string `json:"globalstatustext"`
}

var ErrInvalidApiResponse error = errors.New("invalid api response")

func validateCredential(url, key string) error {
	rgxMd5 := regexp.MustCompile("^[a-f0-9]{32}$")
	if !rgxMd5.MatchString(key) {
		return errors.New("invalid api key")
	}

	rgxUrl := regexp.MustCompile("^https?://.*")
	if !rgxUrl.MatchString(url) {
		return errors.New("invalid base url")
	}
	return nil
}

func (r *waktu) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return err
	}

	*r = waktu(t)
	return nil
}

func (r waktu) MarshalJSON() ([]byte, error) {
	t := time.Time(r)
	s := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(s), nil
}

func (r *cacah) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	intBal, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	uintBal := uint(intBal)
	*r = cacah(uintBal)
	return nil
}

func (r *resi) UnmarshalJSON(b []byte) error {
	strId := string(b)
	*r = resi(strId)
	return nil
}
