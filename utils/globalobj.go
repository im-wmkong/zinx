package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/zicafe"
)

type GlobalObj struct {
	TcpServer         zicafe.IServer
	Host              string
	TcpPort           int
	Name              string
	Version           string
	MaxConn           int
	MaxPackageSize    uint32
	WorkerPoolSize    uint32
	MaxWorkerPoolSize uint32
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Host:              "0.0.0.0",
		TcpPort:           8999,
		Name:              "ZinxServerApp",
		Version:           "V0.4",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,
		MaxWorkerPoolSize: 1024,
	}

	GlobalObject.Reload()
}

func (o *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
