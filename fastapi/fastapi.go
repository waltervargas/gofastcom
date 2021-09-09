package fastclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type Client struct {
	apiURL string
	baseURL string
	maxRequests int
	maxTime int
	ocaURLs []string
	userAgent string
}

var (
	errScriptNotFound = errors.New("unable to get script path")
	errTokenNotFound = errors.New("unable to get token")
)

const (
	baseURL    = "https://fast.com"
	userAgent  = "github.com/waltervargas/gofastcom"
	apiURL     = "https://api.fast.com/netflix/speedtest/v2"
)

func New(maxReq, maxTime int) (*Client, error) {
	c := Client{
		baseURL: baseURL,
		userAgent: userAgent,
		maxRequests: maxReq,
		maxTime: maxTime,
	}

	t, err := getToken(baseURL)
	if err != nil {
		return nil, errTokenNotFound
	}
	c.apiURL = fmt.Sprintf("%s?https=true&token=%s", apiURL, t)

	ocas, err := c.getOCAs()
	if err != nil {
		return nil, err
	}
	c.ocaURLs = ocas

	return &c, nil
}

// TODO: Get list of urls from base url
func (c Client) getOCAs() ([]string, error) {
	var ocas []string

	resp, err := http.Get(c.apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return ocas, nil
}

func getToken(baseURL string) (string, error) {
	base, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer b.Body.Close()

	path, err := getScriptPath(resp.Body)
	if err != nil {
		return "", err
	}

	, err := http.Get(fmt.Sprintf("%s%s", url, path))
	if err != nil {
		return "", err
	}
	defer resp.


	t := getTokenFromScriptBody(b)
}

func getTokenFromScriptBody(b []byte) []byte {
	match := []byte("https:!0,token:\"")
	if bytes.Contains(b, match) {
		index := bytes.Index(b, match)
		if len(b[index+len(match):]) > 32 {
			return b[index+len(match):index+len(match)+32]
		}
	}
	return nil
}

func getScriptPath(r io.Reader) ([]byte, error) {
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err().Error() != "EOF" {
				return nil, z.Err()
			}
			return nil, errScriptNotFound
		}

		if tt == html.EndTagToken {
			continue
		}

		n, hasAttr := z.TagName()
		if hasAttr && len(n) == 6 && n[0] == 's' && n[5] == 't' {
			k, v, _ := z.TagAttr()
			if len(k) == 3 && k[0] == 's' {
				return v, nil
			}
		}
	}
}
