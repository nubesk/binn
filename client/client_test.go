package client

import (
	"time"
	"context"
	"testing"
	"net/http/httptest"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/binn/binn"
	"github.com/binn/server"
)


func TestPostBottle(t *testing.T) {
	cfg := binn.DefaultConfig()
	cfg.SetDeliveryCycle(time.Duration(1) * time.Millisecond)

	idStorage := binn.DefaultIDStorage()
	storage := binn.NewContainerStorage(true, time.Duration(10)*time.Minute, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10)*time.Minute),
	)

	engine := binn.NewEngine(cfg, storage)

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()
	
	ts := httptest.NewServer(server.BottleHandlerFunc(engine))

	cli := NewClient(ts.URL, 0)
	err := cli.Post("1c7a8201-cdf7-11ec-a9b3-0242ac110004", "Post a Bottle")

	assert.Nil(t, err)
}

func TestGetBottle(t *testing.T) {
	cfg := binn.DefaultConfig()
	cfg.SetDeliveryCycle(time.Duration(1) * time.Millisecond)

	idStorage := binn.DefaultIDStorage()
	storage := binn.NewContainerStorage(true, time.Duration(10)*time.Minute, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10)*time.Minute),
	)
	storage.Add(binn.NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"",
		nil,
	))
	
	engine := binn.NewEngine(cfg, storage)
	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()
	
	ts := httptest.NewServer(server.BottleHandlerFunc(engine))
	defer ts.Close()
	
	cli := NewClient(ts.URL, time.Duration(100 * time.Millisecond))
	b, _ := cli.Get()

	assert.Equal(t, "", b.Message().Text)
	assert.NotEqual(t, "1c7a8201-cdf7-11ec-a9b3-0242ac110004", b.ID())
}
