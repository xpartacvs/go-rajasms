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
	GetId() string
	GetMsisdn() string
	GetPrice() uint
	GetError() error
}

func (r report) GetError() error {
	if r.Code != 10 {
		return errors.New(strings.ToLower(r.Message))
	}
	return nil
}

func (r report) GetId() string {
	if err := r.GetError(); err != nil {
		return ""
	}
	return string(r.Id)
}

func (r report) GetMsisdn() string {
	if err := r.GetError(); err != nil {
		return ""
	}
	return r.Msisdn
}

func (r report) GetPrice() uint {
	if err := r.GetError(); err != nil {
		return 0
	}
	return uint(r.Price)
}
