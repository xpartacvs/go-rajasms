package rajasms

import "time"

type inquiry struct {
	apiResponse
	Balance cacah `json:"Balance,omitempty"`
	Expired waktu `json:"Expired,omitempty"`
}

type inquiryResponse struct {
	Inquiries []inquiry `json:"balance_respon"`
}

type Inquiry interface {
	GetBalance() uint
	GetExpiry() time.Time
}

func (i inquiry) GetBalance() uint {
	return uint(i.Balance)
}

func (i inquiry) GetExpiry() time.Time {
	return time.Time(i.Expired)
}
