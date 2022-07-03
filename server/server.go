package server

import (
	"os"
	"log"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/binn/binn"
)


var Logger = log.New(os.Stderr, "[SERVER] ", log.LstdFlags)

type Config struct{
	sendEmptySec int
}

type responseMessage struct {
	Text string `json:"text"`
}

type requestMessage struct {
	Text string `json:"text"`
}

type requestBottle struct {
	ID        string           `json:"id"`
	Message   *responseMessage `json:"message"`
	ExpiredAt *time.Time       `json:"expired_at"`
}

type responseBottle struct {
	ID        string           `json:"id"`
	Message   *responseMessage `json:"message"`
	ExpiredAt *time.Time       `json:"expired_at"`
}

func NewConfig(sendEmptySec int) *Config {
	return &Config{ sendEmptySec: sendEmptySec }
}

func (c *Config) SendEmptySec() int {
	return c.sendEmptySec
}

func containerToResponse(c binn.Container) *responseBottle {
	return &responseBottle{
		ID:        c.ID(),
		Message:   &responseMessage{
			Text: c.Message().Text,
		},
		ExpiredAt: c.ExpiredAt(),
	}
}

func requestToContainer(req *requestBottle) binn.Container {
	return binn.NewBottle(req.ID, req.Message.Text, req.ExpiredAt)
}

func logf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

func BottleGetHandlerFunc(engine *binn.Engine, sendEmptySec int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")

		ticker := time.NewTicker(time.Duration(sendEmptySec) * time.Second)

		outCh := engine.GetOutChan()

		flusher, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var c binn.Container;
	Loop:
		for {
			select {
			case <- r.Context().Done():
				break Loop
			case c = <-outCh:
				res := containerToResponse(c)
				if byte_, err := json.Marshal(res); err == nil {
					byte_ = append(byte_, 10)
					byte_ = append(byte_, 10)
					if _, err := w.Write(byte_); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logf("%d %s", http.StatusInternalServerError, "Failed to write response")
						return
					}
					logf("send a container(id=%#v message=%#v)", c.ID(), c.Message().Text)
					flusher.Flush()
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					logf("%d %s", http.StatusInternalServerError, "Failed to decode response")
					return
				}
			case _ = <-ticker.C:
				if _, err := w.Write([]byte{10, 10}); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logf("%d %s", http.StatusInternalServerError, "Failed to write empty lines")
					return
				}
				flusher.Flush()
			default:
				break
			}
		}
		return
	}
}

func BottlePostHandlerFunc(engine *binn.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var req requestBottle
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logf("%d %s", http.StatusBadRequest, fmt.Sprintf("payload is invalid format, %s", string(body)))
			return
		}

		c := requestToContainer(&req)

		inCh := engine.GetInChan()
		inCh <- c

		w.WriteHeader(http.StatusNoContent)

		logf(
			"%d %s", http.StatusNoContent,
			fmt.Sprintf("id: %#v message: %#v", c.ID(), c.Message().Text),
		)
	}
}

func BottleHandlerFunc(engine *binn.Engine, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST")

		var handler http.HandlerFunc
		if (r.Method == http.MethodGet) {
			handler = BottleGetHandlerFunc(engine, cfg.sendEmptySec)
		} else if (r.Method == http.MethodPost) {
			handler = BottlePostHandlerFunc(engine)
		}

		handler(w, r)
	}
}


func Server(engine *binn.Engine, addr string, cfg *Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/bottle", BottleHandlerFunc(engine, cfg))

	return &http.Server{
		Addr: addr,
		Handler: mux,
	}
}
