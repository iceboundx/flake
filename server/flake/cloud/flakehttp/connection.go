package flakehttp

import (
	"encoding/json"
	"fmt"
	"icebound.cc/flake/cloud/controller"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

const DefaultLearnerDeviceModelName = "defaultlearnermodel"

type Worker struct {
	Ct        *controller.Controller
	ipport    string
	nameSpace string
}

var instance *Worker
var once sync.Once

func GetWorkerInstance(ipport string, nameSpace string) *Worker {
	once.Do(func() {
		ct, err := controller.NewController(false, nameSpace)
		if err != nil {
			panic(err)
			return
		}
		nowWk := &Worker{
			Ct:        ct,
			ipport:    ipport,
			nameSpace: nameSpace,
		}
		instance = nowWk
		http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				fmt.Println(err.Error())
			}
			deviceId := r.PostForm.Get("id")
			edgeNode := r.PostForm.Get("node")
			key := r.PostForm.Get("key")

			if deviceId == "" || edgeNode == "" || key == "" {
				w.Write([]byte("Param Error!"))
				fmt.Println(deviceId, edgeNode, key)
				return
			}
			res, err := nowWk.Authenticate(deviceId, key)
			if res == false {
				w.Write([]byte("Key Error!"))
				fmt.Println(deviceId, edgeNode, key)
				return
			}
			err = nowWk.NewLearnerDeviceConnection(deviceId, edgeNode)
			if err == nil || strings.Contains(err.Error(), "already exists") {
				fmt.Println("Device Connected")
				w.Write([]byte("Connect Ok"))
			} else {
				w.Write([]byte("err: " + err.Error()))
			}
		})
	})
	return instance
}

func (wk *Worker) NewLearnerDeviceConnection(DeviceId string, EdgeNode string) error {
	lst, err := wk.Ct.GetDeviceModelList()
	if err != nil {
		return err
	}
	flag := false
	for _, y := range lst {
		if y.Name == DefaultLearnerDeviceModelName {
			flag = true
			break
		}
	}
	if !flag {
		dm := controller.NewLearnerDeviceModel(DefaultLearnerDeviceModelName)
		err := wk.Ct.AddDeviceModel(dm)
		if err != nil {
			return err
		}
	}
	dv := controller.NewLearnerDeviceInstance(DeviceId, DeviceId,
		DefaultLearnerDeviceModelName, EdgeNode)
	return wk.Ct.AddDevice(dv)
}

func (wk *Worker) RunServer() {
	http.ListenAndServeTLS(wk.ipport, "server.crt", "server.key", nil)
}

func (wk *Worker) Authenticate(id string, key string) (bool, error) {
	f, err := os.OpenFile("keyfile.json", os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	contentByte, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		fmt.Println(err2.Error())
		return false, err
	}
	var inter interface{}
	err = json.Unmarshal(contentByte, &inter)
	if err != nil {
		fmt.Println("error in translating,", err.Error())
		return false, err
	}
	keys, ok := inter.(map[string]interface{})
	if ok {
		res, ok2 := keys[id]
		if ok2 {
			return res == key, nil
		}
	}
	return false, nil
}
