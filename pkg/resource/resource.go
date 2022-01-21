package resource

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/ysicing/kube-resource/pkg/kube"
	"github.com/ysicing/kube-resource/pkg/output"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Res struct {
	Namespace string                      `json:"namespace" yaml:"namespace"`
	Name      string                      `json:"name" yaml:"name"`
	Type      string                      `json:"type" yaml:"type"`
	Resources corev1.ResourceRequirements `json:"resources" yaml:"resources"`
}

func FetchAndPrint(kubeContext, kubeConfig, namespace, labels, format string) error {
	kubecfg := &kube.ClientConfig{
		KubeCtx:    kubeContext,
		KubeConfig: kubeConfig,
	}
	client, err := kube.New(kubecfg)
	if err != nil {
		fmt.Printf("Error connecting to Kubernetes: %v\n", err)
		return err
	}
	deployList := getDeploy(client, namespace, labels)

	switch strings.ToLower(format) {
	case "json":
		return output.EncodeJSON(os.Stdout, deployList)
	case "yaml":
		return output.EncodeYAML(os.Stdout, deployList)
	default:
		table := uitable.New()
		table.AddRow("NAMESPACE", "NAME", "TYPE", "CPU REQUESTS", "CPU LIMITS", "MEMORY REQUESTS", "MEMORY LIMITS")
		for _, d := range deployList {
			table.AddRow(d.Namespace, d.Name, d.Type, d.Resources.Requests.Cpu(), d.Resources.Limits.Cpu(), d.Resources.Requests.Memory(), d.Resources.Limits.Memory())
		}
		return output.EncodeTable(os.Stdout, table)
	}
}

func getDeploy(clientset kubernetes.Interface, namespace, label string) []Res {
	deploys, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		return nil
	}
	var res []Res
	for _, d := range deploys.Items {
		r := Res{
			Namespace: d.Namespace,
			Name:      d.Name,
			Type:      "Deployment",
		}
		var rc, rm, lc, lm int64
		for _, c := range d.Spec.Template.Spec.Containers {
			rc = rc + c.Resources.Requests.Cpu().MilliValue()*int64(*d.Spec.Replicas)
			rm = rm + c.Resources.Requests.Memory().Value()*int64(*d.Spec.Replicas)
			lc = lc + c.Resources.Limits.Cpu().MilliValue()*int64(*d.Spec.Replicas)
			lm = lm + c.Resources.Limits.Memory().Value()*int64(*d.Spec.Replicas)
		}
		r.Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(lc, resource.DecimalSI),
				"memory": *resource.NewQuantity(lm, resource.BinarySI),
			},
			Requests: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(rc, resource.DecimalSI),
				"memory": *resource.NewQuantity(rm, resource.BinarySI),
			},
		}
		res = append(res, r)
	}
	return res
}
