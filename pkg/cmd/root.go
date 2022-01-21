package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ysicing/kube-resource/pkg/resource"
)

var kubeContext string
var kubeConfig string
var output string
var labels string
var namespace string

var rootCmd = &cobra.Command{
	Use:   "kube-resource",
	Short: "kube-resource provides an overview of the resource",
	Run: func(cmd *cobra.Command, args []string) {
		resource.FetchAndPrint(kubeContext, kubeConfig, namespace, labels, output)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&kubeContext, "context", "", "", "context to use for Kubernetes config")
	rootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	rootCmd.PersistentFlags().StringVarP(&labels, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "only include rosource from this namespace")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
