package rajasms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type client struct {
	baseUrl string
	md5Key  string
}

type Client interface {
	GetInquiry() (Inquiry, error)
}

func (c client) GetInquiry() (Inquiry, error) {
	reqBody, err := json.Marshal(
		map[string]string{
			"apikey": c.md5Key,
		},
	)
	if err != nil {
		return nil, errors.New("error on marshal the request body")
	}

	endpoint := fmt.Sprintf("%s%s", c.baseUrl, "/sms/api_sms_masking_balance_json.php")
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	inquiryResp := &inquiryResponse{}
	err = json.Unmarshal(bytesResp, inquiryResp)
	if err != nil {
		return nil, err
	}

	if len(inquiryResp.Inquiries) <= 0 {
		return nil, ErrInvalidApiResponse
	}

	for _, i := range inquiryResp.Inquiries {
		if i.Code == 10 {
			return i, nil
		} else {
			return nil, errors.New(strings.ToLower(i.Message))
		}
	}

	return nil, ErrInvalidApiResponse
}

func NewClient(serverUrl, apikey string) Client {
	if err := validateCredential(serverUrl, apikey); err != nil {
		panic(err)
	}
	return &client{
		baseUrl: serverUrl,
		md5Key:  apikey,
	}
}
