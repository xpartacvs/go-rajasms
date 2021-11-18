package rajasms

import (
	"regexp"
	"time"
)

type sms struct {
	Msisdn   string `json:"number"`
	Body     string `json:"message"`
	Schedule string `json:"sendingdatetime,omitempty"`
}

// Batch type bisa dianggap semacam manifest untuk mengirimkan satu atau pun sekelompok SMS.
type Batch []sms

var (
	rgxMsisdn = regexp.MustCompile(`^(\+?62|0)8[1-9]{2}\d{6,9}`)
)

// AddSMS adalah method untuk menambahkan SMS baru kedalam Batch.
// Argument number harusnya dalam format MSISDN nomor selular Indonesia (+628.., 628.., 08..).
// Argument message adalah pesan SMS yang hendak dikirimkan. Satu pesan SMS terdiri dari maksimum
// 160 karakter. Jika lebih saldo yang tepotong pada akun Anda akan dihitung berdasarkan kelipatannya.
func (e *Batch) AddSMS(number, message string) {
	if rgxMsisdn.MatchString(number) {
		s := sms{
			Msisdn: number,
			Body:   message,
		}
		*e = append(*e, s)
	}
}

// AddScheduledSMS sama dengan method AddSMS, hanya saja ditambah dengan stempel waktu kapan SMS hendak dikirmkan.
// Location yang disarankan untuk argumen sendtime adalah "Asia/Jakarta" (zona waktu WIB)
func (e *Batch) AddScheduledSMS(number, message string, sendtime time.Time) {
	if rgxMsisdn.MatchString(number) {
		s := sms{
			Msisdn:   number,
			Body:     message,
			Schedule: sendtime.Format(mySQLDateTime),
		}
		*e = append(*e, s)
	}
}

// Reset adalah method untuk mengosongkan Batch
func (e *Batch) Reset() {
	*e = nil
}
