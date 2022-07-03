package binn

import (
	"os"
	"log"
	"fmt"
	"time"
	"context"
)

var Logger = log.New(os.Stderr, "[ENGINE] ", log.LstdFlags)

type GenerateContainerHandlerFunc func (cs ContainerKeeper) error

type Engine struct {
	cfg      *Config
	storage  ContainerKeeper
	inCh     chan Container
	outCh    chan Container
	generateContainerHandler GenerateContainerHandlerFunc
}

func NewEngine(cfg *Config, storage ContainerKeeper) *Engine {
	return &Engine{
		cfg:     cfg,
		storage: storage,
		inCh:    make(chan Container, 1),
		outCh:   make(chan Container, 1),
		generateContainerHandler: DefaultGenerateContainerHandlerFunc(),
	}
}

func DefaultEngine() *Engine {
	return NewEngine(
		DefaultConfig(),
		DefaultContainerStorage(),
	)
}

func (e *Engine) logf(format string, v ...interface{}) {
	if !e.GetConfig().Debug() {
		return
	}
	
	Logger.Printf(format, v...)
}

func (e *Engine) GetConfig() *Config {
	return e.cfg
}

func (e *Engine) GetInChan() chan Container {
	return e.inCh
}

func (e *Engine) GetOutChan() chan Container {
	return e.outCh
}

func (e *Engine) SetGenerateContainerHandler(h GenerateContainerHandlerFunc) {
	e.generateContainerHandler = h
}

func DefaultGenerateContainerHandlerFunc() GenerateContainerHandlerFunc {
	return func(cs ContainerKeeper) error {
		b := NewBottle("", GenerateID(), nil)
		err := cs.Add(b)
		if err != nil {
			return err
		}
		return nil
	}
}

func (e *Engine) Run(ctx context.Context) {
	go func() {
		if !e.cfg.Validation() {
			return
		}
		
		t := time.NewTicker(e.cfg.GenerateCycle())
		defer t.Stop()
	Loop:
		for{
			select {
			case <- ctx.Done():
				break Loop
			case <- t.C:
				err := e.generateContainerHandler(e.storage)
				if err != nil {
					e.logf("%s", err)
					break
				}
				e.logf("The engine generate a empty container")
			}
		}
	}()

	go func() {
	Loop:
		for {
			select {
			case <- ctx.Done():
				break Loop
			case c := <- e.inCh:
				// ignore a error intentionally
				// it is not necessary for a user to tell a error
				// it hides whether server received bottle or not
				if err := e.storage.Add(c); err == nil {
					e.logf(fmt.Sprintf("The engine add a container(id=%#v message=%#v)",
						c.ID(),
						c.Message().Text,
					))
				} else {
					e.logf("Failed: %s", err)
				}
			default:
				break
			}
		}
	}()

	go func() {
		t := time.NewTicker(e.cfg.DeliveryCycle())
		defer t.Stop()

	Loop:
		for {
			select {
			case <- ctx.Done():
				break Loop
			case <- t.C:
				c, err := e.storage.Get()
				if err != nil {
					break
				}

				e.logf(fmt.Sprintf("The engine delivery a container(id=%#v message=%#v)",
					c.ID(),
					c.Message().Text,
				))
				e.outCh <- c
			default:
				break
			}
		}
	}()
}
