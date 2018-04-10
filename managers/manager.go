package managers

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Manager interface {
	Cloner
	CreateExample(name string)
	Update(name string)
	List()
	Delete(names ...string)
}

type Cloner interface {
	Clone(toBeCloned string, cloneName string)
}

func getClientSet(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}