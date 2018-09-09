package core

import (
	"context"
	"fmt"


	"github.com/docker/libcompose/docker/ctx"
	composeProject "github.com/docker/libcompose/project"
	composeOptions "github.com/docker/libcompose/project/options"
)


func (project *Project) assembleComposeObject() (*composeProject.Project, *ctx.Context, error) {
	var composeObj *composeProject.Project
	var composeCtx *ctx.Context
	var err error
	var tmp composeProject.APIProject

	for _, repository := range project.Repositories {

		_, hasPath := repository.GetComposePath()
		if !hasPath {
			continue
		}

		if composeObj == nil {
			tmp, composeCtx, err = repository.GetComposeObject()
			composeObj = tmp.(*composeProject.Project)
			if err != nil {
				return nil, nil, err
			}

		} else {
			composeBytes, err := repository.GetComposeBytes()
			if err != nil {
				return nil, nil, err
			}
			composeObj.Load(composeBytes)
		}

	}

	// create a new network and add all services to this network
	// lan := "stitch-lan"
	// composeCtx.Project.AddNetworkConfig(lan, &composeConfig.NetworkConfig{})
	// for _, config := range composeCtx.Project.ServiceConfigs.All() {
	// 	network := &composeYaml.Network{}
	// 	network.Name = "stitch-lan"
	// 	network.RealName = project.SafeName() + "_" + lan
	// 	config.Networks.Networks = append(config.Networks.Networks, network)
	// }
	//
	//
	// fmt.Println("Services...")
	// for name, config := range composeObj.ServiceConfigs.All() {
	// 	fmt.Printf("   %-15s %+v\n", name, config)
	// }
	//
	// fmt.Println("Volumes...")
	// for name, config := range composeObj.VolumeConfigs {
	// 	fmt.Printf("   %-15s %+v\n", name, config)
	// }
	//
	// fmt.Println("Networks...")
	// for name, config := range composeObj.NetworkConfigs {
	// 	fmt.Printf("   %-15s %+v\n", name, config)
	// }

	return composeObj, composeCtx, nil
}



// func (project *Project) ComposeUp(services ...string) error {
// 	if err := project.Compose.Create(context.Background(), composeOptions.Create{}, services...); err != nil {
// 		return err
// 	}
//
// 	if err := project.Compose.Up(context.Background(), composeOptions.Up{}, services...); err != nil {
// 		return err
// 	}
//
// 	project.Compose.Log(context.Background(), true)
// 	// wait forever
// 	<-make(chan interface{})
//
//
// 	return nil
// }

func (project *Project) ComposeLogs(services ...string) error {
	fmt.Println("Following logs...")
	project.Compose.Log(context.Background(), true, services...)
	// wait forever
	<-make(chan interface{})
	return nil
}

func (project *Project) ComposeUp(services ...string) error {
	// create twice to fix circular links
	// if err := project.Compose.Create(context.Background(), composeOptions.Create{}, services...); err != nil {
	// 	return err
	// }
	// if err := project.Compose.Create(context.Background(), composeOptions.Create{}, services...); err != nil {
	// 	return err
	// }

	if err := project.Compose.Up(context.Background(), composeOptions.Up{}); err != nil {
		return err
	}

	return project.ComposeLogs(services...)
}

func (project *Project) ComposeDown() error {
	return project.Compose.Down(context.Background(), composeOptions.Down{})
}