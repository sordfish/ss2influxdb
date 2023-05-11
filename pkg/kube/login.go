package kube

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Login() (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error

	// try to use the local kubeconfig file first
	kubeconfigPath := os.Getenv("KUBECONFIG")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Printf("Failed to load kubeconfig file at %s: %v\n", kubeconfigPath, err)
		}
	} else {

		kubeconfigPath = home + "/.kube/config"
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			log.Printf("Failed to load kubeconfig file at %s: %v\n", kubeconfigPath, err)
		}
	}

	// if the local kubeconfig file wasn't found, use the in-cluster config
	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("Failed to create in-cluster config: %v", err)
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
		return clientset, err
	}

	return clientset, err
}
