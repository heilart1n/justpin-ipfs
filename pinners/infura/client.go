package infura

import (
	"github.com/heilart1n/justpin-ipfs/config"
	"net/http"
)

const (
	ApiUrl     = "https://ipfs.infura.io:5001"
	ClientName = "Infura"
)

// Client represents an Infura configuration. If there is no Apikey or
// Secret, it will make API calls using anonymous requests.
type Client struct {
	*http.Client
	cfg        config.Config
	clientName string
}

func NewClient(cfg config.Config, httpClient *http.Client) *Client {
	return &Client{cfg: cfg, clientName: ClientName, Client: httpClient}
}
