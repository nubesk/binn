package main

import (
	"os"
	"fmt"
	"time"
	"context"
	"strconv"

	"github.com/binn/server"
	"github.com/binn/binn"
)

func printEngineConfig(cfg *binn.Config) {
	fmt.Printf("%18s: %d\n", "Seed", cfg.Seed())
	fmt.Printf("%18s: %f\n", "Delivery cycle sec", cfg.DeliveryCycle().Seconds())
	fmt.Printf("%18s: %t\n", "Enable validation", cfg.Validation())
	fmt.Printf("%18s: %f\n", "Generate cycle sec", cfg.GenerateCycle().Seconds())
	fmt.Printf("%18s: %t\n", "Enable debug", cfg.Debug())
}

func printServerConfig(cfg *server.Config) {
	fmt.Printf("%18s: %d\n", "Send empty sec", cfg.SendEmptySec())
}

func loadEnvAsInt(key string, defaultValue int) int {
	var v int;
	if s := os.Getenv(key); s != "" {
		v, _ = strconv.Atoi(s)
	} else {
		v = defaultValue
	}
	return v
}

func loadEnvAsBool(key string, defaultValue bool) bool {
	var v bool;
	if s := os.Getenv(key); s != "" {
		v, _ = strconv.ParseBool(s)
	} else {
		v = defaultValue
	}
	return v
}

func loadEngineConfigFromEnv() *binn.Config {
	seed := loadEnvAsInt("BINN_SEED", 42)
	deliveryCycleSec := loadEnvAsInt("BINN_DELIVERY_CYCLE_SEC", 20)
	enableValidation := loadEnvAsBool("BINN_ENABLE_VALIDATION", true)
	generateCycleSec := loadEnvAsInt("BINN_GENERATE_CYCLE_SEC", 20)
	enableDebug := loadEnvAsBool("BINN_ENABLE_DEBUG", true)
	
	return binn.NewConfig(seed, time.Duration(deliveryCycleSec) * time.Second, enableValidation,
		time.Duration(generateCycleSec) * time.Second, enableDebug)
}

func loadServerConfigFromEnv() *server.Config {
	sendEmptySec := loadEnvAsInt("BINN_SEND_EMPTY_SEC", 29)
	return server.NewConfig(sendEmptySec)
}

func main() {
	ecfg := loadEngineConfigFromEnv()
	scfg := loadServerConfigFromEnv()

	idStorage := binn.DefaultIDStorage()
	storage := binn.NewContainerStorage(true, time.Duration(10)*time.Minute, idStorage)

	engine := binn.NewEngine(
		ecfg,
		storage,
	)

	engine.SetGenerateContainerHandler(func(cs binn.ContainerKeeper) error {
		id := binn.GenerateID()
		err := idStorage.Add(id, time.Now().Add(time.Duration(10)*time.Minute))
		if err != nil {
			return err
		}
		err = cs.Add(binn.NewBottle(id, "", nil))
		if err != nil {
			return err
		}
		return nil
	})

	ctx, cancelFunc := context.WithCancel(context.Background())
	engine.Run(ctx)
	defer cancelFunc()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	printEngineConfig(ecfg)
	printServerConfig(scfg)

	srv := server.Server(engine, fmt.Sprintf(":%s", port), scfg)
	srv.ListenAndServe()
	defer srv.Close()
}
