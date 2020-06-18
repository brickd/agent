package client

import (
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Gateway struct {
	*logrus.Entry

	client *Client
}

func NewGateway(l *logrus.Entry, client *Client) *Gateway {
	return &Gateway{
		Entry:  l,
		client: client,
	}
}

func (g *Gateway) Init() error {
	return nil
}

func (g *Gateway) RunHTTP(ctx context.Context) error {
	if err := g.Init(); err != nil {
		return err
	}

	http.HandleFunc("/push", func(wr http.ResponseWriter, r *http.Request) {
		bb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			g.Error("Error reading body from the http push", err)
			return
		}

		err = g.client.PublishAsGateway(bb)
		if err != nil {
			g.Error("An error occured when publishing message:", err)
		}
	})

	g.Info("Starting the HTTP Server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}

	return nil
}
