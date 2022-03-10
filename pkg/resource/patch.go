package resource

import (
	"log"
	"os"
	"sort"
	"time"

	"github.com/gosuri/uitable"
	"github.com/ysicing/kube-resource/pkg/kube"
	"github.com/ysicing/kube-resource/pkg/output"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type PatchOption struct {
	Selector   string
	KubeCtx    string
	KubeConfig string
	Skip       int
}

func (p *PatchOption) Validate() {
}

func (p *PatchOption) RunResourcePatch() error {
	selector := labels.Everything()
	var err error
	if len(p.Selector) > 0 {
		selector, err = labels.Parse(p.Selector)
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
	log.Print("....")
	ds := k.GetDeploy(selector)
	sort.Sort(NewNodeMetricsSorter(ds, ""))
	table := uitable.New()
	table.AddRow("Name", "Namespace", "CPU分配", "内存分配", "存活时间")
	// count := 1
	for _, d := range ds {
		dc := d.Spec.Template.Spec.Containers[0]
		if dc.Resources.Requests.Cpu().AsApproximateFloat64()*1000.0 <= 100.0 || dc.Resources.Requests.Cpu().AsApproximateFloat64()*1000.0 >= 500.0 {
			continue
		}
		// if count > 10 {
		// 	break
		// }
		table.AddRow(d.Name, d.Namespace, dc.Resources.Requests.Cpu().String(), dc.Resources.Requests.Memory().String(), time.Since(d.CreationTimestamp.Time))
		// if err := k.Patch(d.Namespace, d.Name); err != nil {
		// 	log.Println(err)
		// 	break
		// }
		// count++
		// time.Sleep(time.Second * 2)
	}
	return output.EncodeTable(os.Stdout, table)
}

type DeploySorter struct {
	deploy []appsv1.Deployment
	sortBy string
}

func (n *DeploySorter) Len() int {
	return len(n.deploy)
}

func (n *DeploySorter) Swap(i, j int) {
	n.deploy[i], n.deploy[j] = n.deploy[j], n.deploy[i]
}

func (n *DeploySorter) Less(i, j int) bool {
	switch n.sortBy {
	case "cpu":
		return n.deploy[i].Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().MilliValue() > n.deploy[j].Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
	case "memory":
		return n.deploy[i].Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value() > n.deploy[j].Spec.Template.Spec.Containers[0].Resources.Requests.Memory().Value()
	default:
		return time.Since(n.deploy[i].CreationTimestamp.Time).Hours() < time.Since(n.deploy[j].CreationTimestamp.Time).Hours()
	}
}

func NewNodeMetricsSorter(deploys []appsv1.Deployment, sortBy string) *DeploySorter {
	return &DeploySorter{
		deploy: deploys,
		sortBy: sortBy,
	}
}
