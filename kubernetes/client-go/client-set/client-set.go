package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	coreV1 := clientSet.CoreV1()
	pod, err := coreV1.Pods("kube-system").Get(context.TODO(), "coredns-6d8c4cb4d-qk64f", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("pod name is: ", pod.Name)
		fmt.Printf("pod.Status.PodIP: %v\n", pod.Status.PodIP)
	}
}
