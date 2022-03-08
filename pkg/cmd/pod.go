package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/kube-resource/pkg/resource"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	KRPodExample = templates.Examples(`
	kubectl kr pod
	kubectl kr pod -l app=my-nginx
	kubectl kr pod -l app=my-nginx -o json
	kubectl kr pod -l app=my-nginx -n default -o yaml
	`)
)

func podCmd() *cobra.Command {
	o := resource.PodOption{}
	podCmd := &cobra.Command{
		Use:                   "pod [NAME | -l label]",
		Short:                 "pod provides an overview of the pod",
		DisableFlagsInUseLine: true,
		Example:               KRPodExample,
		Aliases:               []string{"pods", "po"},
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate()
			o.RunResourcePod()
		},
	}
	podCmd.PersistentFlags().StringVarP(&o.KubeCtx, "context", "", "", "context to use for Kubernetes config")
	podCmd.PersistentFlags().StringVarP(&o.KubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	podCmd.PersistentFlags().StringVarP(&o.Output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	podCmd.PersistentFlags().StringVarP(&o.LabelSelector, "label", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	podCmd.PersistentFlags().StringVarP(&o.FieldSelector, "field", "f", "", "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. -f key1=value1,key2=value2)")
	podCmd.PersistentFlags().StringVarP(&o.SortBy, "sortBy", "s", "cpu", "sort by cpu or memory")
	podCmd.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "only include rosource from this namespace")
	return podCmd
}

func init() {
	rootCmd.AddCommand(podCmd())
}
