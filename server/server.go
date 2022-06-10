package server

import (
	"log"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/binn/binn"
)

var Logger = log.Default()


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
	if Logger == nil {
		log.Printf(format, v...)
		return
	}
	Logger.Printf(format, v...)
}

func BottleGetHandlerFunc(engine *binn.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		outCh := engine.GetOutChan()

		var c binn.Container;
	Loop:
		for {
			select {
			case <- r.Context().Done():
				return
			case c = <-outCh:
				break Loop
			default:
				break
			}
		}

		res := containerToResponse(c)

		if byte_, err := json.Marshal(res); err == nil {
			if _, err := w.Write(byte_); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logf("[%d] %s", http.StatusInternalServerError, "server error")
				return
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logf("[%d] %s", http.StatusInternalServerError, "server error")
			return
		}

		logf(
			"[%d] %s", http.StatusOK,
			fmt.Sprintf("id: %#v message: %#v", c.ID(), c.Message().Text))
		return
	}
}

func BottlePostHandlerFunc(engine *binn.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		var req requestBottle
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logf("[%d] %s", http.StatusBadRequest, fmt.Sprintf("payload is invalid format, %s", string(body)))
			return
		}

		c := requestToContainer(&req)

		inCh := engine.GetInChan()
		inCh <- c

		w.WriteHeader(http.StatusNoContent)

		logf(
			"[%d] %s", http.StatusNoContent,
			fmt.Sprintf("id: %#v message: %#v", c.ID(), c.Message().Text),
		)
	}
}

func BottleHandlerFunc(engine *binn.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST")

		var handler http.HandlerFunc
		if (r.Method == http.MethodGet) {
			handler = BottleGetHandlerFunc(engine)
		} else if (r.Method == http.MethodPost) {
			handler = BottlePostHandlerFunc(engine)
		}

		handler(w, r)
	}
}


func Server(engine *binn.Engine, addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/bottle", BottleHandlerFunc(engine))

	return &http.Server{
		Addr: addr,
		Handler: mux,
	}
}
