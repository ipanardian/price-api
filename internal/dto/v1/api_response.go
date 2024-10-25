package dtoV1

type ResponseWrapper struct {
	Error         error       `json:"-"`
	Status        int         `json:"status"`
	StatusNumber  string      `json:"status_number"`
	StatusCode    string      `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Timestamp     int64       `json:"ts"`
	Data          interface{} `json:"data"`
}
