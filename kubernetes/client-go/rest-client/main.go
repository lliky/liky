package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	config.GroupVersion = &v1.SchemeGroupVersion
	config.APIPath = "/api"
	config.NegotiatedSerializer = scheme.Codecs
	// client
	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	// get data
	pod := v1.Pod{}
	err = client.Get().Namespace("kube-system").Resource("pods").Name("kube-controller-manager-master01").Do(context.TODO()).Into(&pod)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("pod.Name: %v\n", pod.Name)
	}
}
