package rajasms

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Client adalah type utama yang akan berinteraksi dengan API RajaSMS
type Client struct {
	baseUrl *url.URL
	key     string
}

// NewClient adalah fungsi untuk membuat type Client dalam wujud pointer.
// Argumen baseUrl adalah URL server akun RajaSMS Anda dan ini haruslah diawali dengan http:// atau https://.
// Argumen apiKey adalah API key dari akun RajaSMS Anda.
func NewCient(baseUrl, apiKey string) (*Client, error) {
	if !regexp.MustCompile(`^https?://.*`).MatchString(baseUrl) {
		return nil, errors.New("invalid base url")
	}

	bUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	bUrl.Fragment = ""
	bUrl.RawQuery = ""
	bUrl.Path = strings.TrimRight(bUrl.Path, "/ ")

	return &Client{baseUrl: bUrl, key: apiKey}, nil
}

// AccountInfo adalah method untuk mendapatkan detail informasi dari akun RajaSMS Anda berupa
// saldo terakhir dan tanggal kedaluarsa akun.
func (c Client) AccountInfo() (*Account, error) {
	reqBody, err := json.Marshal(
		map[string]string{
			"apikey": c.key,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.baseUrl.String()+"/sms/api_sms_masking_balance_json.php", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var buff struct {
		Inquiries []struct {
			Code    uint8  `json:"globalstatus"`
			Status  string `json:"globalstatustext"`
			Balance cacah  `json:"Balance,omitempty"`
			Expired waktu  `json:"Expired,omitempty"`
		} `json:"balance_respon"`
	}
	err = json.Unmarshal(bytesResp, &buff)
	if err != nil {
		return nil, err
	}

	if len(buff.Inquiries) <= 0 {
		return nil, errorMap[204]
	}

	if buff.Inquiries[0].Code != 10 {
		return nil, errorMap[buff.Inquiries[0].Code]
	}

	return &Account{
		Balance: uint64(buff.Inquiries[0].Balance),
		Expiry:  time.Time(buff.Inquiries[0].Expired),
	}, nil
}

// Send adalah method untuk mengirimkan SMS yang ada didalam Batch.
// SMS yang diterima oleh server RajaSMS akan langsung dikirimkan sesegera mungkin ke nomor tujuan.
// Jika tidak terjadi error, method ini akan mengembalikan nilai berupa slice dari type Report berisi
// informasi hasil pegniriman untuk masing-masing SMS.
func (c Client) Send(batch Batch) ([]Report, error) {
	return c.SendWithCallbackURL(batch, "")
}

// SendWithCallbackURL adalah method yang berfungsi sama dengan method Send hanya saja kita melampirkan
// URL (webhook) yang akan digunakan oleh server RajaSMS untuk mengupdate info terkait perubahan status
// kiriman masing-masing SMS. Argumen callbackUrl haruslah diawali dengan http:// atau https://.
func (c Client) SendWithCallbackURL(batch Batch, callbackUrl string) ([]Report, error) {
	if len(batch) <= 0 {
		return nil, errors.New("batch tidak boleh kosong")
	}

	cb := ""
	if len(strings.TrimSpace(callbackUrl)) > 0 {
		if !regexp.MustCompile(`^https?://.*`).MatchString(callbackUrl) {
			return nil, errors.New("invalid callback url")
		}
		cb = strings.TrimSpace(callbackUrl)
	}

	var data struct {
		ApiKey      string `json:"apikey"`
		CallbackUrl string `json:"callbackurl,omitempty"`
		Packages    Batch  `json:"datapacket"`
	}

	data.ApiKey = c.key
	data.Packages = batch
	if len(cb) > 0 {
		data.CallbackUrl = cb
	}

	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reqPost, err := http.NewRequest(
		http.MethodPost,
		c.baseUrl.String()+"/sms/api_sms_masking_send_json.php",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}

	reqPost.Header.Add("Content-Length", strconv.FormatInt(int64(len(reqBody)), 10))
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

	var reports struct {
		Responses []struct {
			Code   uint8  `json:"globalstatus"`
			Status string `json:"globalstatustext"`
			Data   []struct {
				Response struct {
					Code   uint8  `json:"sendingstatus"`
					Status string `json:"sendingstatustext"`
					Id     uint64 `json:"sendingid"`
					Msisdn string `json:"number"`
					Price  uint   `json:"price"`
				} `json:"packet"`
			} `json:"datapacket"`
		} `json:"sending_respon"`
	}

	err = json.Unmarshal(bytesResp, &reports)
	if err != nil {
		return nil, err
	}

	if len(reports.Responses) != 1 {
		return nil, ErrGlobalNoData
	}

	if reports.Responses[0].Code != 10 {
		return nil, errorMap[reports.Responses[0].Code]
	}

	if len(reports.Responses[0].Data) <= 0 {
		return nil, ErrGlobalNoData
	}

	var ret []Report
	for _, d := range reports.Responses[0].Data {
		r := Report{Msisdn: d.Response.Msisdn}
		if d.Response.Code == 10 {
			r.Error = nil
			r.Id = strconv.FormatUint(d.Response.Id, 10)
			r.Price = d.Response.Price
		} else {
			r.Error = errorMap[d.Response.Code]
			r.Price = 0
		}
		ret = append(ret, r)
	}

	return ret, nil
}
