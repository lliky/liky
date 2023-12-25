package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liky/client-go-demo/11/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. config
	// 2. client
	// 3. informer
	// 4. add event handler
	// 5. informer.Start

	// 集群外部获取 config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	fmt.Println(config)
	if err != nil {
		// 集群内部获取 config
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln("can't get config")
		}
		config = inClusterConfig
		context.Background()
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can't crate client")
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)

	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()

	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	// 所有全部同步完成
	factory.WaitForCacheSync(stopCh)

	controller.Run(stopCh)
}
