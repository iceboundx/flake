package main

import (
	"fmt"
	"icebound.cc/flake/cloud/flakehttp"
	"strings"
	"time"
)

const (
	GetCmd    = "get"
	SetCmd    = "set"
	DemoStart = "demo"
)

var wk *flakehttp.Worker

func DealSetCmd(src string, spec string) {
	if src == "demo" {
		demo(spec)
		return
	}
	specs := strings.Split(spec, "=")
	if len(specs) != 2 {
		fmt.Println("spec error")
		return
	}
	err := wk.Ct.ChangeDeviceTwins(src, specs[0], specs[1])
	if err != nil {
		fmt.Println("err:", err.Error())
	} else {
		fmt.Println("Success")
	}
}

func DealGetCmd(src string, spec string) {

	if src == "device" {
		if spec == "list" {
			de, err := wk.Ct.GetDeviceList()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				for _, x := range de {
					fmt.Println("Device:" + x.Name)
					fmt.Println("Description:" + x.Desc)
					fmt.Println("Status:" + x.State)
					fmt.Println("")
				}
			}
		} else {
			x, err := wk.Ct.GetDevice(spec)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Device:" + x.Name)
				fmt.Println("Description:" + x.Desc)
				fmt.Println("Status:" + x.State)
				fmt.Printf("LastCommTime:%v\n", x.TimeStamp)
				fmt.Println("EdgeNode:" + x.EdgeNode)
				fmt.Println("")
			}
		}
	} else if src == "node" {
		nd, err := wk.Ct.GetNodeList()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			for _, x := range nd {
				fmt.Println("Device:" + x.Name)
				fmt.Println("Description:" + x.Desc)
				fmt.Println("Status:" + x.State)
				fmt.Println("NodeType:" + x.NodeVersion)
				fmt.Println("")
			}
		}
	}

}
func demo(cmd string) {
	err := wk.Ct.ChangeDeviceTwins("learner0", "taskinfo", cmd)
	if err != nil {
		fmt.Println("err:", err.Error())
	} else {
		fmt.Println("learner0 Success")
	}
	err = wk.Ct.ChangeDeviceTwins("learner1", "taskinfo", cmd)
	if err != nil {
		fmt.Println("err:", err.Error())
	} else {
		fmt.Println("learner1 Success")
	}
}

func main() {
	wk = flakehttp.GetWorkerInstance("0.0.0.0:443", "default")
	go wk.RunServer()
	for {
		time.Sleep(200 * 1000)
	}

	/*for {
		fmt.Println("Please Input Commands:")
		cmd := ""
		src := ""
		spec := ""
		_, err := fmt.Scan(&cmd, &src, &spec)
		if err != nil {
			fmt.Println("Input Error:", err.Error())
		}
		switch cmd {
		case SetCmd:
			DealSetCmd(src, spec)
		case GetCmd:
			DealGetCmd(src, spec)
		}
	}*/

}
