package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	APIGroup, slices, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	fmt.Printf("APIGroup: %v\n\n\n", APIGroup)
	for _, v := range slices {
		gvStr := v.GroupVersion
		gv, err := schema.ParseGroupVersion(gvStr)
		if err != nil {
			panic(err)
		}
		fmt.Println("#############################################")
		fmt.Printf("GV string: [%v]\nGV struct [%#v]\n\n", gvStr, gv)
		for _, res := range v.APIResources {
			fmt.Printf("%v\n", res.Name)
		}
	}
}
