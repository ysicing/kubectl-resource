package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/kube-resource/pkg/resource"
)

func patchCmd() *cobra.Command {
	o := resource.PatchOption{}
	p := &cobra.Command{
		Use:   "patch",
		Short: "patch provides an overview of the patch",
		RunE: func(cmd *cobra.Command, args []string) error {
			// o.Validate()
			return o.RunNodePatch()
		},
	}
	p.PersistentFlags().StringVarP(&o.KubeCtx, "context", "", "", "context to use for Kubernetes config")
	p.PersistentFlags().StringVarP(&o.KubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	// p.PersistentFlags().StringVarP(&o.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	// p.PersistentFlags().IntVarP(&o.Skip, "skip", "", 100, "Skip the first N CPU")
	p.PersistentFlags().StringVarP(&o.Name, "name", "n", "", "name")
	p.PersistentFlags().StringVarP(&o.Namespace, "namespace", "", "", "namespace")
	p.PersistentFlags().BoolVarP(&o.NoStop, "nostop", "", false, "no stop")
	return p
}

func init() {
	rootCmd.AddCommand(patchCmd())
}
