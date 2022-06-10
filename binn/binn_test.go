package binn

import (
	"time"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	_ = NewEngine(
		DefaultConfig(),
		NewContainerStorage(false, 0, nil),
	)
}

func TestDefaultEngine(t *testing.T) {
	engine := DefaultEngine()
	cfg := engine.GetConfig()

	assert.Equal(t, 42, cfg.Seed())
	assert.Equal(t, 15.0, cfg.DeliveryCycle().Minutes())
	assert.Equal(t, 15.0, cfg.GenerateCycle().Minutes())
}

func TestAddBottle(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetDeliveryCycle(time.Duration(1) * time.Millisecond)
	
	engine := NewEngine(
		cfg,
		DefaultContainerStorage(),
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()

	bottle := NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		nil,
	)

	ch := engine.GetInChan()
	ch <- bottle
}

func TestGetBottle(t *testing.T) {
	idStorage := DefaultIDStorage()
	storage := NewContainerStorage(true, 0, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10) * time.Minute),
	)
	_ = storage.Add(NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		nil,
	))

	cfg := DefaultConfig()
	cfg.SetDeliveryCycle(time.Duration(1) * time.Millisecond)

	engine := NewEngine(
		cfg,
		storage,
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()

	outCh := engine.GetOutChan()
	bottle := <- outCh

	assert.Equal(t, "This is a Test Message", bottle.Message().Text)
	assert.NotEqual(t, "1c7a8201-cdf7-11ec-a9b3-0242ac110004", bottle.ID())
}

func TestBottleGetDeley(t *testing.T) {
	idStorage := DefaultIDStorage()
	storage := NewContainerStorage(true, 0, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10) * time.Minute),
	)
	storage.Add(NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		nil,
	))

	cfg := DefaultConfig()
	cfg.SetDeliveryCycle(time.Duration(15) * time.Millisecond)

	engine := NewEngine(
		cfg,
		storage,
	)

	ctx, cancelFunc := context.WithCancel(context.Background())
	outCh := engine.GetOutChan()
	begin := time.Now()
	engine.Run(ctx)
	defer cancelFunc()
	<- outCh
	end := time.Now()
	elapsed := end.Sub(begin)

	assert.GreaterOrEqual(t, 15 * time.Millisecond, elapsed.Milliseconds())
}

func TestGetEmptyBottleIsGenerated(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SetDeliveryCycle(time.Duration(10) * time.Millisecond)
	cfg.SetGenerateCycle(time.Duration(1) * time.Millisecond)

	idStorage := DefaultIDStorage()
	engine := NewEngine(
		cfg,
		NewContainerStorage(true, time.Duration(10) * time.Minute, idStorage),
	)
	engine.SetGenerateContainerHandler(func(cs ContainerKeeper) error {
		id := GenerateID()
		err := idStorage.Add(id, time.Now().Add(time.Duration(10) * time.Minute))
		if err != nil {
			return err
		}
		err = cs.Add(NewBottle(id, "", nil))
		if err != nil {
			return err
		}
		return err
	})
	
	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(100 * time.Millisecond))
	engine.Run(ctx)
	defer cancelFunc()

	outCh := engine.GetOutChan()
	bottle := <- outCh

	assert.Equal(t, "", bottle.Message().Text)
}
