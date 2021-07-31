package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Property struct {
	Name       string
	Desc       string
	DefaultVal string
}

func ChangeYaml(sourceFilePath string, outputFilePath string, changeMap map[string]interface{}) error {
	filebytes, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}
	result := make(map[interface{}]interface{})
	err = yaml.Unmarshal(filebytes, &result)
	if err != nil {
		return err
	}
	mapChange(result, changeMap)
	text, err := yaml.Marshal(&result)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFilePath, text, 0777)
}

func mapChange(sourceMap map[interface{}]interface{}, changeMap map[string]interface{}) {
	if sourceMap == nil || changeMap == nil {
		return
	}
	for k, v := range changeMap {
		if _, ok := sourceMap[k]; ok {
			val, isStr := v.(string)
			if isStr == true {
				sourceMap[k] = val
				fmt.Println("Change ", k, " ", val)
				continue
			} else {
				x := sourceMap[k].(map[interface{}]interface{})
				y := changeMap[k].(map[string]interface{})
				mapChange(x, y)
			}
		} else {
			sourceMap[k] = changeMap[k]
			fmt.Println("Change ", k, " ", changeMap[k])
		}
	}
}

func GenerateModelYaml(nameSpace string, name string, properties []Property) (string, error) {
	propertiesArray := make([]map[string]interface{}, 0, 1)
	for _, x := range properties {
		tmp := map[string]interface{}{
			"name":        x.Name,
			"description": x.Desc,
			"type": map[string]interface{}{
				"string": map[string]interface{}{
					"accessMode":   AccessModeReadWrite,
					"defaultValue": x.DefaultVal,
				},
			},
		}
		propertiesArray = append(propertiesArray, tmp)
	}
	mp := map[string]interface{}{
		"apiVersion": "devices.kubeedge.io/v1alpha2",
		"kind":       "DeviceModel",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": nameSpace,
		},
		"spec": map[string]interface{}{
			"properties": propertiesArray,
		},
	}
	text, err := yaml.Marshal(&mp)
	if err != nil {
		return "", err
	}
	fileName := DeviceModelYamlPrefix + name + ".yaml"
	err = ioutil.WriteFile(fileName, text, 0777)
	return fileName, nil
}

func GenerateDeviceYaml(nameSpace string, deviceName string, desc string, edgeNode string, modelYamlPath string) (string, error) {

	properties, err := GetModelProperties(modelYamlPath)
	if err != nil {
		return "", err
	}
	modelName, err := GetModelName(modelYamlPath)
	if err != nil {
		return "", err
	}
	prop := make([]map[string]interface{}, 0, 1)
	for _, x := range properties {
		prop = append(prop, map[string]interface{}{
			"propertyName": x.Name,
			"desired": map[string]interface{}{
				"metadata": map[string]interface{}{
					"type": "string",
				},
				"value": x.DefaultVal,
			},
			"reported": map[string]interface{}{
				"metadata": map[string]interface{}{
					"type": "string",
				},
				"value": "",
			},
		})
	}
	mp := map[string]interface{}{
		"apiVersion": "devices.kubeedge.io/v1alpha2",
		"kind":       "Device",
		"metadata": map[string]interface{}{
			"name":      deviceName,
			"namespace": nameSpace,
			"labels": map[string]interface{}{
				"description": desc,
			},
		},
		"spec": map[string]interface{}{
			"deviceModelRef": map[string]interface{}{
				"name": modelName,
			},
			"nodeSelector": map[string]interface{}{
				"nodeSelectorTerms": []map[string]interface{}{
					{
						"matchExpressions": []map[string]interface{}{
							{
								"key":      "",
								"operator": "In",
								"values": []string{
									edgeNode,
								},
							},
						},
					},
				},
			},
		},
		"status": map[string]interface{}{
			"twins": prop,
		},
	}
	text, err := yaml.Marshal(&mp)
	if err != nil {
		return "", err
	}
	fileName := DeviceYamlPrefix + deviceName + ".yaml"
	err = ioutil.WriteFile(fileName, text, 0777)
	return fileName, nil
}

func GetModelProperties(filePath string) ([]Property, error) {
	filebytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(filebytes, &result)
	if err != nil {
		return nil, err
	}
	raw := result["spec"].(map[interface{}]interface{})["properties"]
	pro := raw.([]interface{})
	ret := make([]Property, 0, 1)
	for _, x := range pro {
		now := x.(map[interface{}]interface{})
		typ := now["type"].(map[interface{}]interface{})
		typeString := typ["string"].(map[interface{}]interface{})
		ret = append(ret, Property{
			Name:       now["name"].(string),
			Desc:       now["description"].(string),
			DefaultVal: typeString["defaultValue"].(string),
		})
	}
	return ret, nil
}

func GetModelName(filePath string) (string, error) {
	filebytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(filebytes, &result)
	if err != nil {
		return "", err
	}
	metaData := result["metadata"].(map[interface{}]interface{})
	return metaData["name"].(string), nil
}
