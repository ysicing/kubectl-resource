package kube

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/metricsutil"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsV1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type ClientConfig struct {
	QPS        float32
	Burst      int
	KubeCtx    string
	KubeConfig string
}

type KubeClient struct {
	apiClient     kubernetes.Interface
	metricsClient *metrics.Clientset
}

func NewKubeClient(cc *ClientConfig) (*KubeClient, error) {
	client, metricsClient, err := New(cc)
	if err != nil {
		return nil, err
	}

	return &KubeClient{
		apiClient:     client,
		metricsClient: metricsClient,
	}, nil
}

// New returns a kubernetes client.
// It tries first with in-cluster config, if it fails it will try with out-of-cluster config.
func New(cc *ClientConfig) (client kubernetes.Interface, metricsClient *metrics.Clientset, err error) {
	// client, err = NewInCluster(cc)
	// if err == nil {
	// 	return
	// }
	client, metricsClient, err = NewFromConfig(cc)
	if err != nil {
		return
	}

	return
}

// NewFromConfig returns a new out-of-cluster kubernetes client.
func NewFromConfig(cc *ClientConfig) (client kubernetes.Interface, metricsClient *metrics.Clientset, err error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if cc.KubeConfig == "" {
		cc.KubeConfig = os.Getenv("KUBECONFIG")
		if cc.KubeConfig == "" {
			dir, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, err
			}
			cc.KubeConfig = filepath.Join(dir, ".kube", "config")
		}
	}

	loadingRules.ExplicitPath = cc.KubeConfig

	var config *rest.Config

	if cc.KubeCtx != "" {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			loadingRules,
			&clientcmd.ConfigOverrides{CurrentContext: cc.KubeCtx},
		).ClientConfig()
	} else {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", cc.KubeConfig)
	}

	if err != nil {
		return
	}

	cc.apply(config)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = clientset

	metricsClient, err = metrics.NewForConfig(config)
	if err != nil {
		return
	}

	return client, metricsClient, nil
}

// NewInCluster returns a new in-cluster kubernetes client.
func NewInCluster(cc *ClientConfig) (client kubernetes.Interface, metricsClient *metrics.Clientset, err error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}

	cc.apply(config)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = clientset

	metricsClient, err = metrics.NewForConfig(config)
	if err != nil {
		return
	}

	return
}

func (cc *ClientConfig) apply(config *rest.Config) {
	if cc.QPS > 0.0 {
		config.QPS = cc.QPS // the default is rest.DefaultQPS which is 5.0
	}

	if cc.Burst > 0 {
		config.Burst = cc.Burst // the default is rest.DefaultBurst which is 10
	}
}

// GetNodes
func (k *KubeClient) GetNodes(resourceName string, selector labels.Selector) (map[string]corev1.Node, error) {
	nodes := make(map[string]corev1.Node)
	if len(resourceName) > 0 {
		node, err := k.apiClient.CoreV1().Nodes().Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		nodes[node.Name] = *node
	} else {
		nodeList, err := k.apiClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			return nil, err
		}
		for _, i := range nodeList.Items {
			nodes[i.Name] = i
		}
	}
	return nodes, nil
}

// GetActivePodByNodename
func (k *KubeClient) GetActivePodByNodename(node corev1.Node) (*corev1.PodList, error) {
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + node.Name +
		",status.phase!=" + string(corev1.PodSucceeded) +
		",status.phase!=" + string(corev1.PodFailed))

	if err != nil {
		return nil, err
	}
	activePods, err := k.apiClient.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{FieldSelector: fieldSelector.String()})
	if err != nil {
		return nil, err
	}
	return activePods, err
}

// GetActivePodByPodname
func (k *KubeClient) GetPodByPodname(podName string, namespace string) (*corev1.Pod, error) {
	pod, err := k.apiClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, err
}

type NodeResources struct {
	NodeName            string `json:"nodeName" yaml:"nodeName"`
	NodeIP              string `json:"nodeIP" yaml:"nodeIP"`
	CPUUsages           string `json:"cpuUsages" yaml:"cpuUsages"`
	CPURequests         string `json:"cpuRequests" yaml:"cpuRequests"`
	CPULimits           string `json:"cpuLimits" yaml:"cpuLimits"`
	CPUCapacity         string `json:"cpuCapacity" yaml:"cpuCapacity"`
	CPURequestsFraction string `json:"cpuRequestsFraction" yaml:"cpuRequestsFraction"`
	CPULimitsFraction   string `json:"cpuLimitsFraction" yaml:"cpuLimitsFraction"`

	MemoryUsages           string `json:"memoryUsages" yaml:"memoryUsages"`
	MemoryRequests         string `json:"memoryRequests" yaml:"memoryRequests"`
	MemoryLimits           string `json:"memoryLimits" yaml:"memoryLimits"`
	MemoryCapacity         string `json:"memoryCapacity" yaml:"memoryCapacity"`
	MemoryRequestsFraction string `json:"memoryRequestsFraction" yaml:"memoryRequestsFraction"`
	MemoryLimitsFraction   string `json:"memoryLimitsFraction" yaml:"memoryLimitsFraction"`

	AllocatedPods int    `json:"allocatedPods" yaml:"allocatedPods"`
	PodCapacity   int64  `json:"podCapacity" yaml:"podCapacity"`
	PodFraction   string `json:"podFraction" yaml:"podFraction"`

	Age string `json:"age" yaml:"age"`
}

// NodeResources
func (k *KubeClient) GetNodeResources(sortBy string, selector labels.Selector) ([]NodeResources, error) {
	//resources := make(map[string]map[string]interface{})
	var resources []NodeResources
	var nodenames []string

	metrics, err := k.GetNodeMetricsFromMetricsAPI("", selector)
	if err != nil {
		return nil, err
	}
	//判断是否排序
	if len(sortBy) > 0 {
		sort.Sort(metricsutil.NewNodeMetricsSorter(metrics.Items, sortBy))
	}
	for _, i := range metrics.Items {
		nodenames = append(nodenames, i.Name)
	}

	nodes, err := k.GetNodes("", selector)
	if err != nil {
		return nil, err
	}

	for _, nodename := range nodenames {
		//resource := make(map[string]interface{})
		var resource NodeResources
		activePodsList, err := k.GetActivePodByNodename(nodes[nodename])
		if err != nil {
			return nil, err
		}
		NodeMetricsList, err := k.GetNodeMetricsFromMetricsAPI(nodename, selector)
		if err != nil {
			return nil, err
		}

		resource.NodeName = nodename
		resource.NodeIP = nodes[nodename].Status.Addresses[0].Address
		resource.Age = time.Since(nodes[nodename].CreationTimestamp.Time).String()
		noderesource, err := getNodeAllocatedResources(nodes[nodename], activePodsList, NodeMetricsList)
		if err != nil {
			log.Printf("Couldn't get allocated resources of %s node: %s\n", nodename, err)
		}
		resource.CPUUsages = noderesource.CPUUsages.String()
		resource.CPURequests = noderesource.CPURequests.String()
		resource.CPULimits = noderesource.CPULimits.String()
		resource.CPUCapacity = noderesource.CPUCapacity.String()
		resource.CPURequestsFraction = ExceedsCompare(float64ToString(noderesource.CPURequestsFraction))
		resource.CPULimitsFraction = float64ToString(noderesource.CPULimitsFraction)

		resource.MemoryUsages = noderesource.MemoryUsages.String()
		resource.MemoryRequests = noderesource.MemoryRequests.String()
		resource.MemoryLimits = noderesource.MemoryLimits.String()
		resource.MemoryCapacity = noderesource.MemoryCapacity.String()
		resource.MemoryRequestsFraction = ExceedsCompare(float64ToString(noderesource.MemoryRequestsFraction))
		resource.MemoryLimitsFraction = float64ToString(noderesource.MemoryLimitsFraction)

		resource.AllocatedPods = noderesource.AllocatedPods
		resource.PodCapacity = noderesource.PodCapacity
		resource.PodFraction = ExceedsCompare(float64ToString(noderesource.PodFraction))
		resources = append(resources, resource)
	}
	return resources, err
}

type PodsResources struct {
	Name                 string `json:"name" yaml:"name"`
	Namespace            string `json:"namespace" yaml:"namespace"`
	CPUUsages            string `json:"cpuUsages" yaml:"cpuUsages"`
	CPURequests          string `json:"cpuRequests" yaml:"cpuRequests"`
	CPULimits            string `json:"cpuLimits" yaml:"cpuLimits"`
	CPUUsagesFraction    string `json:"cpuUsagesFraction" yaml:"cpuUsagesFraction"`
	MemoryUsages         string `json:"memoryUsages" yaml:"memoryUsages"`
	MemoryRequests       string `json:"memoryRequests" yaml:"memoryRequests"`
	MemoryLimits         string `json:"memoryLimits" yaml:"memoryLimits"`
	MemoryUsagesFraction string `json:"memoryUsagesFraction" yaml:"memoryUsagesFraction"`
}

func (k *KubeClient) GetPodResources(podmetrics []metricsapi.PodMetrics, namespace string, sortBy string) ([]PodsResources, error) {
	var resources []PodsResources

	//判断是否排序
	if len(sortBy) > 0 {
		allNamespaces := true
		if len(namespace) > 0 {
			allNamespaces = false
		}
		sort.Sort(metricsutil.NewPodMetricsSorter(podmetrics, allNamespaces, sortBy))
	}
	for _, podmetric := range podmetrics {
		var resource PodsResources
		pod, err := k.GetPodByPodname(podmetric.Name, podmetric.Namespace)
		if err != nil {
			return nil, err
		}

		resource.Name = podmetric.Name
		resource.Namespace = podmetric.Namespace
		podresource, err := getPodAllocatedResources(pod, &podmetric)
		if err != nil {
			return nil, err
		}

		resource.CPUUsages = podresource.CPUUsages.String()
		resource.CPUUsagesFraction = ExceedsCompare(float64ToString(podresource.CPUUsagesFraction))
		resource.CPURequests = podresource.CPURequests.String()
		resource.CPULimits = podresource.CPULimits.String()

		resource.MemoryUsages = podresource.MemoryUsages.String()
		resource.MemoryUsagesFraction = ExceedsCompare(float64ToString(podresource.MemoryUsagesFraction))
		resource.MemoryRequests = podresource.MemoryRequests.String()
		resource.MemoryLimits = podresource.MemoryLimits.String()
		resources = append(resources, resource)
	}

	return resources, nil
}

// PodMetricses returns all pods' usage metrics
func (k *KubeClient) PodMetricses() (*metricsV1beta1api.PodMetricsList, error) {
	podMetricses, err := k.metricsClient.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return podMetricses, nil
}

// GetNodeMetricsFromMetricsAPI
func (k *KubeClient) GetNodeMetricsFromMetricsAPI(resourceName string, selector labels.Selector) (*metricsapi.NodeMetricsList, error) {
	var err error
	versionedMetrics := &metricsV1beta1api.NodeMetricsList{}
	mc := k.metricsClient.MetricsV1beta1()
	nm := mc.NodeMetricses()
	if resourceName != "" {
		m, err := nm.Get(context.TODO(), resourceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		versionedMetrics.Items = []metricsV1beta1api.NodeMetrics{*m}
	} else {
		versionedMetrics, err = nm.List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			return nil, err
		}
	}
	metrics := &metricsapi.NodeMetricsList{}

	err = metricsV1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// GetPodMetricsFromMetricsAPI
func (k *KubeClient) GetPodMetricsFromMetricsAPI(namespace string, labelSelector labels.Selector, fieldSelector fields.Selector) (*metricsapi.PodMetricsList, error) {
	var err error
	ns := metav1.NamespaceAll
	if len(namespace) > 0 {
		ns = namespace
	}

	versionedMetrics := &metricsV1beta1api.PodMetricsList{}
	versionedMetrics, err = k.metricsClient.MetricsV1beta1().PodMetricses(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector.String(), FieldSelector: fieldSelector.String()})
	if err != nil {
		return nil, err
	}
	metrics := &metricsapi.PodMetricsList{}
	err = metricsV1beta1api.Convert_v1beta1_PodMetricsList_To_metrics_PodMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
