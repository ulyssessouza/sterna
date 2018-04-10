package managers

import (
	"log"
	coreV1 "k8s.io/api/core/v1"
	tCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewServiceManager(kubeconfig string, namespace string) *ServiceManager {
	clientset := getClientSet(kubeconfig)
	return &ServiceManager{ServiceInterface: clientset.CoreV1().Services(namespace), namespace: namespace}
}

type ServiceManager struct {
	tCoreV1.ServiceInterface
	namespace        string
}

func (d *ServiceManager) Clone(toBeCloned string, cloneName string) {
	service, err := d.Get(toBeCloned, metaV1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get latest version of Deployment: %v", err)
	}
	log.Printf("Cloning service... %s -> %s\n", service.ObjectMeta.Name, cloneName)
	service.ObjectMeta.Name = cloneName
	service.ResourceVersion = ""
	cloneService, err := d.ServiceInterface.Create(service)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Cloned service %q.\n", cloneService.GetObjectMeta().GetName())
}

func (d *ServiceManager) CreateExample(name string) {
	service := &coreV1.Service{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Labels: map[string]string{
				"service": name,
			},
			Name: name,
		},
		Spec: coreV1.ServiceSpec{
			Ports: []coreV1.ServicePort{
				{
					Port:     80,
					NodePort: 30080,
					Protocol: coreV1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"service": name,
			},
			Type: "NodePort",
		},
	}

	log.Println("Creating service...")
	result, err := d.ServiceInterface.Create(service)
	if err != nil {
		panic(err)
	}
	log.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
}

func (d *ServiceManager) Update(name string) {
	log.Fatalln("TODO Implement Update")
}

func (d *ServiceManager) List() {
	services, err := d.ServiceInterface.List(metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Printf("There are %d pods in the cluster\n", len(services.Items))

	for _, s := range services.Items {
		log.Printf("Service: %s\n", s.Name)
		for p, _ := range s.Spec.Ports {
			log.Println("- Port:", s.Spec.Ports[p].Port)
			log.Println("- NodePort:", s.Spec.Ports[p].NodePort)
		}
	}
}

func (d *ServiceManager) Delete(names ...string) {
	log.Println("Deleting service...")
	deletePolicy := metaV1.DeletePropagationForeground
	for _, name := range names {
		if err := d.ServiceInterface.Delete(name, &metaV1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		log.Printf("Deleted service '%s'\n", name)
	}

	log.Printf("Delete service finished with success\n")
}