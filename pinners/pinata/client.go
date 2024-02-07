package pinata

import (
	"github.com/heilart1n/justpin-ipfs/config"
	"net/http"
)

const (
	PinFileUrl = "https://api.pinata.cloud/pinning/pinFileToIPFS"
	PinHashUrl = "https://api.pinata.cloud/pinning/pinByHash"
	ClientName = "Pinata"
	IPFSUrl    = "https://gateway.pinata.cloud/ipfs/%s"
)

// Client Pinata represents a Pinata configuration.
type Client struct {
	*http.Client
	cfg        config.Config
	clientName string
}

func NewClient(cfg config.Config, httpClient *http.Client) *Client {
	return &Client{cfg: cfg, clientName: ClientName, Client: httpClient}
}
