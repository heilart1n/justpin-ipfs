package web3storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heilart1n/justpin-ipfs/file"
	httpretry "github.com/heilart1n/justpin-ipfs/http"
	"github.com/heilart1n/justpin-ipfs/pinners"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func (client *Client) Name() string {
	return ClientName
}

// PinFile pins content to Web3Storage by providing a file path, it returns an IPFS
// hash and an error.
func (client *Client) PinFile(fp string) (pinners.Result, error) {
	f, err := file.NewSerialFile(fp)
	if err != nil {
		return nil, err
	}
	f.MapDirectory(file.RandString(32, "lower"))

	mfr, err := file.CreateMultiForm(f, true)
	if err != nil {
		return nil, err
	}
	boundary := "multipart/form-data; boundary=" + mfr.Boundary()

	return client.pinFile(mfr, boundary)
}

// PinWithReader pins content to Web3Storage by given io.Reader, it returns an IPFS hash and an error.
func (client *Client) PinWithReader(rd io.Reader) (pinners.Result, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	fn := file.RandString(6, "lower")

	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("file", fn)
		if err != nil {
			return
		}

		if _, err = io.Copy(part, rd); err != nil {
			return
		}
	}()

	return client.pinFile(r, m.FormDataContentType())
}

// PinWithBytes pins content to Web3Storage by given byte slice, it returns an IPFS hash and an error.
func (client *Client) PinWithBytes(buf []byte) (pinners.Result, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	fn := file.RandString(6, "lower")

	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("file", fn)
		if err != nil {
			return
		}

		if _, err = part.Write(buf); err != nil {
			return
		}
	}()

	return client.pinFile(r, m.FormDataContentType())
}

func (client *Client) pinFile(r io.Reader, boundary string) (pinners.Result, error) {
	endpoint := APIUrl + "/upload"

	req, err := http.NewRequest(http.MethodPost, endpoint, r)
	if err != nil {
		return nil, err
	}
	client.setAuth(req)

	req.Header.Add("Content-Type", boundary)
	httpClient := httpretry.NewClient(client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out addEvent
	if err := json.Unmarshal(data, &out); err != nil {
		var e *json.SyntaxError
		if errors.As(err, &e) {
			return nil, fmt.Errorf("json syntax error at byte offset %d", e.Offset)
		}
		return nil, err
	}

	return client.NewResult(out.Cid), nil
}

// PinHash pins content to Web3Storage by giving an IPFS hash, it returns the result and an error.
// Note: unsupported
func (client *Client) PinHash(hash string) (bool, error) {
	return false, fmt.Errorf("not yet supported")
}

// PinDir pins a directory to the Pinata pinning service.
// It alias to PinFile.
func (client *Client) PinDir(name string) (pinners.Result, error) {
	return client.PinFile(name)
}

func (client *Client) setAuth(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+client.cfg.Apikey)
}

func (client *Client) Pin(path interface{}) (result pinners.Result, err error) {
	err = fmt.Errorf("unsupported pinner")
	switch v := path.(type) {
	case string:
		_, err = os.Lstat(v)
		if err != nil {
			return
		}
		result, err = client.PinFile(v)
	case io.Reader:
		result, err = client.PinWithReader(v)
	case []byte:
		result, err = client.PinWithBytes(v)
	}
	if err != nil {
		err = fmt.Errorf("%s: %w", client.Name(), err)
	}
	return result, err
}
