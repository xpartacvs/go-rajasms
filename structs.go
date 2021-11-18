package rajasms

import "time"

// Account adalah struct yang berisi informasi keadaan akun RajaSMS Anda yang terdiri dari jumlah saldo dan tanggal kedaluarsa.
// Account adalah type yang dikembalikan oleh (Client).AccountInfo() jika tidak ada error. Khusus untuk field Expiry,
// secara default location dari field tersebut adalah Asia/Jakarta (WIB)
type Account struct {
	// Jumlah saldo akun
	Balance uint64
	// Waktu kedaluarsa akun dalam WIB
	Expiry time.Time
}

// Report adalah struct yang berisi informasi hasil pengiriman SMS ke pool RajaSMS.
// Report adalah type yang dikembalikan oleh (Client).Send() ataupun (Client).SendWithCallbak() jika tidak ada error.
// Penting untuk meng-cast field Error terlebih dahulu, jika field Error bernilai selain nil maka filed-field lainnya
// idealnya akan bernilai non-zero value.
type Report struct {
	// Jika terjadi error, nilainya bukan nil.
	Error error
	// Report ID dari pengiriman
	Id string
	// Nomor hape
	Msisdn string
	// Harga yang dikenakan pada pengiriman SMS terkait.
	Price uint
}
