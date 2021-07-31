package controller

import "icebound.cc/flake/utils"

func NewLearnerDeviceModel(name string) DeviceModel {
	return DeviceModel{
		ResMeta: ResMeta{
			Name:  name,
			State: "",
			Desc:  "LearningDeviceModel",
		},
		Properties: []utils.Property{
			utils.Property{
				Name:       utils.DeviceStatusName,
				Desc:       "DeviceStatus",
				DefaultVal: "",
			},
			utils.Property{
				Name:       utils.DeviceTaskInfoName,
				Desc:       "DeviceTaskInfo",
				DefaultVal: "",
			},
			utils.Property{
				Name:       utils.DeviceInfoName,
				Desc:       "DeviceInfo",
				DefaultVal: "",
			},
			utils.Property{
				Name:       utils.DeviceTaskParam,
				Desc:       "TaskParam",
				DefaultVal: "python nn.py",
			},
		},
	}
}

func NewLearnerDeviceInstance(name string, desc string, modelName string, edgeNode string) Device {
	return Device{
		ResMeta: ResMeta{
			Name:  name,
			State: utils.DeviceStatusNotReady,
			Desc:  desc,
		},
		ModelName:   modelName,
		EdgeNode:    edgeNode,
		DesireState: utils.DeviceStatusNotReady,
		TimeStamp:   0,
	}
}
