package controller

import (
	"fmt"
	devices "github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha2"
	"icebound.cc/flake/utils"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type ResMeta struct {
	Name  string
	State string
	Desc  string
}

type Device struct {
	ResMeta
	ModelName   string
	EdgeNode    string
	DesireState string
	TimeStamp   int64
}

type Node struct { //amm core
	ResMeta
	NodeVersion string
}

type DeviceModel struct {
	ResMeta
	Properties []utils.Property
}

type Pod struct {
	ResMeta
	Status v1.PodStatus
}

type Controller struct {
	Client    *utils.K8sClient
	nameSpace string
}

func NewController(isInCluster bool, nameSpace string) (*Controller, error) {
	cl, err := utils.NewK8sClient(isInCluster)
	if err != nil {
		fmt.Printf("Controller Created Failed:%s \n", err.Error())
		return nil, err
	}
	return &Controller{
		Client:    cl,
		nameSpace: nameSpace,
	}, nil
}

func (ct *Controller) GetNameSpace() string {
	return ct.nameSpace
}

func (ct *Controller) GetDeviceList() ([]Device, error) {
	ret, err := ct.Client.GetDeviceList(ct.nameSpace)
	if err != nil {
		fmt.Printf("GetDeviceList Failed: %s \n", err.Error())
		return nil, err
	}
	dev := make([]Device, 0, 1)
	for _, x := range ret.Items {
		dev = append(dev, TransDevice(&x))
	}
	return dev, nil
}

func (ct *Controller) GetDevice(deviceName string) (*Device, error) {
	ret, err := ct.Client.GetDevice(ct.nameSpace, deviceName)
	if err != nil {
		fmt.Printf("GetDevice Failed: %s \n", err.Error())
		return nil, err
	}
	ans := TransDevice(ret)
	return &ans, nil
}

func (ct *Controller) GetNodeList() ([]Node, error) {
	nds, err := ct.Client.GetNodeList()
	if err != nil {
		fmt.Printf("GetNodeList Failed: %s \n", err.Error())
		return nil, err
	}
	nodes := make([]Node, 0, 1)
	for _, x := range nds.Items {
		nodes = append(nodes, Node{
			ResMeta: ResMeta{
				Name:  x.Name,
				State: utils.GetNodeStatus(x),
				Desc:  "Created Time:" + x.ObjectMeta.CreationTimestamp.Format(utils.TimeLayoutStr),
			},
			NodeVersion: x.Status.NodeInfo.KubeletVersion,
		})
	}
	return nodes, nil
}

func (ct *Controller) GetDeviceModelList() ([]DeviceModel, error) {
	dm, err := ct.Client.GetDeviceModelList(ct.nameSpace)
	if err != nil {
		fmt.Printf("GetDeviceModelList Failed: %s \n", err.Error())
		return nil, err
	}
	ret := make([]DeviceModel, 0, 1)
	for _, x := range dm.Items {
		nowProp := make([]utils.Property, 0, 1)
		for _, y := range x.Spec.Properties {
			nowProp = append(nowProp, utils.Property{
				Name:       y.Name,
				Desc:       y.Description,
				DefaultVal: y.Type.String.DefaultValue,
			})
		}
		ret = append(ret, DeviceModel{
			ResMeta: ResMeta{
				Name:  x.Name,
				State: "",
				Desc:  "Created Time:" + x.ObjectMeta.CreationTimestamp.Format(utils.TimeLayoutStr),
			},
			Properties: nowProp,
		})
	}
	return ret, nil
}

func (ct *Controller) AddDeviceModel(deviceModel DeviceModel) error {
	fileName, err := utils.GenerateModelYaml(ct.nameSpace, deviceModel.Name, deviceModel.Properties)
	if err != nil {
		fmt.Printf("AddDeviceModel GenerateModelYaml Failed: %s \n", err.Error())
		return err
	}
	return ct.Client.AddDeviceModelFromYaml(ct.nameSpace, fileName)
}

func (ct *Controller) AddDevice(device Device) error {
	fileName, err := utils.GenerateDeviceYaml(ct.nameSpace, device.Name, device.Desc, device.EdgeNode,
		utils.DeviceModelYamlPrefix+device.ModelName+".yaml")
	if err != nil {
		fmt.Printf("AddDeviceModel GenerateModelYaml Failed: %s \n", err.Error())
		return err
	}
	return ct.Client.AddDeviceFromYaml(ct.nameSpace, fileName)
}

func (ct *Controller) ChangeDeviceTwins(deviceName string, key string, value string) error {
	de, err := ct.Client.GetDevice(ct.nameSpace, deviceName)
	if err != nil {
		fmt.Printf("ChangeDeviceTwins GetDevice Failed: %s \n", err.Error())
		return err
	}
	changeMap := map[string]string{}
	for _, x := range de.Status.Twins {
		if x.PropertyName == key {
			changeMap[key] = value
		} else {
			changeMap[x.PropertyName] = x.Desired.Value
		}
	}
	return ct.Client.SetDeviceTwin(ct.nameSpace, deviceName, changeMap)
}
func (ct *Controller) GetPod(appName string) (*Pod, error) {
	nowList, err := ct.Client.GetPodList(ct.nameSpace)
	if err != nil {
		return nil, err
	}
	index := -1
	for i, x := range nowList.Items {
		if strings.Contains(x.Name, appName) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("Controller GetPod Error: No such pod: %s", appName)
	}

	return &Pod{
		ResMeta: ResMeta{
			Name:  nowList.Items[index].Name,
			State: "",
			Desc:  "",
		},
		Status: nowList.Items[index].Status,
	}, nil
}

func TransDevice(device *devices.Device) Device {
	sta, err := utils.GetTwinFromDevice(device, utils.DeviceStatusName, false)
	if err != nil {
		fmt.Println("GetDeviceList DeviceStatus desired not found")
	}
	dsta, err := utils.GetTwinFromDevice(device, utils.DeviceStatusName, false)
	if err != nil {
		fmt.Println("GetDeviceList DeviceStatus reported not found")
	}
	tstamp, err := utils.GetTimeStampFromDevice(device)
	if err != nil {
		fmt.Println("GetDeviceList DeviceStatus TimeStamp not found in:", device.Name)
	}
	return Device{
		ResMeta: ResMeta{
			Name:  device.Name,
			State: sta,
			Desc:  device.ObjectMeta.Labels["description"],
		},
		ModelName:   device.Spec.DeviceModelRef.Name,
		EdgeNode:    device.Spec.NodeSelector.NodeSelectorTerms[0].MatchExpressions[0].Values[0],
		DesireState: dsta,
		TimeStamp:   tstamp,
	}
}
