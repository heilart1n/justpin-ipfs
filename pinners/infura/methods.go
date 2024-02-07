package infura

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/heilart1n/justpin-ipfs/file"
	httpretry "github.com/heilart1n/justpin-ipfs/http"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (client *Client) Name() string {
	return ClientName
}

// PinFile pins content to Infura by providing a file path, it returns an IPFS
// hash and an error.
func (client *Client) PinFile(fp string) (string, error) {
	mfr, err := file.NewMultiFileReader(fp, false, false)
	if err != nil {
		return "", fmt.Errorf("unexpected creates multipart file: %v", err)
	}
	boundary := "multipart/form-data; boundary=" + mfr.Boundary()

	return client.pinFile(mfr, boundary)
}

// PinWithReader pins content to Infura by given io.Reader, it returns an IPFS hash and an error.
func (client *Client) PinWithReader(rd io.Reader) (string, error) {
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

// PinWithBytes pins content to Infura by given byte slice, it returns an IPFS hash and an error.
func (client *Client) PinWithBytes(buf []byte) (string, error) {
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

func (client *Client) pinFile(r io.Reader, boundary string) (string, error) {
	endpoint := ApiUrl + "/api/v0/add?cid-version=1&pin=true"
	httpClient := httpretry.NewClient(client.Client)

	req, err := http.NewRequest(http.MethodPost, endpoint, r)
	if err != nil {
		return "", err
	}
	client.setAuth(req)

	req.Header.Add("Content-Type", boundary)
	req.Header.Set("Content-Disposition", `form-data; name="files"`)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// It limits anonymous requests to 12 write requests/min.
	// https://infura.io/docs/ipfs#section/Rate-Limits/API-Anonymous-Requests
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}

	var out addEvent
	dec := json.NewDecoder(resp.Body)

loop:
	for {
		var evt addEvent
		switch err := dec.Decode(&evt); err {
		case nil:
		case io.EOF:
			break loop
		default:
			return "", err
		}
		out = evt
	}

	return out.Hash, nil
}

// PinHash pins content to Infura by giving an IPFS hash, it returns the result and an error.
func (client *Client) PinHash(hash string) (bool, error) {
	if hash == "" {
		return false, fmt.Errorf("invalid hash: %s", hash)
	}

	endpoint := fmt.Sprintf("%s/api/v0/pin/add?arg=%s", ApiUrl, hash)
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return false, err
	}
	client.setAuth(req)

	httpClient := httpretry.NewClient(client.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// It limits anonymous requests to 12 write requests/min.
	// https://infura.io/docs/ipfs#section/Rate-Limits/API-Anonymous-Requests
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(data, &dat); err != nil {
		var e *json.SyntaxError
		if errors.As(err, &e) {
			return false, fmt.Errorf("json syntax error at byte offset %d", e.Offset)
		}
		return false, err
	}

	if h, ok := dat["Pins"].([]interface{}); ok && len(h) > 0 {
		return h[0] == hash, nil
	}

	return false, fmt.Errorf("pin hash to Infura failed")
}

// PinDir pins a directory to the NFT.Storage pinning service.
// It alias to PinFile.
func (client *Client) PinDir(name string) (string, error) {
	return client.PinFile(name)
}

func (client *Client) setAuth(req *http.Request) {
	if client.cfg.Apikey != "" && client.cfg.Secret != "" {
		req.SetBasicAuth(client.cfg.Apikey, client.cfg.Secret)
	}
}

func (client *Client) Pin(path interface{}) (cid string, err error) {
	err = fmt.Errorf("unsupported pinner")
	switch v := path.(type) {
	case string:
		_, err = os.Lstat(v)
		if err != nil {
			return
		}
		cid, err = client.PinFile(v)
	case io.Reader:
		cid, err = client.PinWithReader(v)
	case []byte:
		cid, err = client.PinWithBytes(v)
	}
	if err != nil {
		err = fmt.Errorf("%s: %w", client.Name(), err)
	}
	return cid, err
}
