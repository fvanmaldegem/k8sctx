/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fvanmaldegem/k8sctx/pkg/helpers"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	newKubeconfigFile string
	userName          string
	clusterName       string
	oldContextName    string
	newContextName    string

	addCmd = &cobra.Command{
		Use:   "add <context-name>",
		Short: "Add a new context",
		Long:  `Add a new context to your config file with a specified name`,
		Args:  cobra.ExactArgs(1),
		Run:   add,
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&newKubeconfigFile, "newKubeconfig", "", "The kubeconfig to copy the context out of. You can leave this empty to open an editor.")
	addCmd.Flags().StringVar(&oldContextName, "context", "", "The name of the old context to copy into the new one. If this is passed, it wil ignore the user and cluster flag.")
	addCmd.Flags().StringVar(&userName, "user", "", "The name of the user to copy into the new context. If a context is given, it will select from there. If there is only one it will default to that one")
	addCmd.Flags().StringVar(&clusterName, "cluster", "", "The name of the cluster to copy into the new context. If a context is given, it will select from there. If there is only one it will default to that one")
}

func add(cmd *cobra.Command, args []string) {
	newContextName = args[0]
	oldConfig := loadConfig()

	if oldContextName != "" {
		context := getFlag(oldContextName, oldConfig.Contexts, "the old name of the context")
		userName = context.AuthInfo
		clusterName = context.Cluster
	}

	user := getFlag(userName, oldConfig.AuthInfos, "the name of the user int the old config")
	cluster := getFlag(clusterName, oldConfig.Clusters, "the name of the cluster in the old config")

	appendToConfig(user, cluster, newContextName)
}

func appendToConfig(user *api.AuthInfo, cluster *api.Cluster, newContextname string) {
	for newNameExists(newContextname) {
		s := strings.Split(newContextName, "-")
		l := len(s) - 1
		i, err := strconv.ParseInt(s[l], 10, 0)
		if err != nil {
			newContextname = fmt.Sprintf("%s-%d", newContextName, 1)
			continue
		}
		newContextname = fmt.Sprintf("%s-%d", strings.Join(s[0:l], "-"), i)
	}
	newContext := api.NewContext()
	newContext.AuthInfo = newContextname
	newContext.Cluster = newContextname

	Kubeconfig.AuthInfos[newContextname] = user.DeepCopy()
	Kubeconfig.Clusters[newContextname] = cluster.DeepCopy()
	Kubeconfig.Contexts[newContextname] = newContext.DeepCopy()

	WriteKubeConfig(Kubeconfig)
	fmt.Printf("Added context '%s'", newContextname)
}

func newNameExists(n string) bool {
	if Force {
		return false
	}

	isFree := Kubeconfig.AuthInfos[n] == nil &&
		Kubeconfig.Clusters[n] == nil &&
		Kubeconfig.Contexts[n] == nil

	if !isFree {
		return !helpers.PromptConfirm(fmt.Sprintf("Do you want to override '%s'", n), Force)
	}

	return false
}

func getFlag[T any](flag string, m map[string]*T, label string) *T {
	items := maps.Keys(m)

	if flag == "" {
		if len(items) == 1 {
			flag = items[0]
		} else {
			flag, _ = promptNames(items, label)
		}
	}

	r := m[flag]
	if r == nil {
		flag = items[0]
		r = m[flag]
	}

	return r
}

func loadConfig() *api.Config {
	config := &api.Config{}
	if newKubeconfigFile == "" {
		config = helpers.OpenAndRead()
	} else {
		config = clientcmd.GetConfigFromFileOrDie(KubeconfigFile)
	}

	return config
}

func promptNames(items []string, label string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, v, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return v, nil
}
