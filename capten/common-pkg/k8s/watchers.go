package k8s

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

const (
	// Sync every 60 sec for the updates.
	defaultResyncTime = 60
)

func RegisterDynamicInformers(resEvtHandler cache.ResourceEventHandler,
	client dynamic.Interface, gvr schema.GroupVersionResource,
) error {
	kubeInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(
		client,
		time.Duration(defaultResyncTime)*time.Second,
	)

	informer := kubeInformerFactory.ForResource(gvr).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    resEvtHandler.OnAdd,
		UpdateFunc: resEvtHandler.OnUpdate,
		DeleteFunc: resEvtHandler.OnDelete,
	})

	stop := make(chan struct{})
	kubeInformerFactory.Start(stop)

	if !cache.WaitForCacheSync(stop, kubeInformerFactory.ForResource(gvr).Informer().HasSynced) {
		return errors.New(fmt.Sprintf("Failed to cache the %v", gvr))
	}

	return nil
}
