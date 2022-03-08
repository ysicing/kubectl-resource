package resource

import (
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/ysicing/kube-resource/pkg/kube"
	"github.com/ysicing/kube-resource/pkg/output"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type PodOption struct {
	Namespace     string
	LabelSelector string
	FieldSelector string
	SortBy        string
	QPS           float32
	Burst         int
	KubeCtx       string
	KubeConfig    string
	Output        string
}

func (p *PodOption) Validate() {
	if len(p.SortBy) > 0 {
		if p.SortBy != "cpu" {
			p.SortBy = "memory"
		}
	}
}

func (p *PodOption) RunResourcePod() error {
	labelSelector := labels.Everything()
	var err error
	if len(p.LabelSelector) > 0 {
		labelSelector, err = labels.Parse(p.LabelSelector)
		if err != nil {
			return err
		}
	}
	fieldSelector := fields.Everything()
	if len(p.FieldSelector) > 0 {
		fieldSelector, err = fields.ParseSelector(p.FieldSelector)
		if err != nil {
			return err
		}
	}
	cfg := kube.ClientConfig{
		KubeCtx:    p.KubeCtx,
		KubeConfig: p.KubeConfig,
	}
	k, err := kube.NewKubeClient(&cfg)
	if err != nil {
		return err
	}
	metrics, err := k.GetPodMetricsFromMetricsAPI(p.Namespace, labelSelector, fieldSelector)
	if err != nil {
		return err
	}
	if len(metrics.Items) == 0 {
		return nil
	}
	data, err := k.GetPodResources(metrics.Items, p.Namespace, p.SortBy)
	if err != nil {
		return err
	}
	switch strings.ToLower(p.Output) {
	case "json":
		return output.EncodeJSON(os.Stdout, data)
	case "yaml":
		return output.EncodeYAML(os.Stdout, data)
	default:
		table := uitable.New()
		table.AddRow("Name", "IP", "CPU使用", "CPU分配", "CPU限制", "CPU容量", "内存使用", "内存分配", "内存限制", "内存容量", "pod数", "pod容量")
		for _, d := range data {
			table.AddRow(d)
		}
		return output.EncodeTable(os.Stdout, table)
	}
}
