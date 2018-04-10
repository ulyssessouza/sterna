package main

import (
	"os"
	"bufio"
	"flag"
	"fmt"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/homedir"

	"github.com/ulyssessouza/sterna/managers"
)

func main1() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	var deploymentManager managers.Manager = managers.NewDeploymentManager(*kubeconfig, apiv1.NamespaceDefault)

	const demoDeploymentName = "demo-deployment"
	var demoDeploymentNameCloned = fmt.Sprintf("%s-%s", demoDeploymentName, "cloned")

	deploymentManager.CreateExample(demoDeploymentName)
	prompt()
	deploymentManager.Clone(demoDeploymentName, demoDeploymentNameCloned)
	prompt()
	deploymentManager.Update(demoDeploymentName)
	prompt()
	deploymentManager.List()
	prompt()
	deploymentManager.Delete(demoDeploymentName, demoDeploymentNameCloned)
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
