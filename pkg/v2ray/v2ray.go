package v2ray

import (
	"bytes"
	"embed"
	"sync"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/infra/conf/serial"
	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"
)

var (
	mu      sync.Mutex
	_server core.Server
)

//go:embed v2ray_config.json
var conf embed.FS

func StartV2Ray(data []byte) (core.Server, error) {
	cf := &v4.Config{}
	c, err := serial.DecodeJSONConfig(bytes.NewReader(data))
	common.Must(err)
	*cf = *c

	coreConfig, err := cf.Build()
	if err != nil {
		return nil, err
	}

	server, err := core.New(coreConfig)
	if err != nil {
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
