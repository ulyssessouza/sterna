package main

import (
	"bufio"
	"flag"
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

	prompt("Will execute pre-migrations scripts")

	// Execute pre-clone scripts
	for _, m := range c.Migrations {
		if strings.TrimSpace(m.PreCloneScript) != "" {
			execCommand(m.PreCloneScript)
		}
	}

	prompt("Will execute cloning/update steps")

	// Clone supported controllers
	for _, m := range c.Migrations {
		switch m.MigrationType {
		case managers.SERVICE:
			if m.Inplace {
				managers.NewServiceManager(*kubeconfig, metaV1.NamespaceDefault).Update(m.Name, m.Selector)
			} else {
				managers.NewServiceManager(*kubeconfig, metaV1.NamespaceDefault).Clone(m.Name, m.ClonedName)
			}
		case managers.DEPLOYMENT:
			managers.NewDeploymentManager(*kubeconfig, apiv1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		case managers.STATEFUL_SET:
			managers.NewStatefulSetManager(*kubeconfig, metaV1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		case managers.PDB_SET:
			managers.NewPodDisruptionBudgetManager(*kubeconfig, metaV1.NamespaceDefault).Clone(m.Name, m.ClonedName)
		}
	}

	prompt("Will execute post-migration scripts")

	// Execute post-clone scripts
	for _, m := range c.Migrations {
		if strings.TrimSpace(m.PostCloneScript) != "" {
			execCommand(m.PostCloneScript)
		}
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

func prompt(msg string) {
	if *forced {
		return
	}
	log.Printf("%s -> Press Return key to continue.", msg)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
