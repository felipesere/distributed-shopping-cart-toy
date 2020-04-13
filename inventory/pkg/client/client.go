package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/felipesere/inventory/v0/pkg/products"
	"io/ioutil"
	"net/http"
	"time"
)

type client struct {
	peers   []string
	self    string
	timeout time.Duration
}

type RemoteInventory interface {
	Available(category string) ([]products.RemoteProduct, error)
}

func New(self string, peers []string, timeout time.Duration) RemoteInventory {
	return &client {
		self: self,
		peers: peers,
		timeout: timeout,
	}
}

func (c *client) Available(category string) ([]products.RemoteProduct, error) {
	var remoteProducts []products.RemoteProduct

	for _, peer := range c.peers {
		ps, err := c.queryPeer(peer, category)
		if err != nil {
			return remoteProducts, fmt.Errorf("unable to get products from peer %q: %w", peer, err)
		}
		remoteProducts = append(remoteProducts, ps...)
	}

	return remoteProducts, nil
}

func (c *client) queryPeer(peer, category string) ([]products.RemoteProduct, error) {
	var ps []products.RemoteProduct
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/inventory/available", peer), nil)
	client := http.Client{}
	response, err := client.Do(req)
	cancel()
	if err != nil {
		return  ps, err
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return  ps, fmt.Errorf("unable to read body: %w", err)
	}

	var receivedProducts []products.Product
	_ = json.Unmarshal(content, &receivedProducts)

	for _, product := range receivedProducts {
		ps = append(ps, products.RemoteProduct{Product: product})
	}

	return ps, nil
}
