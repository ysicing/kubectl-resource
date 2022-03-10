package kube

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *KubeClient) GetDeploy(selector labels.Selector) []appsv1.Deployment {
	listopt := metav1.ListOptions{}
	if selector != nil {
		listopt.LabelSelector = selector.String()
	}

	deploys, err := k.apiClient.AppsV1().Deployments("").List(context.TODO(), listopt)
	if err != nil {
		return nil
	}
	return deploys.Items
}

func (k *KubeClient) Patch(ns, name string) error {
	d, err := k.apiClient.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	d.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("100m"),
		corev1.ResourceMemory: resource.MustParse("256Mi"),
	}
	_, err = k.apiClient.AppsV1().Deployments(ns).Update(context.TODO(), d, metav1.UpdateOptions{})
	return err
}
