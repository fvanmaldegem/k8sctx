/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list <context-name>",
		Short: "List all contexts or list information about a specific context.",
		Long:  `List all contexts or list information about a specific context. The active context will be marked with a '*'.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				list(args[0])
			} else {
				listAll()
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func listAll() {
	for n := range Kubeconfig.Contexts {
		l := n
		if Kubeconfig.CurrentContext == n {
			l = fmt.Sprintf("%s *", l)
		}
		fmt.Println(l)
	}
}

func list(contextName string) {
	context := Kubeconfig.Contexts[contextName]
	if context == nil {
		fmt.Printf("context '%s' was not found in kubeconfig '%s", contextName, KubeconfigFile)
		os.Exit(1)
	}

	fmt.Printf(`Name: %s
Cluster: %s
User: %s
`, contextName, context.Cluster, context.AuthInfo)

}
