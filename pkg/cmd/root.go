/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	// Force flag
	Force bool
	// in-memory kubeconfig file
	Kubeconfig *api.Config
	// absolute location of the kubeconfig file
	KubeconfigFile string
	// root command
	rootCmd = &cobra.Command{
		Use:   "k8sctx",
		Short: "The K8S Context creator",
		Long: `k8sctx is a command line application to properly manage your Kubernetes Contexts.

		Author: Floris van Maldegem`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initKubeConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&KubeconfigFile, "kubeconfig", "c", "", "The kubeconfig file to use (defaults to $HOME/.kube/config)")
	rootCmd.PersistentFlags().BoolVarP(&Force, "force", "f", false, "Specifies if the action needs to be forced. Will perform action without overwriting")

}

func initKubeConfig() {
	var err error
	if KubeconfigFile == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("cant find homefolder")
			os.Exit(1)
		}
		KubeconfigFile = path.Join(homedir, ".kube", "config")
	}
	Kubeconfig, err = clientcmd.LoadFromFile(KubeconfigFile)
	if err != nil {
		fmt.Printf("cant read '%s'\n", KubeconfigFile)
		os.Exit(1)
	}
}

func WriteKubeConfig(newKubeconfig *api.Config) {
	clientcmd.WriteToFile(*newKubeconfig, KubeconfigFile)
	fmt.Printf("saved config to: '%s'\n", KubeconfigFile)
}
