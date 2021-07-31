package main

import (
	"fmt"
	"icebound.cc/flake/cloud/flake"
	"time"
)

func main() {
	clients := []string{
		"learner0",
		"learner1",
		"learner2",
		"learner3",
	}
	dm2, err := flake.NewDemo2(clients)

	if err != nil {
		panic(err)
		return
	}
	var inp string

	fmt.Println("Press any key to let apps ready!")
	fmt.Scanln(&inp)

	err = dm2.InitApps()
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("Press any key to let device ready!")
	fmt.Scanln(&inp)

	dm2.ClientChange("learner0", "taskparam", "python_nn.py_flakeedge1.com:2333_edge1")
	dm2.ClientChange("learner1", "taskparam", "python_nn.py_flakeedge1.com:2333_edge1")

	dm2.ClientChange("learner2", "taskparam", "python_nn.py_flakeedge2.com:2333_edge2")
	dm2.ClientChange("learner3", "taskparam", "python_nn.py_flakeedge2.com:2333_edge2")

	fmt.Println("Change Learners to Ready")
	dm2.ClientReady()

	fmt.Println("Press any key to run!")
	fmt.Scanln(&inp)
	fmt.Println("Run!")
	dm2.ClientRun()
	fmt.Println("Press any key to stop!")
	fmt.Scanln(&inp)
	dm2.ClientStop()

	fmt.Println("Press any key to run again!")
	fmt.Scanln(&inp)
	fmt.Println("Run!")
	dm2.ClientRun()

	for {
		time.Sleep(200 * 1000)
	}
}

//ct, err := controller.NewController(false, "default")
//	if err != nil {
//	panic(err.Error())
//}
/*err = ct.AddDeviceModel(controller.DeviceModel{
	ResMeta: controller.ResMeta{
		Name:  "cttest",
		State: "",
		Desc:  "",
	},
	Properties: []utils.Property{
		{Name: utils.DeviceStatusName, Desc: "status filed", DefaultVal: "Ready"},
		{Name: "modelfiled", Desc: "test filed", DefaultVal: "nonono"},
	},
})
if err != nil {
	panic(err.Error())
}
time.Sleep(10000)
deviceMList, err = ct.GetDeviceModelList()
if err != nil {
	panic(err.Error())
}
fmt.Printf("%+v\n", deviceMList)
err = ct.AddDevice(controller.Device{
	ResMeta: controller.ResMeta{
		Name:  "testname",
		State: "Ready",
		Desc:  "Testtt",
	},
	ModelName:   "cttest",
	EdgeNode:    "raspberrypi",
	DesireState: "Ready",
	TimeStamp:   0,
})
if err != nil {
	panic(err.Error())
}*/
