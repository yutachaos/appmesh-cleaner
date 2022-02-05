package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appmesh"
	"log"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	meshName := flag.String("mesh-name", "", "please specify -mesh-name flag")
	flag.Parse()
	if len(*meshName) == 0 {
		return fmt.Errorf("you need specify -mesh-name flag")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}
	svc := appmesh.NewFromConfig(cfg)

	virtualServiceList, err := svc.ListVirtualServices(context.TODO(), &appmesh.ListVirtualServicesInput{MeshName: meshName})
	if err != nil {
		return fmt.Errorf("failed to ListVirtualServices. %v", err)
	}

	for _, virtualService := range virtualServiceList.VirtualServices {
		fmt.Printf("Delete VirtualService: %s\n", *virtualService.VirtualServiceName)
		_, err = svc.DeleteVirtualService(context.TODO(), &appmesh.DeleteVirtualServiceInput{MeshName: meshName, VirtualServiceName: virtualService.VirtualServiceName})
		if err != nil {
			return fmt.Errorf("failed to DeleteVirtualService, %v", err)
		}
	}

	virtualRouterList, err := svc.ListVirtualRouters(context.TODO(), &appmesh.ListVirtualRoutersInput{MeshName: meshName})
	if err != nil {
		return fmt.Errorf("failed to ListVirtualRouters, %v", err)
	}

	for _, virtualRouter := range virtualRouterList.VirtualRouters {
		fmt.Printf("Delete VirtualRouter: %s\n", *virtualRouter.VirtualRouterName)
		routeList, err := svc.ListRoutes(context.TODO(), &appmesh.ListRoutesInput{MeshName: meshName, VirtualRouterName: virtualRouter.VirtualRouterName})
		if err != nil {
			return fmt.Errorf("failed to ListRoutes, %v", err)
		}
		for _, route := range routeList.Routes {
			fmt.Printf("Delete Route: %s\n", *route.RouteName)
			_, err = svc.DeleteRoute(context.TODO(), &appmesh.DeleteRouteInput{MeshName: meshName, VirtualRouterName: virtualRouter.VirtualRouterName, RouteName: route.RouteName})
			if err != nil {
				return fmt.Errorf("failed to DeleteRoute, %v", err)
			}
		}
		_, err = svc.DeleteVirtualRouter(context.TODO(), &appmesh.DeleteVirtualRouterInput{MeshName: meshName, VirtualRouterName: virtualRouter.VirtualRouterName})
		if err != nil {
			return fmt.Errorf("failed to DeleteVirtualRouter, %v", err)
		}
	}

	virtualNodeList, err := svc.ListVirtualNodes(context.TODO(), &appmesh.ListVirtualNodesInput{MeshName: meshName})
	if err != nil {
		return fmt.Errorf("failed to ListVirtualNodes, %v", err)
	}

	for _, virtualNode := range virtualNodeList.VirtualNodes {
		fmt.Printf("Delete VirtualNode: %s\n", *virtualNode.VirtualNodeName)
		_, err = svc.DeleteVirtualNode(context.TODO(), &appmesh.DeleteVirtualNodeInput{MeshName: meshName, VirtualNodeName: virtualNode.VirtualNodeName})
		if err != nil {
			return fmt.Errorf("failed to DeleteVirtualNode, %v", err)
		}
	}
	fmt.Printf("Delete Mesh: %s\n", *meshName)

	_, err = svc.DeleteMesh(context.TODO(), &appmesh.DeleteMeshInput{MeshName: meshName})
	if err != nil {
		return fmt.Errorf("failed to DeleteMesh %v", err)
	}
	return nil
}
