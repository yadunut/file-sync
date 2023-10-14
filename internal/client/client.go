package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yadunut/file-sync/internal/contracts"
	"github.com/yadunut/file-sync/internal/util"
	"go.uber.org/zap"
)

type Client struct {
	Log    *zap.SugaredLogger
	Config util.Config
}

func NewClient(log *zap.SugaredLogger, config util.Config) *Client {
	return &Client{Log: log, Config: config}
}

func (c *Client) get(url string, into any) error {
	res, err := http.Get(fmt.Sprintf("http://%s/%s", c.Config.GetUrl(), url))
	if err != nil {
		c.Log.Fatal("is the server running?", err)
		c.Log.Fatal(err)
	}
	return json.NewDecoder(res.Body).Decode(into)
}

func (c *Client) post(url string, req any, res any) error {
	responseBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/%s", c.Config.GetUrl(), url), "application/json", bytes.NewBuffer(responseBody))
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(&res)
}

func (c *Client) Version() (contracts.VersionRes, error) {
	var v contracts.VersionRes
	c.get("version", &v)
	return v, nil
}

func (c *Client) WatchUp(req contracts.WatchUpReq) (res contracts.WatchUpRes, err error) {
	err = c.post("watch/up", req, &res)
	if err != nil {
		return contracts.WatchUpRes{Success: false}, err
	}
	return res, nil
}

func (c *Client) WatchDown(req contracts.WatchDownReq) (res contracts.WatchDownRes, err error) {
	err = c.post("watch/down", req, &res)
	if err != nil {
		return contracts.WatchDownRes{Success: false}, err
	}
	return res, nil
}

func (c *Client) WatchList() (res contracts.WatchListRes, err error) {
	err = c.get("watch", &res)
	if err != nil {
		return contracts.WatchListRes{Success: false}, err
	}
	return res, nil
}
