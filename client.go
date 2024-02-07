package justpin_ipfs

import (
	"github.com/heilart1n/justpin-ipfs/config"
	"github.com/heilart1n/justpin-ipfs/pinners"
	"github.com/heilart1n/justpin-ipfs/pinners/infura"
	"github.com/heilart1n/justpin-ipfs/pinners/nftstorage"
	"github.com/heilart1n/justpin-ipfs/pinners/pinata"
	"github.com/heilart1n/justpin-ipfs/pinners/web3storage"
	"net/http"
)

type ClientName string

const (
	ClientNameInfura      ClientName = "Infura"
	ClientNameNFTStorage  ClientName = "NFTStorage"
	ClientNamePinata      ClientName = "Pinata"
	ClientNameWeb3Storage ClientName = "Web3Storage"
)

type Pinners struct {
	Infura      pinners.Pinner
	NFTStorage  pinners.Pinner
	Pinata      pinners.Pinner
	Web3Storage pinners.Pinner
}

func NewPinners(infuraCfg, nftStorageCfg, pinataCfg, web3StorageCfg config.Config) *Pinners {
	return &Pinners{
		Infura:      infura.NewClient(infuraCfg, http.DefaultClient),
		NFTStorage:  nftstorage.NewClient(nftStorageCfg, http.DefaultClient),
		Pinata:      pinata.NewClient(pinataCfg, http.DefaultClient),
		Web3Storage: web3storage.NewClient(web3StorageCfg, http.DefaultClient),
	}
}
