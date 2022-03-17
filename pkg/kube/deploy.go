package kube

import (
	"context"
	"log"

	"github.com/ergoapi/util/ptr"
	"github.com/golang-module/carbon/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
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

func (k *KubeClient) GetDeployByName(ns, name string) (*appsv1.Deployment, error) {
	d, err := k.apiClient.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (k *KubeClient) PatchReplicas(ns, name string, nostop bool) error {
	d, err := k.GetDeployByName(ns, name)
	if err != nil {
		return err
	}
	ct := d.CreationTimestamp.Time

	if !nostop {
		y := carbon.Yesterday().Carbon2Time()
		if ct.Before(y) {
			log.Println("创建时间在昨天之前, 副本不改变")
			d.Spec.Replicas = ptr.Int32Ptr(0)
		}
	}
	if d.Spec.Template.Spec.NodeSelector == nil {
		log.Println("没有设置节点选择器, 设置")
		d.Spec.Template.Spec.NodeSelector = map[string]string{
			"k8s.easycorp.work/pool": "free",
		}
		d.Spec.Template.Spec.Tolerations = []corev1.Toleration{
			{
				Operator: "Exists",
			},
		}
	}
	if d.Spec.Template.Spec.Containers[0].ReadinessProbe == nil {
		log.Println("没有检测到ReadinessProbe, 添加")
		d.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/misc-status",
					Port: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      2,
			PeriodSeconds:       5,
			FailureThreshold:    10,
		}
	} else {
		if d.Spec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds >= 25 {
			log.Println("ReadinessProbe的InitialDelaySeconds大于25，改为25")
			d.Spec.Template.Spec.Containers[0].ReadinessProbe.InitialDelaySeconds = 5
			d.Spec.Template.Spec.Containers[0].ReadinessProbe.FailureThreshold = 10
			d.Spec.Template.Spec.Containers[0].ReadinessProbe.PeriodSeconds = 5
		}
	}
	_, err = k.apiClient.AppsV1().Deployments(ns).Update(context.TODO(), d, metav1.UpdateOptions{})
	log.Println("更新副本数")
	return err
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
