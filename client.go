package arhc

import (
	"fmt"
	"io"
	"net/http"

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

func (c *Client) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {

	wsUrl := c.url + url
	req, err = http.NewRequest(method, wsUrl, body)
	if err != nil {
		return
	}

	req.Header.Add(authHeaderName, c.token.Token)

	return
}
