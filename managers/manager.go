package managers

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	SERVICE      = "Service"
	DEPLOYMENT   = "Deployment"
	STATEFUL_SET = "StatefulSet"
	PDB_SET = "PodDisruptionBudget"
)

type Manager interface {
	Cloner
	CreateExample(name string)
	Update(oldSelector string, newSelector string)
	List()
	Delete(names ...string)
}

type Cloner interface {
	Clone(toBeCloned string, cloneName string)
	CloneInline(toBeCloned string, cloneName string, inplace bool)
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
