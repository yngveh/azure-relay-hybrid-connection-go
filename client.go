package arhc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
)

const (
	authHeaderName = "ServiceBusAuthorization"
)

type Config struct {
	Namespace      string
	ConnectionName string
	KeyName        string
	Key            string
}

type Client struct {
	token  *auth.Token
	url    string
	config *Config
}

func NewClient(c *Config) (*Client, error) {

	provider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(c.KeyName, c.Key))
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s/%s", c.Namespace, c.ConnectionName)

	token, err := provider.GetToken(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		token:  token,
		url:    url,
		config: c,
	}, nil
}

func (c *Client) Get(uri string) ([]byte, error) {

	r, statusCode, err := c.httpRequest(c.url+uri, "GET", nil)
	if err != nil {
		return nil, err
	}

	debug("StatusCode: %d ", *statusCode)

	if *statusCode > 299 {
		return nil, fmt.Errorf("statusCode: %d", *statusCode)
	}

	//var response *Response
	//if err := json.Unmarshal(r, response); err != nil {
	//	return nil, err
	//}

	debug("Response: %s", string(r))

	return r, nil
}

func (c *Client) httpRequest(uri, method string, payload []byte, ) ([]byte, *int, error) {

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add(authHeaderName, c.token.Token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		debug("Body: %s", string(body))
	}

	return body, &resp.StatusCode, err
}
