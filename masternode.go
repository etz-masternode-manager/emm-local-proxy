package main

import (
	"errors"
	"io/ioutil"
	"os/exec"
)

const homePath = "/home/emm"
const masternodesPath = homePath + "/masternodes"

type MasternodeRequest struct {
	Command  string `json:"command"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	NodePort int    `json:"nodePort"`
	NodeKey  string `json:"nodeKey"`
}

func MasternodeRequestProcess(request MasternodeRequest) (string, error) {
	switch request.Command {
	case "start":
		return "", masternodeStart(request.Name)
	case "stop":
		return "", masternodeStop(request.Name)
	case "create":
		return "", masternodeCreate(request.Name, request.Port, request.NodePort, request.NodeKey)
	case "delete":
		return "", masternodeDelete(request.Name)
	case "nodeKey":
		return masternodeNodeKey(request.Name)
	}
	return "", errors.New("Unknown command: " + request.Name)
}

func masternodeStart(name string) error {
	masternodeStop(name)
	//cmd := exec.Command("nohup", "./start.sh >/dev/null 2>&1&")
	cmd := exec.Command(masternodesPath + "/" + name + "/bin/geth")
	cmd.Dir = masternodesPath + "/" + name
	cmd.Start()
	return nil
}

func masternodeStop(name string) error {
	return nil
}

func masternodeCreate(name string, port int, nodePort int, nodeKey string) error {
	return nil
}

func masternodeDelete(name string) error {
	return nil
}

func masternodeUpdate(name string) error {
	hostGethUpdate()

	return nil
}

func masternodeNodeKey(name string) (string, error) {
	dat, err := ioutil.ReadFile(masternodesPath + name + "/data/geth/nodekey")
	if err != nil {
		return "", errors.New("can't read node key")
	}

	return string(dat), nil
}

func hostGethUpdate() {

}
