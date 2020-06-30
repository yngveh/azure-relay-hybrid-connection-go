package arhc

type Request struct {
	Address     string            `json:"address"`
	ID          string            `json:"id"`
	Headers     map[string]string `json:"requestHeaders"`
	Target      string            `json:"requestTarget"`
	Method      string            `json:"method"`
	Body        bool              `json:"body"`
	requestBody []byte
}

type RequestObj struct {
	Request Request `json:"request"`
}

func (r *Request) SetRequestBody(b []byte) {
	r.Body = true
	r.requestBody = b
}

func (r *Request) GetRequestBody() []byte {
	return r.requestBody
}