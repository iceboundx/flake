package flake

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"icebound.cc/flake/cloud/controller"
	"io/ioutil"
	"strconv"
)

type Aggcore struct {
	Name           string
	SourceYamlfile string
	Yamlfile       string
	Rounds         int
	Port           int
	Cmd            []string
	podName        string
}

type EdgeAggcore struct {
	Aggcore
	K        int
	CloudUrl string
	CertName string
}

func (ac *Aggcore) DeployAggcore(ctl *controller.Controller) error {
	return ctl.Client.ApplyYaml(ctl.GetNameSpace(), ac.Yamlfile)
}

func (ac *Aggcore) CheckState(ctl *controller.Controller) (string, error) {
	nowpod, err := ctl.GetPod(ac.podName)
	if err != nil {
		return "", err
	}
	return string(nowpod.Status.Phase), nil
}

func NewCloudAggcore(name string, yamlfile string, rounds int, port int) (*Aggcore, error) {
	cmd := []string{"python", "server.py", strconv.Itoa(rounds), "0.0.0.0:" + strconv.Itoa(port)}
	podName, err := generateAggcoreYaml(yamlfile, name+".yaml", rounds, port, cmd)
	if err != nil {
		return nil, err
	}
	return &Aggcore{
		Name:           name,
		SourceYamlfile: yamlfile,
		Yamlfile:       name + ".yaml",
		Rounds:         rounds,
		Port:           port,
		Cmd:            cmd,
		podName:        podName,
	}, nil
}

func NewEdgeAggcore(name string, yamlfile string, K int, rounds int, port int, cloudUrl string, certName string) (*EdgeAggcore, error) {
	cmd := []string{"python3", "server_edge.py", strconv.Itoa(K), strconv.Itoa(rounds),
		"0.0.0.0:" + strconv.Itoa(port), cloudUrl, certName}
	podName, err := generateAggcoreYaml(yamlfile, name+".yaml", rounds, port, cmd)
	if err != nil {
		return nil, err
	}
	return &EdgeAggcore{
		Aggcore: Aggcore{
			Name:           name,
			SourceYamlfile: yamlfile,
			Yamlfile:       name + ".yaml",
			Rounds:         rounds,
			Port:           port,
			Cmd:            cmd,
			podName:        podName,
		},
		K:        K,
		CloudUrl: cloudUrl,
		CertName: certName,
	}, nil
}

func generateAggcoreYaml(soucefile string, filename string, round int, port int, cmd []string) (string, error) {
	filebytes, err := ioutil.ReadFile(soucefile)
	if err != nil {
		return "", err
	}
	result := make(map[interface{}]interface{})
	err = yaml.Unmarshal(filebytes, &result)
	if err != nil {
		return "", err
	}
	meta := result["metadata"].(map[interface{}]interface{})
	name := meta["name"].(string)
	spec1 := result["spec"].(map[interface{}]interface{})
	templ := spec1["template"].(map[interface{}]interface{})
	spec2 := templ["spec"].(map[interface{}]interface{})
	cont := ((spec2["containers"].([]interface{}))[0]).(map[interface{}]interface{})
	a := cmd
	b := ((cont["ports"].([]interface{}))[0]).(map[interface{}]interface{})
	b["containerPort"] = port
	b["hostPort"] = port
	cont["command"] = a
	fmt.Println(result)

	text, err := yaml.Marshal(&result)
	if err != nil {
		return "", err
	}
	return name, ioutil.WriteFile(filename, text, 0777)

}
