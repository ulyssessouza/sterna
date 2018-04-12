package managers

import (
	"fmt"
	"log"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tV1Beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	aV1Beta1 "k8s.io/api/apps/v1beta1"
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
	statefulSet.Status = aV1Beta1.StatefulSetStatus{}
	statefulSet.ObjectMeta.Name = cloneName
	statefulSet.ResourceVersion = ""
	statefulSet.ObjectMeta.SelfLink = ""
	statefulSet.ObjectMeta.UID = ""
	statefulSet.ObjectMeta.Labels["app"] = cloneName
	statefulSet.Spec.ServiceName = cloneName
	statefulSet.Spec.Selector.MatchLabels["app"] = cloneName
	statefulSet.Spec.Template.ObjectMeta.Labels["app"] = cloneName
	statefulSet.Spec.Template.Spec.InitContainers[0].Args[1] = fmt.Sprintf("-service=%s", cloneName)
	statefulSet.Spec.Template.Spec.Containers[0].Name = cloneName
	statefulSet.Spec.Template.Spec.Containers[0].Command[0] = "/bin/bash"
	statefulSet.Spec.Template.Spec.Containers[0].Command[1] = "-ecx"
	statefulSet.Spec.Template.Spec.Containers[0].Command[2] = fmt.Sprintf(`
# The use of qualified "hostname -f" is crucial:
# Other nodes aren't able to look up the unqualified hostname.
CRARGS=("start" "--logtostderr" "--insecure" "--host" "$(hostname -f)" "--http-host" "0.0.0.0")
# We only want to initialize a new cluster (by omitting the join flag)
# if we're sure that we're the first node (i.e. index 0) and that
# there aren't any other nodes running as part of the cluster that
# this is supposed to be a part of (which indicates that a cluster
# already exists and we should make sure not to create a new one).
# It's fine to run without --join on a restart if there aren't any
# other nodes.
if [ ! "$(hostname)" == "%s-0" ] || \
   [ -e "/cockroach/cockroach-data/cluster_exists_marker" ]
then
	# We don't join cockroachdb in order to avoid a node attempting
	# to join itself, which currently doesn't work
	# (https://github.com/cockroachdb/cockroach/issues/9625).
	CRARGS+=("--join" "%s-public")
fi
exec /cockroach/cockroach ${CRARGS[*]}
`, cloneName, cloneName)

	statefulSet.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].PodAffinityTerm.LabelSelector.MatchExpressions[0].Values[0] = cloneName

	cloneNewStatefulSet, err := d.StatefulSetInterface.Create(statefulSet)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Cloned StatefulSet %q.\n", cloneNewStatefulSet.GetObjectMeta().GetName())
}

func (d *StatefulSetManager) CreateExample(name string) {
	log.Fatalln("TODO Implement CreateExample for StatefulSet")
}

func (d *StatefulSetManager) Update(name string, _ string) {
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
