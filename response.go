package arhc

type Response struct {
	RequestID         string            `json:"requestId"`
	StatusCode        string            `json:"statusCode"`
	StatusDescription string            `json:"statusDescription"`
	Headers           map[string]string `json:"responseHeaders"`
	Body              bool              `json:"body"`
	responseBody      []byte
}

type ResponseObj struct {
	Response Response `json:"response"`
}

func (r *Response) SetResponseBody(b []byte) {
	r.Body = true
	r.responseBody = b
}

func (r *Response) GetResponseBody() []byte {
	return r.responseBody
}
