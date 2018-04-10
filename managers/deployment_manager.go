package managers

import (
	appsV1 "k8s.io/api/apps/v1"
	apiV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedV1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"fmt"
	"k8s.io/client-go/util/retry"
	"log"
)

func NewDeploymentManager(kubeconfig string, namespace string) *DeploymentManager {
	clientset := getClientSet(kubeconfig)
	return &DeploymentManager{DeploymentInterface: clientset.AppsV1().Deployments(apiV1.NamespaceDefault), namespace: namespace}
}

type DeploymentManager struct {
	typedV1.DeploymentInterface
	namespace string
}

func (d *DeploymentManager) CreateExample(name string) {
	deployment := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiV1.PodSpec{
					Containers: []apiV1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiV1.ContainerPort{
								{
									Name:          "nginx-httpport",
									Protocol:      apiV1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
					RestartPolicy: "Always",
				},
			},
		},
	}

	log.Println("Creating deployment...")
	result, err := d.DeploymentInterface.Create(deployment)
	if err != nil {
		panic(err)
	}
	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func (d *DeploymentManager) Clone(toBeCloned string, cloneName string) {
	resultDeployment, getErr := d.DeploymentInterface.Get(toBeCloned, metaV1.GetOptions{})
	if getErr != nil {
		log.Fatalf("Failed to get latest version of Deployment: %v", getErr)
	}

	log.Printf("Cloning deployment... %s -> %s\n", resultDeployment.ObjectMeta.Name, cloneName)
	resultDeployment.ObjectMeta.Name = cloneName
	resultDeployment.ResourceVersion = ""
	cloneDeployment, err := d.DeploymentInterface.Create(resultDeployment)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Cloned deployment %q.\n", cloneDeployment.GetObjectMeta().GetName())
}

func (d *DeploymentManager) Update(name string) {
	log.Println("Updating deployment...")
	//    You have two options to Update() this Deployment:
	//
	//    1. Modify the "deployment" variable and call: Update(deployment).
	//       This works like the "kubectl replace" command and it overwrites/loses changes
	//       made by other clients between you Create() and Update() the object.
	//    2. Modify the "result" returned by Get() and retry Update(result) until
	//       you no longer get a conflict error. This way, you can preserve changes made
	//       by other clients between Create() and Update(). This is implemented below
	//			 using the retry utility package included with client-go. (RECOMMENDED)
	//
	// More Info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#concurrency-control-and-consistency

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := d.DeploymentInterface.Get(name, metaV1.GetOptions{})
		if getErr != nil {
			log.Fatalf("Failed to get latest version of Deployment: %v", getErr)
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := d.DeploymentInterface.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	log.Println("Updated deployment...")
}

func (d *DeploymentManager) List() {
	log.Printf("Listing deployments in namespace %q:\n", d.namespace)
	list, err := d.DeploymentInterface.List(metaV1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		log.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func (d *DeploymentManager) Delete(names ...string) {
	log.Println("Deleting deployment...")
	deletePolicy := metaV1.DeletePropagationForeground
	for _, name := range names {
		if err := d.DeploymentInterface.Delete(name, &metaV1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		log.Printf("Deleted deployment '%s'\n", name)
	}

	log.Printf("Delete deployment finished with success\n")
}

func int32Ptr(i int32) *int32 { return &i }
