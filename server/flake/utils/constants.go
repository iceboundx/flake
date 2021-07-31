package utils

const (
	// MergePatchType is patch type
	MergePatchType = "application/merge-patch+json"
	JsonPatchType  = "application/json-patch+json"
	// ResourceTypeDevices is plural of device resource in apiserver
	ResourceTypeDevices     = "devices"
	ResourceTypeDeviceModel = "devicemodels"
)

const (
	AccessModeReadOnly  = "ReadOnly"
	AccessModeReadWrite = "ReadWrite"
)

const (
	DeviceModelYamlPrefix = "devicemodel_"
	DeviceYamlPrefix      = "device_"
)

const (
	DeviceStatusNotReady         = "NotReady"
	DeviceStatusOnline           = "Online"
	DeviceStatusOffline          = "Offline"
	DeviceStatusReadyForLearning = "ReadyForLearning"
	DeviceStatusLearning         = "Learning"
)

const TimeLayoutStr = "2006-01-02 15:04:05"

const (
	DeviceStatusName   = "status"
	DeviceTaskInfoName = "taskinfo"
	DeviceInfoName     = "deviceinfo"
	DeviceTaskParam    = "taskparam"
)

const (
	PodRunning = "Running"
)
