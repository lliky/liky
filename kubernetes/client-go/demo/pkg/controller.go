package pkg

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	serviceInformer "k8s.io/client-go/informers/core/v1"
	ingressInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	serviceV1 "k8s.io/client-go/listers/core/v1"
	ingressV1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	workNum  = 5
	maxRetry = 10
)

type controller struct {
	client        kubernetes.Interface
	serviceLister serviceV1.ServiceLister
	ingressLister ingressV1.IngressLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) updateService(oldObj, newObj interface{}) {
	// todo compare annotation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

func (c *controller) deleteIngress(obj interface{}) {

	ingress := obj.(*networkingV1.Ingress)

	ownerReference := metav1.GetControllerOf(ingress)
	fmt.Printf("reference: %+v\n", ownerReference)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "Service" {
		return
	}
	c.enqueue(obj)
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.Add(key)
}

func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.work, time.Minute, stopCh)
	}
	<-stopCh
}

func (c *controller) work() {
	for c.processNextItem() {

	}
}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	key := item.(string)

	err := c.syncService(key)
	if err != nil {
		c.handlerError(item, err)
	}
	return true
}

func (c *controller) syncService(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	log.Printf("svc: %+v", service)
	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if ok && errors.IsNotFound(err) {
		// create ingress
		ig := c.constructIngress(service)
		_, err := c.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ig, metav1.CreateOptions{})
		if err != nil {
			return err
		}

	} else if !ok && ingress != nil {
		// delete ingress
		err := c.client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *controller) constructIngress(service *coreV1.Service) *networkingV1.Ingress {
	ingress := networkingV1.Ingress{}
	ingress.Namespace = service.Namespace
	ingress.Name = service.Name

	ingress.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(service, coreV1.SchemeGroupVersion.WithKind("Service")),
	}
	icn := "nginx"
	pathType := networkingV1.PathTypePrefix
	ingress.Spec = networkingV1.IngressSpec{
		IngressClassName: &icn,
		Rules: []networkingV1.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: networkingV1.IngressRuleValue{
					HTTP: &networkingV1.HTTPIngressRuleValue{
						Paths: []networkingV1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: networkingV1.IngressBackend{
									Service: &networkingV1.IngressServiceBackend{
										Name: service.Name,
										Port: networkingV1.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &ingress
}

func (c *controller) handlerError(item interface{}, err error) {
	if c.queue.NumRequeues(item) <= maxRetry {
		c.queue.AddRateLimited(item)
		return
	}

	runtime.HandleError(err)
	c.queue.Forget(item)
}
func NewController(client kubernetes.Interface, serviceInformer serviceInformer.ServiceInformer, ingressInformer ingressInformer.IngressInformer) controller {
	c := controller{
		client:        client,
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return c
}
