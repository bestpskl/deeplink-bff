package constant

type Code string

func (c Code) String() string {
	return string(c)
}

var (
	CodeSuccess                    Code = "DL0000"
	CodeInvalidCommonFields        Code = "DL4000"
	CodePartnerConfigNotExist      Code = "DL4091"
	CodeInvalidDynamicFields       Code = "DL4092"
	CodeDuplicatePartnerTxnRef     Code = "DL4093"
	CodeSessionValidUntilTooOld    Code = "DL4094"
	CodeTransactionNotExist        Code = "DL4040"
	CodeInvalidDeeplink            Code = "DL4020"
	CodeDeeplinkExpired            Code = "DL4021"
	CodeInvalidDeeplinkTransaction Code = "DL4023"
	CodeUnprocessEntity            Code = "DL4222"
	CodeInternal                   Code = "DL9999"
)
