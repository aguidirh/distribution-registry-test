package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/distribution/distribution/v3/configuration"
	"github.com/distribution/distribution/v3/registry"

	_ "github.com/distribution/distribution/v3/registry/storage/driver/filesystem"
)

func main() {
	registry, err := setupRegistry()
	if err != nil {
		println("error on setupRegistry")
		println(err)
	}

	// run registry server
	var errchan chan error
	go func() {
		errchan <- registry.ListenAndServe()
	}()
	select {
	case err = <-errchan:
		println("error on gofunc")
		println(err)
	default:
	}

	time.Sleep(5 * time.Second)

	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		println("error on dial")
	}
	fmt.Fprintf(conn, "GET /v2/ ")

}

func setupRegistry() (*registry.Registry, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	config.Storage = map[string]configuration.Parameters{"filesystem": map[string]interface{}{}}

	return registry.NewRegistry(context.Background(), config)
}

func getConfig() (*configuration.Configuration, error) {
	configurationPath := "./config-example.yaml"

	fp, err := os.Open(configurationPath)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	config, err := configuration.Parse(fp)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", configurationPath, err)
	}

	return config, nil
}
