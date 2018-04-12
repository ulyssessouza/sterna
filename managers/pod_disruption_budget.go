package managers

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tV1Beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	"log"
)

func NewPodDisruptionBudgetManager(kubeconfig string, namespace string) *PodDisruptionBudgetManager {
	clientset := getClientSet(kubeconfig)
	return &PodDisruptionBudgetManager{
		PodDisruptionBudgetInterface: clientset.PolicyV1beta1().PodDisruptionBudgets(namespace),
		namespace:                    namespace,
	}
}

type PodDisruptionBudgetManager struct {
	tV1Beta1.PodDisruptionBudgetInterface
	namespace string
}

func (d *PodDisruptionBudgetManager) Clone(toBeCloned string, cloneName string) {
	d.CloneInline(toBeCloned, cloneName, false)
}

func (d *PodDisruptionBudgetManager) CloneInline(toBeCloned string, cloneName string, _ bool) {
	pdb, err := d.Get(toBeCloned, metaV1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get latest version of Deployment: %v", err)
	}
	log.Printf("Cloning PodDisruptionBudget... %s -> %s\n", pdb.ObjectMeta.Name, cloneName)
	pdb.ObjectMeta.Name = cloneName
	pdb.ObjectMeta.Labels["app"] = cloneName
	pdb.Spec.Selector.MatchLabels["app"] = cloneName
	pdb.ResourceVersion = ""
	pdb.ObjectMeta.SelfLink = ""
	pdb.ObjectMeta.UID = ""
	pdb, err = d.PodDisruptionBudgetInterface.Create(pdb)

	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Cloned PodDisruptionBudget %q.\n", pdb.GetObjectMeta().GetName())
}

func (d *PodDisruptionBudgetManager) CreateExample(name string) {
	log.Fatalln("TODO Implement CreateExample")
}

func (d *PodDisruptionBudgetManager) Update(name string) {
	log.Fatalln("TODO Implement Update")
}

func (d *PodDisruptionBudgetManager) List() {
	log.Fatalln("TODO Implement List")
}

func (d *PodDisruptionBudgetManager) Delete(names ...string) {
	log.Fatalln("TODO Implement Delete")
}
