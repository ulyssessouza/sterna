package main

import (
	"flag"
	"github.com/ulyssessouza/sterna/managers"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

func mainServices() {
	var kubeconfig *string
	home := homedir.HomeDir()
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var serviceManager managers.Manager = managers.NewServiceManager(*kubeconfig, metaV1.NamespaceDefault)

	serviceManager.List()
	log.Println("Will create a new Service")
	prompt()

	serviceManager.CreateExample("nginx")

	log.Printf("Created Service\n")

	prompt()
	serviceManager.List()
	log.Printf("Deleting Service\n")
	prompt()
	serviceManager.Delete("nginx")
	log.Printf("Service deleted\n")
}
