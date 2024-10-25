package constant

const (
	RequestSuccess = 1
	RequestFailure = 0
)

const (
	InternalError        = "50000"
	ApiMaintenance       = "50001"
	WebsocketMaintenance = "50002"
	GrpcMaintenance      = "50003"
	InvalidRequest       = "50004"
	StatusUnauthorized   = "50005"
	TooManyRequest       = "50006"
	StatusSuspended      = "50007"

	RestrictedUser   = "50100"
	PriceNotFound    = "50101"
	PriceUnavailable = "50102"
)
