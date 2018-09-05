package core

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/docker/libcompose/config"
	composeProject "github.com/docker/libcompose/project"
	"gopkg.in/yaml.v2"
)

type Compose struct {
	Path string
}

type ComposeConfig struct {
	Services map[string]*config.ServiceConfig          `yaml:"services,omitempty"`
	Volumes  map[string]*config.VolumeConfig  `yaml:"volumes,omitempty"`
	Networks map[string]*config.NetworkConfig  `yaml:"networks,omitempty"`
}

func (compose *Compose) runCmd(cmdStr string, args... string) error {

	allArgs := cleanArgs(append([]string{cmdStr}, args...))

	cmd := exec.Command("docker-compose", allArgs...)

	cmd.Dir = compose.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (compose *Compose) Stop(args []string) error {
	return compose.runCmd("stop", args...)
}

func (compose *Compose) Start(args []string) error {
	return compose.runCmd("start", args...)
}

func (compose *Compose) Up(args []string) error {
	return compose.runCmd("up", args...)
}

func (compose *Compose) Down(args []string) error {
	return compose.runCmd("down", args...)
}

func (compose *Compose) Merge(repositories []Repository) error {

	// approach 1: use libcompose to stitch a fictional yaml file. The problem with this is that
	// all relative path references are wrong... This is not an easy fix

	var composeObj composeProject.APIProject
	// var composeCtx *ctx.Context
	var err error

	for _, repository := range repositories {

		_, hasPath := repository.GetComposePath()
		if !hasPath {
			continue
		}

		if composeObj == nil {
			composeObj, _, err = repository.GetComposeObject()
			if err != nil {
				return err
			}

		} else {
			composeBytes, err := repository.GetComposeBytes()
			if err != nil {
				return err
			}
			composeObj.Load(composeBytes)
		}

	}

	composeProjectObj := composeObj.(*composeProject.Project)

	theConfig := ComposeConfig{
		Services: composeProjectObj.ServiceConfigs.All(),
		Volumes: composeProjectObj.VolumeConfigs,
		Networks: composeProjectObj.NetworkConfigs,
	}

	content, err := yaml.Marshal(theConfig)
	if err != nil {
		panic(err)
	}
	// content := service_content + "\n" + volume_content + "\n" + network_content + "\n"
	if err := ioutil.WriteFile("merged.yaml", content, 0644); err != nil {
		panic(err)
	}


	// approach 2: blindly merge yaml... obviously not a good idea
	// var rawComposeObjs []map[string]interface{}
	//
	// for _, path := range paths {
	// 	rawComposeObj := make(map[string]interface{})
	// 	content, err := ioutil.ReadFile(path)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if err := yaml.Unmarshal(content, &rawComposeObj); err != nil {
	// 		panic(err)
	// 	}
	// 	rawComposeObjs = append(rawComposeObjs, rawComposeObj)
	// }
	//
	// merged := make(map[string]interface{})
	// for _, obj := range rawComposeObjs {
	// 	mergo.Merge(&merged, obj)
	// }
	//
	// content, err := yaml.Marshal(merged)
	// if err != nil {
	// 	panic(err)
	// }
	// if err := ioutil.WriteFile("merged.yaml", content, 0644); err != nil {
	// 	panic(err)
	// }



	return nil
}