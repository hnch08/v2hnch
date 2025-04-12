package v2ray

import (
	"bytes"
	"embed"
	"fmt"
	"sync"

	core "github.com/v2fly/v2ray-core/v5"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/outbound"
	_ "github.com/v2fly/v2ray-core/v5/main/formats"
)

var (
	mu      sync.Mutex
	_server core.Server
)

//go:embed v2ray_config.json
var conf embed.FS

func StartV2Ray(data []byte) (core.Server, error) {
	configFormat := "json"
	reader := bytes.NewReader(data)
	config, err := core.LoadConfig(configFormat, reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// coreConfig, err := cf.Build()
	// if err != nil {
	// 	return nil, err
	// }

	server, err := core.New(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return server, nil
}

func Start() error {
	configure, err := conf.ReadFile("v2ray_config.json")
	// fmt.Println(string(configure))
	if err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()

	if _server != nil {
		_server.Close()
		_server = nil
	}

	server, err := StartV2Ray(configure)
	if err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}

	_server = server
	return nil
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if _server != nil {
		_server.Close()
		_server = nil
	}
}
