package utils

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	devices "github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha2"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

// DeviceStatus is used to patch device status
type DeviceStatus struct {
	Status devices.DeviceStatus `json:"status"`
}

// NewDeviceClient is used to create a restClient for crd
func NewDeviceClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(addDeviceCrds)

	err := schemeBuilder.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	config := *cfg
	config.APIPath = "/apis"
	config.GroupVersion = &devices.SchemeGroupVersion
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		log.Fatalf("Failed to create REST Client due to error %v", err)
		return nil, err
	}

	return client, nil
}

func addDeviceCrds(scheme *runtime.Scheme) error {
	// Add Device
	scheme.AddKnownTypes(devices.SchemeGroupVersion, &devices.Device{}, &devices.DeviceList{})
	v1.AddToGroupVersion(scheme, devices.SchemeGroupVersion)

	// Add DeviceModel
	scheme.AddKnownTypes(devices.SchemeGroupVersion, &devices.DeviceModel{}, &devices.DeviceModelList{})
	v1.AddToGroupVersion(scheme, devices.SchemeGroupVersion)

	return nil
}

func (cl *K8sClient) GetDeviceList(nameSpace string) (res *devices.DeviceList, err error) {
	raw, err := cl.CrdClient.Get().
		Namespace(nameSpace).Resource(ResourceTypeDevices).DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cl *K8sClient) GetDevice(nameSpace string, deviceName string) (res *devices.Device, err error) {
	raw, err := cl.CrdClient.Get().
		Namespace(nameSpace).Resource(ResourceTypeDevices).Name(deviceName).DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cl *K8sClient) AddDeviceFromYaml(nameSpace string, filePath string) error {
	return cl.ApplyYaml(nameSpace, filePath)
}

//TODO: AddDevice From Object
/*
func (cl *K8sClient) AddDevice(nameSpace string, device devices.Device) error {
	device.
}
*/

func (cl *K8sClient) DeleteDevice(nameSpace string, deviceName string) error {
	res := cl.CrdClient.Delete().
		Namespace(nameSpace).Resource(ResourceTypeDevices).Name(deviceName).Do(context.TODO())
	return res.Error()
}

func (cl *K8sClient) GetDeviceModelList(nameSpace string) (res *devices.DeviceModelList, err error) {
	raw, err := cl.CrdClient.Get().
		Namespace(nameSpace).Resource(ResourceTypeDeviceModel).DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cl *K8sClient) GetDeviceModel(nameSpace string, modelName string) (res *devices.DeviceModel, err error) {
	raw, err := cl.CrdClient.Get().
		Namespace(nameSpace).Resource(ResourceTypeDeviceModel).Name(modelName).DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cl *K8sClient) DeleteDeviceModel(nameSpace string, modelName string) error {
	res := cl.CrdClient.Delete().
		Namespace(nameSpace).Resource(ResourceTypeDeviceModel).Name(modelName).Do(context.TODO())
	return res.Error()
}

func (cl *K8sClient) AddDeviceModelFromYaml(nameSpace string, filePath string) error {
	return cl.ApplyYaml(nameSpace, filePath)
}

func (cl *K8sClient) SetDeviceTwin(nameSpace string, deviceName string, changeMap map[string]string) error {

	status := buildTwinData(changeMap)
	deviceStatus := &DeviceStatus{Status: status}
	body, err := json.Marshal(deviceStatus)
	if err != nil {
		log.Printf("Failed to marshal device status %v", deviceStatus)
		return err
	}
	result := cl.CrdClient.Patch(MergePatchType).Namespace(nameSpace).
		Resource(ResourceTypeDevices).Name(deviceName).Body(body).Do(context.TODO())
	if result.Error() != nil {
		log.Printf("Failed to patch device status %v of device %v in namespace %v \n error:%+v",
			deviceStatus, deviceName, nameSpace, result.Error())
		return result.Error()
	} else {
		log.Printf("Turn %s %s", body, deviceName)
	}

	return nil
}

func buildTwinData(changeMap map[string]string) devices.DeviceStatus {
	metadata := map[string]string{
		"timestamp": strconv.FormatInt(time.Now().Unix()/1e6, 10),
		"type":      "string",
	}
	twins := make([]devices.Twin, 0, 1)
	for k, v := range changeMap {
		twins = append(twins, devices.Twin{PropertyName: k,
			Desired:  devices.TwinProperty{Value: v, Metadata: metadata},
			Reported: devices.TwinProperty{Value: "", Metadata: metadata}})
	}

	devicestatus := devices.DeviceStatus{Twins: twins}
	return devicestatus
}
