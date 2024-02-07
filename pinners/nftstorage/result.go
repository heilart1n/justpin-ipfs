package nftstorage

import "fmt"

type Result struct {
	hash string
	link string
}

func (result *Result) GetHash() string {
	return result.hash
}

func (result *Result) GetLink() string {
	return result.link
}

func newResult(hash string) *Result {
	return &Result{hash: hash, link: fmt.Sprintf(IPFSUrl, hash)}
}
