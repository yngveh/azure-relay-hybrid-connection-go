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

	client := &Client{
		config: c,
		url:    fmt.Sprintf("https://%s/%s", c.Namespace, c.ConnectionName),
	}

	provider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(c.KeyName, c.Key))
	if err != nil {
		return nil, err
	}

	client.token, err = provider.GetToken(client.url)

	return client, err
}

func (c *Client) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {

	req, err = http.NewRequest(method, c.url+url, body)
	if err != nil {
		return
	}

	req.Header.Add(authHeaderName, c.token.Token)

	return
}
