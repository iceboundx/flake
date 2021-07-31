package utils

import (
	"context"
	"fmt"
	devices "github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha2"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
)

type K8sClient struct {
	Client    *kubernetes.Clientset
	CrdClient *rest.RESTClient
	conf      *rest.Config
}

func NewK8sClient(isInCluster bool) (*K8sClient, error) {
	var config *rest.Config
	var confErr error
	if isInCluster {
		config, confErr = rest.InClusterConfig()
		if confErr != nil {
			panic(confErr.Error())
		}
	} else {
		config, confErr = OutClusterConf()
		if confErr != nil {
			panic(confErr.Error())
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	crdClient, err := NewDeviceClient(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{
		Client:    clientset,
		CrdClient: crdClient,
		conf:      config,
	}, nil
}

func (cl *K8sClient) GetPodList(nameSpace string) (*v1.PodList, error) {
	return cl.Client.CoreV1().Pods(nameSpace).List(context.Background(), metav1.ListOptions{})
}

func (cl *K8sClient) GetPod(nameSpace string, podName string) (pod *v1.Pod, err error) {
	pod, err = cl.Client.CoreV1().Pods(nameSpace).Get(context.Background(), podName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, nameSpace)
		return nil, err
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, nameSpace, statusError.ErrStatus.Message)
		return nil, err
	} else if err != nil {
		panic(err.Error())
	} else {
		return pod, nil
	}
}

func (cl *K8sClient) DeletePod(nameSpace string, podName string) (err error) {
	err = cl.Client.CoreV1().Pods(nameSpace).Delete(context.Background(), podName, metav1.DeleteOptions{})
	return err
}

func (cl *K8sClient) DeployApp(nameSpace string, yamlName string) error {
	fmt.Println("Try to deploy apps")
	err := cl.ApplyYaml(nameSpace, yamlName)
	return err
}

func (cl *K8sClient) DeleteApp(nameSpace string, appName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := cl.Client.AppsV1().Deployments(nameSpace).Delete(context.Background(), appName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	return err
}

func (cl *K8sClient) GetAppList(nameSpace string) (*appv1.DeploymentList, error) {
	return cl.Client.AppsV1().Deployments(nameSpace).List(context.Background(), metav1.ListOptions{})
}

func (cl *K8sClient) GetApp(nameSpace string, appName string) (*appv1.Deployment, error) {
	return cl.Client.AppsV1().Deployments(nameSpace).Get(context.Background(), appName, metav1.GetOptions{})
}

func (cl *K8sClient) GetNodeList() (*v1.NodeList, error) {
	return cl.Client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
}

func GetTwinFromDevice(device *devices.Device, proName string, isDesired bool) (string, error) {
	for _, x := range device.Status.Twins {
		if x.PropertyName == proName {
			if isDesired {
				return x.Desired.Value, nil
			} else {
				return x.Reported.Value, nil
			}
		}
	}
	return "", fmt.Errorf("Can't Find property: %s", proName)
}
func GetTimeStampFromDevice(device *devices.Device) (int64, error) {
	for _, x := range device.Status.Twins {
		if x.PropertyName == DeviceStatusName { //TODO
			ret, err := strconv.ParseInt(x.Reported.Metadata["timestamp"], 10, 64)
			if err != nil {
				return 0, err
			} else {
				return ret, nil
			}
		}
	}
	return 0, fmt.Errorf("Can't Find property: %s", "state")
}

func GetNodeStatus(node v1.Node) string {
	for _, x := range node.Status.Conditions {
		if x.Status == "True" {
			return string(x.Type)
		}
	}
	return "NotReady"
}
