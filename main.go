package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"

	"github.com/ulyssessouza/sterna/managers"
	"github.com/ulyssessouza/sterna/migration"
	"os/exec"
)

var forced *bool

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	var configFilePath = flag.String("config", "sterna.yml", "Configuration file path")
	forced = flag.Bool("force", false, "Apply the changes without asking for confirmation")
	flag.Parse()

	c := migration.Config{}
	c.Load(*configFilePath)

	// Execute pre-clone scripts
	for _, m := range c.Migrations {
		execCommand(m.PreCloneScript)
	}

	prompt()

	// Clone supported controllers
	for _, m := range c.Migrations {
		switch m.MigrationType {
		case managers.SERVICE:
			managers.NewServiceManager(*kubeconfig, metaV1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		case managers.DEPLOYMENT:
			managers.NewDeploymentManager(*kubeconfig, apiv1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		case managers.STATEFUL_SET:
			managers.NewStatefulSetManager(*kubeconfig, metaV1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		case managers.PDB_SET:
			managers.NewPodDisruptionBudgetManager(*kubeconfig, metaV1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		}
	}

	prompt()

	// Execute post-clone scripts
	for _, m := range c.Migrations {
		execCommand(m.PostCloneScript)
	}
}

func execCommand(c string) {
	split := strings.Split(c, " ")
	out, err := exec.Command(split[0], split[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Script output is: %s\n", out)
}

func prompt() {
	if *forced {
		return
	}
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
