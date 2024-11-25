package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	unstructured, err := dynamicClient.Resource(gvr).Namespace("kube-system").Get(context.Background(), "kube-proxy-2z7xf", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	pod := v1.Pod{}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured.UnstructuredContent(), &pod); err != nil {
		return
	}
	fmt.Printf("pod.Name: %v\n", pod.Name)
	fmt.Printf("pod.Status.PodIP: %v\n", pod.Status.PodIP)
}
