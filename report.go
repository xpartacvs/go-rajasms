package rajasms

import (
	"errors"
	"strings"
)

type report struct {
	Code    uint8  `json:"sendingstatus"`
	Message string `json:"sendingstatustext"`
	Id      resi   `json:"sendingid"`
	Msisdn  string `json:"number"`
	Price   cacah  `json:"price"`
}

type packet struct {
	Report report `json:"packet"`
}

type reportResponse struct {
	apiResponse
	Packets []packet `json:"datapacket"`
}

type reportWrapper struct {
	Sending []reportResponse `json:"sending_respon"`
}

type Report interface {
	GetId() (string, error)
	GetMsisdn() (string, error)
	GetPrice() (uint, error)
}

func (r report) GetError() error {
	if r.Code != 10 {
		return errors.New(strings.ToLower(r.Message))
	}
	return nil
}

func (r report) GetId() (string, error) {
	if err := r.GetError(); err != nil {
		return "", err
	}
	return string(r.Id), nil
}

func (r report) GetMsisdn() (string, error) {
	return r.Msisdn, r.GetError()
}

func (r report) GetPrice() (uint, error) {
	if err := r.GetError(); err != nil {
		return 0, err
	}
	return uint(r.Price), nil
}
