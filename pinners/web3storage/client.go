package web3storage

import (
	"github.com/heilart1n/justpin-ipfs/config"
	"net/http"
)

const (
	APIUrl     = "https://api.web3.storage"
	ClientName = "Web3Storage"
	IPFSUrl    = "https://w3s.link/ipfs/%s"
)

// Client Web3Storage represents a Web3Storage configuration.
type Client struct {
	*http.Client
	cfg        config.Config
	clientName string
}

func NewClient(cfg config.Config, httpClient *http.Client) *Client {
	return &Client{cfg: cfg, clientName: ClientName, Client: httpClient}
}
