package rajasms

import "errors"

var (
	// List global error terkait api call RajaSMS secara umum
	ErrGlobalInvalidJSON          error = errors.New("json post error")
	ErrGlobalCheckBalance         error = errors.New("error check balance")
	ErrGlobalUnregisteredAPIKey   error = errors.New("apikey not register")
	ErrGlobalRateLimitExceed      error = errors.New("ratelimit exceeded")
	ErrGlobalUnregisteredClientIP error = errors.New("ip address not register")
	ErrGlobalBalanceExpired       error = errors.New("expire balance")
	ErrGlobalMaxDataLimirExceeded error = errors.New("maximum data")
	ErrGlobalNoData               error = errors.New("no data")

	// List error yang berkaitan dengan api call pengiriman SMS
	ErrSendingInvalidNumber      error = errors.New("invalid msisdn")
	ErrSendingInvalidBody        error = errors.New("invalid sms body")
	ErrSendingInsuficientBalance error = errors.New("insuficient balance")
	ErrSendingSystemError        error = errors.New("system error")

	errorMap map[uint8]error = map[uint8]error{
		20:  ErrGlobalInvalidJSON,
		25:  ErrGlobalCheckBalance,
		30:  ErrGlobalUnregisteredAPIKey,
		35:  ErrGlobalRateLimitExceed,
		40:  ErrGlobalUnregisteredClientIP,
		50:  ErrGlobalBalanceExpired,
		204: ErrGlobalNoData,
		55:  ErrGlobalMaxDataLimirExceeded,
		60:  ErrSendingInvalidNumber,
		70:  ErrSendingInvalidBody,
		80:  ErrSendingInsuficientBalance,
		90:  ErrSendingSystemError,
	}
)
