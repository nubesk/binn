package server

import (
	"io"
	"time"
	"bytes"
	"context"
	"testing"
	"encoding/json"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/binn/binn"
)

func TestHandleGetBottle (t *testing.T) {
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
	
	engine := binn.NewEngine(
		cfg,
		storage,
	)

	handler := BottleGetHandlerFunc(engine, 10)

	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(100 * time.Millisecond))
	engine.Run(ctx)
	defer cancelFunc()

	req := httptest.NewRequest("GET", "http://example.com/api/bottle", nil)
	reqCtx, _ := context.WithTimeout(
		context.Background(),
		time.Duration(10 * time.Millisecond))
	req = req.WithContext(reqCtx)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	var rb responseBottle
	if err := json.Unmarshal(body[20:], &rb); err != nil {
		assert.Failf(t, "failed", "%w", err)
		return
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/event-stream; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.NotEqual(t, "1c7a8201-cdf7-11ec-a9b3-0242ac110004", rb.ID)
	assert.Equal(t, "", rb.Message.Text)
}

func TestHandlePostBottle (t *testing.T) {
	cfg := binn.DefaultConfig()
	cfg.EnableDebug()

	idStorage := binn.DefaultIDStorage()
	storage := binn.NewContainerStorage(true, time.Duration(10)*time.Minute, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10)*time.Minute),
	)

	engine := binn.NewEngine(
		cfg,
		storage,
	)

	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(100 * time.Millisecond))
	engine.Run(ctx)
	defer cancelFunc()

	handler := BottlePostHandlerFunc(engine)

	reqBody := bytes.NewBufferString("{\"id\":\"1c7a8201-cdf7-11ec-a9b3-0242ac110004\",\"message\":{\"text\":\"Post a Bottle\"}}")
	req := httptest.NewRequest("POST", "http://example.com/api/bottle", reqBody)
	w := httptest.NewRecorder()
	handler(w, req)

	// wait for adding a bottle to storage
	time.Sleep(time.Duration(10) * time.Millisecond)

	resp := w.Result()
	gottenBottle, _ := storage.Get()

	assert.Equal(t, 204, resp.StatusCode)
	assert.Equal(t, "Post a Bottle", gottenBottle.Message().Text)
}
