package flake

import (
	"fmt"
	"icebound.cc/flake/cloud/controller"
	"icebound.cc/flake/utils"
	"time"
)

type Demo2 struct {
	cloud *Aggcore
	edge1 *EdgeAggcore
	edge2 *EdgeAggcore

	clients []string
	ct      *controller.Controller
}

func NewDemo2(clients []string) (*Demo2, error) {
	ct, err := controller.NewController(false, "default")
	if err != nil {
		return nil, err
	}
	err = ct.ChangeDeviceTwins("learner0", "taskinfo", "ready")
	if err != nil {
		fmt.Println("err:", err.Error())
	} else {
		fmt.Println("learner0 Success")
	}
	cloud, err := NewCloudAggcore("cloudAggcore", "cloud/yaml/cloudaggcore.yaml", 10, 6666)
	if err != nil {
		return nil, err
	}
	edge1, err := NewEdgeAggcore("edge1Aggcore", "cloud/yaml/edge1aggcore.yaml", 2, 30, 2333, "flakecloud.com:6666", "edge1")
	if err != nil {
		return nil, err
	}
	edge2, err := NewEdgeAggcore("edge2Aggcore", "cloud/yaml/edge2aggcore.yaml", 2, 30, 2333, "flakecloud.com:6666", "edge2")
	if err != nil {
		return nil, err
	}
	return &Demo2{
		cloud:   cloud,
		edge1:   edge1,
		edge2:   edge2,
		clients: clients,
		ct:      ct,
	}, nil
}

func (dm *Demo2) ClientsChange(key string, value string) error {
	for _, x := range dm.clients {
		err := dm.ct.ChangeDeviceTwins(x, key, value)
		if err != nil {
			fmt.Println("err:", err.Error())
			return err
		} else {
			fmt.Println(x + " Success")
		}
	}
	return nil
}

func (dm *Demo2) ClientChange(client string, key string, value string) error {
	err := dm.ct.ChangeDeviceTwins(client, key, value)
	if err != nil {
		fmt.Println("err:", err.Error())
		return err
	} else {
		fmt.Println(client + " Success")
	}
	return nil
}

func (dm *Demo2) ClientReady() error {
	return dm.ClientsChange("taskinfo", "ready")
}

func (dm *Demo2) ClientRun() error {
	return dm.ClientsChange("taskinfo", "run")
}
func (dm *Demo2) ClientStop() error {
	return dm.ClientsChange("taskinfo", "stop")
}

func (dm *Demo2) InitApps() error {
	err := dm.cloud.DeployAggcore(dm.ct)
	if err != nil {
		return err
	}
	err = dm.WaitCloudStart(dm.cloud)
	if err != nil {
		return err
	}
	fmt.Println("Cloud start!")

	err = dm.edge1.DeployAggcore(dm.ct)
	if err != nil {
		return err
	}
	err = dm.WaitEdgeStart(dm.edge1)
	if err != nil {
		return err
	}
	fmt.Println("Edge1 start!")

	err = dm.edge2.DeployAggcore(dm.ct)
	if err != nil {
		return err
	}
	err = dm.WaitEdgeStart(dm.edge2)
	if err != nil {
		return err
	}
	fmt.Println("Edge2 start!")

	return nil
}

func (dm *Demo2) WatchApps() error {
	for true {
		time.Sleep(time.Duration(500) * time.Millisecond)
		res, err := dm.cloud.CheckState(dm.ct)
		if err != nil {
			return err
		}
		if res != utils.PodRunning {
			break
		}

		res, err = dm.edge1.CheckState(dm.ct)
		if err != nil {
			return err
		}
		if res != utils.PodRunning {
			break
		}

		res, err = dm.edge2.CheckState(dm.ct)
		if err != nil {
			return err
		}
		if res != utils.PodRunning {
			break
		}
	}
	err := dm.ct.Client.DeleteApp(dm.ct.GetNameSpace(), dm.cloud.podName)
	err = dm.ct.Client.DeleteApp(dm.ct.GetNameSpace(), dm.edge1.podName)
	err = dm.ct.Client.DeleteApp(dm.ct.GetNameSpace(), dm.edge2.podName)
	return err
}

func (dm *Demo2) StartDemo() {
	fmt.Println("Demo Start!!")
	fmt.Print("Now, start three apps...")
	err := dm.InitApps()
	if err != nil {
		fmt.Println("Demo2 StartDemo InitApps Error: %s", err.Error())
		return
	}
	fmt.Println("All apps started, now watching...")
	err = dm.WatchApps()
	if err != nil {
		fmt.Println("Demo2 StartDemo WatchApps Error: %s", err.Error())
		return
	}
	fmt.Println("Finished!")
}

func (dm *Demo2) WaitCloudStart(aggcore *Aggcore) error {
	for true {
		time.Sleep(time.Duration(1) * time.Second)
		res, err := aggcore.CheckState(dm.ct)
		if err != nil {
			return err
		}
		if res == utils.PodRunning {
			return nil
		}
	}
	return nil
}

func (dm *Demo2) WaitEdgeStart(aggcore *EdgeAggcore) error {
	for true {
		time.Sleep(time.Duration(1) * time.Second)
		res, err := aggcore.CheckState(dm.ct)
		if err != nil {
			return err
		}
		if res == utils.PodRunning {
			return nil
		}
	}
	return nil
}
