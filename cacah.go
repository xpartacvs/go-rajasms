package rajasms

import (
	"strconv"
	"strings"
)

type cacah uint64

func (ui64 cacah) MarshalJSON() ([]byte, error) {
	ui64Val := uint64(ui64)
	s := `"` + strconv.FormatUint(ui64Val, 10) + `"`
	return []byte(s), nil
}

func (ui64 *cacah) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	ui64Val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*ui64 = cacah(ui64Val)
	return nil
}
