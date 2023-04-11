package main

import (
	"context"
	"flag"
	"fmt"
	myv1alpha1 "github.com/obaydullahmhs/crd/pkg/apis/aadee.apps/v1alpha1"
	klientkluster "github.com/obaydullahmhs/crd/pkg/client/clientset/versioned"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientcluster "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/code-generator"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("Building config from flags failed, %s, trying to build inclusterconfig", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error %s building inclusterconfig", err.Error())
		}
	}

	clientset, err := clientcluster.NewForConfig(config)
	if err != nil {
		log.Printf("getting client set %s\n", err.Error())
	}
	klientset, err := klientkluster.NewForConfig(config)
	if err != nil {
		log.Printf("getting klient set %s\n", err.Error())
	}

	customCRD := v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "klusters.aadee.apps",
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: "aadee.apps",
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:    "v1alpha1",
					Served:  true,
					Storage: true,
					Schema: &v1.CustomResourceValidation{
						OpenAPIV3Schema: &v1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]v1.JSONSchemaProps{
										"name": {
											Type: "string",
										},
										"replicas": {
											Type: "integer",
										},
										"container": {
											Type: "object",
											Properties: map[string]v1.JSONSchemaProps{
												"image": {
													Type: "string",
												},
												"port": {
													Type: "integer",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Scope: "Namespaced",
			Names: v1.CustomResourceDefinitionNames{
				Kind:     "Kluster",
				Plural:   "klusters",
				Singular: "kluster",
				ShortNames: []string{
					"kl",
				},
				Categories: []string{
					"all",
				},
			},
		},
	}
	ctx := context.TODO()
	// delete cr

	_ = clientset.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, customCRD.Name, metav1.DeleteOptions{})
	//creating a new one
	_, err = clientset.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, &customCRD, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating new cr %s\n", err.Error())
	}

	time.Sleep(5 * time.Second)
	log.Println("CRD is Created!")

	log.Println("Press ctrl+c to create a Kluster")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	log.Println("Creating Kluster")
	klObj := &myv1alpha1.Kluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kluster",
		},
		Spec: myv1alpha1.KlusterSpec{
			Name:     "CustomResource",
			Replicas: intptr(2),
			Container: myv1alpha1.ContainerSpec{
				Image: "obaydullahmhs/api-server",
				Port:  3000,
			},
		},
	}
	_, err = klientset.AadeeV1alpha1().Klusters("default").Create(ctx, klObj, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error creating kluster %s\n", err.Error())
	}
	time.Sleep(2 * time.Second)
	log.Println("Kluster created!")

	klusters, err := klientset.AadeeV1alpha1().Klusters("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("listing klusters %s\n", err.Error())
	}
	fmt.Printf("length of klusters is %d and name is %s\n", len(klusters.Items), klusters.Items[0].Name)

	log.Println("Press ctrl+c to clean up!")
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	err = klientset.AadeeV1alpha1().Klusters("default").Delete(ctx, klObj.Name, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}

	err = clientset.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, customCRD.Name, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}

	log.Println("cleaned up")

}
func intptr(i int32) *int32 {
	return &i
}
