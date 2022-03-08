package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/kube-resource/pkg/resource"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	KRNodeExample = templates.Examples(`
	kubectl kr node
	`)
)

func nodeCmd() *cobra.Command {
	o := resource.NodeOption{}
	nodeCmd := &cobra.Command{
		Use:                   "node",
		DisableFlagsInUseLine: true,
		Short:                 "node provides an overview of the node",
		Aliases:               []string{"nodes", "no"},
		Example:               KRNodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate()
			o.RunResourceNode()
		},
	}
	nodeCmd.PersistentFlags().StringVarP(&o.KubeCtx, "context", "", "", "context to use for Kubernetes config")
	nodeCmd.PersistentFlags().StringVarP(&o.KubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	nodeCmd.PersistentFlags().StringVarP(&o.Output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	nodeCmd.PersistentFlags().StringVarP(&o.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	nodeCmd.PersistentFlags().StringVarP(&o.SortBy, "sortBy", "s", "cpu", "sort by cpu or memory")
	return nodeCmd
}

func init() {
	rootCmd.AddCommand(nodeCmd())
}
