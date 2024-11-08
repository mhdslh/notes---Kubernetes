package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	k8sConfig *rest.Config
	scheme    *runtime.Scheme
)

func init() {
	var err error
	k8sConfig, err = config.GetConfig()
	if err != nil {
		panic(err)
	}

	scheme = runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
}

func main() {
	k8sConfig.APIPath = "/api"
	k8sConfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	k8sConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme).WithoutConversion()

	restClient, err := rest.RESTClientFor(k8sConfig)
	if err != nil {
		panic(err)
	}

	lw := cache.NewListWatchFromClient(restClient, "pods", "default", make(fields.Set).AsSelector())
	store := cache.NewStore(func(obj interface{}) (string, error) {
		pod, ok := obj.(*corev1.Pod)
		if !ok {
			return "", fmt.Errorf("type cast to *corev1.Pod failed")
		}
		return pod.GetNamespace() + "/" + pod.GetName(), nil
	})
	reflector := cache.NewReflector(lw, &corev1.Pod{}, store, time.Second)

	stopChannel := make(chan struct{})

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		reflector.Run(stopChannel)
		log.Println("Shutting down reflector")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				log.Println("Pods: ", store.ListKeys())
			case <-stopChannel:
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}
