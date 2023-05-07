package kube

import (
	"context"
	"encoding/json"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetK8sSecret(clientset *kubernetes.Clientset, secret, namespace string) (*corev1.Secret, error) {
	result, err := clientset.CoreV1().Secrets(namespace).Get(context.Background(), secret, metav1.GetOptions{})
	if err != nil {
		return result, err
	}
	return result, err
}

func UpdateK8sSecret(clientset *kubernetes.Clientset, result *corev1.Secret, namespace string, data interface{}, filename string) (*corev1.Secret, error) {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	result.Data[filename] = dataBytes

	update, err := clientset.CoreV1().Secrets(namespace).Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return update, err
}

func CreateK8sSecret(clientset *kubernetes.Clientset, secret, namespace string, data interface{}, filename string) (*corev1.Secret, error) {

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	NewSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret,
			Namespace: namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			filename: dataBytes,
		},
	}

	result, err := clientset.CoreV1().Secrets(namespace).Create(context.Background(), NewSecret, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return result, err

}
