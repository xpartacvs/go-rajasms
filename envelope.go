package rajasms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type sms struct {
	Msisdn   string `json:"number"`
	Message  string `json:"message"`
	Schedule *waktu `json:"sendingdatetime,omitempty"`
}

type envelope struct {
	BaseUrl     string `json:"-"`
	Apikey      string `json:"apikey"`
	CallbackUrl string `json:"callbackurl,omitempty"`
	Packet      []sms  `json:"datapacket"`
}

type Envelope interface {
	SetCallbackUrl(url string) error
	Flush()
	AddSms(msisdn, message string) error
	AddScheduledSms(msisdn, message string, schedule time.Time) error
	Send() ([]Report, error)
}

func NewEnvelope(serverUrl, apikey string) (Envelope, error) {
	if err := validateCredential(serverUrl, apikey); err != nil {
		return nil, err
	}

	return &envelope{
		BaseUrl: serverUrl,
		Apikey:  apikey,
		Packet:  make([]sms, 0, 1000),
	}, nil
}

func (e *envelope) SetCallbackUrl(url string) error {
	if err := validateCredential(url, e.Apikey); err != nil {
		return err
	}
	e.CallbackUrl = url
	return nil
}

func (e *envelope) Flush() {
	e.Packet = make([]sms, 0, 1000)
}

func createSms(msisdn, message string) (*sms, error) {
	rgxMsisdn := regexp.MustCompile(`^(0|62)8[1-9]\d+$`)
	if !rgxMsisdn.MatchString(msisdn) {
		return nil, errors.New("msisdn must begin with 628 or 08")
	}

	if len(message) < 1 {
		return nil, errors.New("message length must greater than 0")
	}

	if len(message) > 480 {
		message = message[:480]
	}

	return &sms{Msisdn: msisdn, Message: message}, nil
}

func (e *envelope) AddSms(msisdn, message string) error {
	s, err := createSms(msisdn, message)
	if err != nil {
		return err
	}
	e.Packet = append(e.Packet, *s)
	return nil
}

func (e *envelope) AddScheduledSms(msisdn, message string, schedule time.Time) error {
	s, err := createSms(msisdn, message)
	if err != nil {
		return err
	}
	rTime := waktu(schedule)
	s.Schedule = &rTime
	e.Packet = append(e.Packet, *s)
	return nil
}

func (e envelope) Send() ([]Report, error) {
	reqBody, err := json.Marshal(e)
	if err != nil {
		return nil, errors.New("error on marshal the request body")
	}

	reqPost, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s", e.BaseUrl, "/sms/api_sms_masking_send_json.php"),
		bytes.NewBuffer(reqBody),
	)

	if err != nil {
		return nil, err
	}

	reqPost.Header.Add("Content-Length", fmt.Sprintf("%d", len(reqBody)))
	reqPost.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqPost)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rptResp := &reportWrapper{}
	err = json.Unmarshal(bytesResp, rptResp)
	if err != nil {
		return nil, err
	}

	if len(rptResp.Sending) <= 0 {
		return nil, ErrInvalidApiResponse
	}

	var ret []Report
	for _, r := range rptResp.Sending {
		if r.Code == 10 {
			if len(r.Packets) <= 0 {
				return nil, ErrInvalidApiResponse
			}
			ret = make([]Report, 0, len(r.Packets))
			for _, p := range r.Packets {
				ret = append(ret, p.Report)
			}
		} else {
			return nil, errors.New(strings.ToLower(r.Message))
		}
	}

	return ret, nil
}
