package managers

import (
	"log"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tV1Beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
)

func NewStatefulSetManager(kubeconfig string, namespace string) *StatefulSetManager {
	clientset := getClientSet(kubeconfig)
	return &StatefulSetManager{StatefulSetInterface: clientset.AppsV1beta1().StatefulSets(namespace), namespace: namespace}
}

type StatefulSetManager struct {
	tV1Beta1.StatefulSetInterface
	namespace string
}

func (d *StatefulSetManager) Clone(toBeCloned string, cloneName string) {
	d.CloneInline(toBeCloned, cloneName, false)
}

func (d *StatefulSetManager) CloneInline(toBeCloned string, cloneName string, _ bool) {
	statefulSet, err := d.Get(toBeCloned, metaV1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get latest version of StatefulSet: %v", err)
	}
	log.Printf("Cloning StatefulSet... %s -> %s\n", statefulSet.ObjectMeta.Name, cloneName)
	statefulSet.ObjectMeta.Name = cloneName
	statefulSet.ResourceVersion = ""
	statefulSet.Spec.Template.ObjectMeta.Labels["app"] = cloneName
	cloneService, err := d.StatefulSetInterface.Create(statefulSet)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Cloned StatefulSet %q.\n", cloneService.GetObjectMeta().GetName())
}

func (d *StatefulSetManager) CreateExample(name string) {
	log.Fatalln("TODO Implement CreateExample for StatefulSet")
}

func (d *StatefulSetManager) Update(name string) {
	log.Fatalln("TODO Implement Update for StatefulSet")
}

func (d *StatefulSetManager) List() {
	log.Fatalln("TODO Implement List for StatefulSet")
}

func (d *StatefulSetManager) Delete(names ...string) {
	log.Println("Deleting service...")
	deletePolicy := metaV1.DeletePropagationForeground
	for _, name := range names {
		if err := d.StatefulSetInterface.Delete(name, &metaV1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		log.Printf("Deleted StatefulSet '%s'\n", name)
	}
	log.Printf("Delete StatefulSet finished with success\n")
}
