package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/yadunut/file-sync/internal/contracts"
)

type Client struct {
	log    *log.Logger
	config contracts.Config
}

func NewClient() *Client {
	return &Client{log: log.Default()}
}

func (c *Client) get(url string) (*http.Response, error) {
	return http.Get(fmt.Sprintf("http://%s/%s", c.config.GetUrl(), url))
}

func (c *Client) Version() {
	res, err := c.get("/version")
	if err != nil {
		log.Println("is the server running?")
		log.Fatal(err)
	}
	var v contracts.Version
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(data, &v)
	fmt.Println(v.Version)
}
