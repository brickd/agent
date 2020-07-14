package httpgateway

import (
	"context"
	"github.com/brickd/agent/internal/brickd"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Gateway struct {
	*logrus.Entry

	client brickd.Conn
}

func New(l *logrus.Entry, client brickd.Conn) *Gateway {
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

		err = g.client.Publish(bb)
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
