/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fvanmaldegem/k8sctx/pkg/helpers"
	"github.com/spf13/cobra"
)

var (
	user      bool
	cluster   bool
	full      bool
	removeCmd = &cobra.Command{
		Use:   "remove <context-name>",
		Short: "Remove a context.",
		Long:  `Removes a context. Pass the --user and --cluster flag to respectively delete this info as well. The can be captured together with the --full flag`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			remove(args[0], user, cluster, full)
		},
	}
)

func init() {
	removeCmd.Flags().BoolVar(&user, "user", false, "define if the user also needs to be removed")
	removeCmd.Flags().BoolVar(&cluster, "cluster", false, "define if the cluster also needs to be removed")
	removeCmd.Flags().BoolVar(&full, "full", false, "This will remove the user and the cluster together with the context")
	rootCmd.AddCommand(removeCmd)
}

func remove(contextName string, removeUser, removeCluster, removeFull bool) {
	newKubeconfig := Kubeconfig.DeepCopy()
	oldContext := Kubeconfig.Contexts[contextName]

	if helpers.PromptRemovalOf("context", Force) {
		fmt.Printf("deleting context '%s'\n", contextName)
		delete(newKubeconfig.Contexts, contextName)
	}

	if (removeUser || removeFull) && helpers.PromptRemovalOf("user", Force) {
		userN := contextName
		if oldContext != nil {
			userN = oldContext.AuthInfo
		}

		fmt.Printf("deleting user '%s'\n", userN)
		delete(newKubeconfig.AuthInfos, userN)
	}

	if (removeCluster || removeFull) && helpers.PromptRemovalOf("cluster", Force) {
		clusterN := contextName
		if oldContext != nil {
			clusterN = oldContext.Cluster
		}

		fmt.Printf("deleting cluster '%s'\n", clusterN)
		delete(newKubeconfig.Clusters, clusterN)
	}

	WriteKubeConfig(newKubeconfig)
}
