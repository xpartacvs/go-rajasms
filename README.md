# Go RajaSMS
Sebuah paket pustaka golang untuk berinteraksi dengan API RajaSMS

## Cara Pakai

Download dulu packagenya ke dalam project:

```bash
go get github.com/xpartacvs/go-rajasms
```

Lalu mulai koding:

```go
package main

import (
    "fmt"
    "log"

    "github.com/xpartacvs/go-rajasms"
)

func main() {
    // Instansiasi client
    client, err := rajasms.NewCient("http://url-server-rajasms.mu", "apikey-akunmu")
    if err != nil {
        log.Fatalln(err.Error())
    }

    // Ambil info akun
    info, err := client.AccountInfo()
    if err != nil {
        log.Fatalln(err.Error())
    }

    // Print info akun
    fmt.Println("Saldo\t\t: ", info.Balance)
    fmt.Println("Kedaluarsa\t: ", info.Expiry.Format("2006-01-02 15:04:05"))

    // Buat batch dan isi dengan SMS
    paket := rajasms.Batch{}
    paket.AddSMS("081xxxxxxxxx", "Testing SMS RajaSMS dari package golang")

    // Isi lagi dengan SMS terjadwal
    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        log.Fatalln(err.Error())
    }
    tenMinutesFromNow := time.Now().In(loc).Add(10*time.Minute)
    paket.AddScheduledSMS("6285xxxxxxxxx", "Testing SMS terjadwal RajaSMS dari package golang", tenMinutesFromNow)

    // Kirim ke server RajaSMS
    reports, err := client.Send(paket)
    if err != nil {
        log.Fatalln(err.Error())
    }

    // Iterate report kiriman
    for _, r := range reports {
        if r.Error == nil {
            fmt.Println("Id\t: ", r.Id)
            fmt.Println("MSISDN\t: ", r.Msisdn)
            fmt.Println("Price\t: ", r.Price)
        }
	}

    // Bisa juga kirim dengan melampirkan callback url untuk pantau perubahan status kiriman
    _, _ = client.SendWithCallbackURL(paket, "http://domain.tld/path/webhook/kamu")
}
```

Demikian manteman, selamat mencoba.