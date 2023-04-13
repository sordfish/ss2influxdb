package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Create a Kubernetes client using the in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Get the "default" namespace
	namespace := "sunsynk"

	// Create a new secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sunsynk-credentials",
			Namespace: namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"username": []byte("my-username"),
			"password": []byte("my-password"),
		},
	}

	// Create the secret
	result, err := clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())

	// Get the secret
	result, err = clientset.CoreV1().Secrets(namespace).Get(context.Background(), "sunsynk-credentials", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got secret %q.\n", result.GetObjectMeta().GetName())

	// Get the secret
	result, err = clientset.CoreV1().Secrets(namespace).Get(context.Background(), "my-secret", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got secret %q.\n", result.GetObjectMeta().GetName())

	// Update the secret
	result.Data["password"] = []byte("new-password")
	_, err = clientset.CoreV1().Secrets(namespace).Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated secret %q.\n", result.GetObjectMeta().GetName())

	// Delete the secret
	err = clientset.CoreV1().Secrets(namespace).Delete(context.Background(), "my-secret", metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted secret %q.\n", "my-secret")
}
