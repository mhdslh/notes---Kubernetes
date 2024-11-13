package main

import (
	"fmt"
	"log/slog"
	"os"
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
	scheme    *runtime.Scheme
	k8sConfig *rest.Config
	logger    *slog.Logger
)

func init() {
	scheme = runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)

	var err error
	k8sConfig, err = config.GetConfig()
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	logger = slog.New(handler)
}

func main() {
	restConfig := rest.CopyConfig(k8sConfig)
	restConfig.APIPath = "/api"
	restConfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	restConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme).WithoutConversion()

	restClient, err := rest.RESTClientFor(restConfig)
	if err != nil {
		logger.Error("Failed to create a REST client", "error", err)
		return
	}

	lw := cache.NewListWatchFromClient(restClient, "pods", "default", make(fields.Set).AsSelector())
	queue := cache.NewFIFO(cache.MetaNamespaceKeyFunc) // cache.MetaNamespaceKeyFunc returns "<namespace>/<name>" for namespace-scoped resources and "<name>" for cluster-scoped resources.
	// The FIFO implementation uses a map to store key-object pairs and a queue to manage the order of keys for processing.
	// When an object is updated, the associated key is not re-added to the queue; instead, the value for that key is updated in the map.
	// During the Pop operation, the key at the front of the queue is removed, the corresponding item is deleted from the map, and then it is processed.
	// In case of a requeue error, the item is re-inserted into both the queue and the map, but only if it doesn't already exist.
	process := func(obj interface{}, isInInitialList bool) error {
		time.Sleep(time.Second)
		// Invoking queue methods within PopProcessFunc could result in a deadlock
		pod, ok := obj.(*corev1.Pod)
		if !ok {
			return fmt.Errorf("type cast to *corev1.Pod failed")
		}
		logger.Info("Processing", "namespace", pod.GetNamespace(), "name", pod.GetName())
		return nil
	}
	config := &cache.Config{
		Queue:         queue,
		ListerWatcher: lw,
		Process:       process,
		ObjectType:    &corev1.Pod{},
	}
	controller := cache.New(config)

	stopChannel := make(chan struct{})
	wg := &sync.WaitGroup{}

	// Items that exist in the store are removed after they are processed.
	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				logger.Info("Store", "objects", queue.ListKeys())
			case <-stopChannel:
				wg.Done()
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		controller.Run(stopChannel)
		wg.Done()
	}()

	<-stopChannel
	wg.Wait()
}
