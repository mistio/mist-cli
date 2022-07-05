package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jmespath/go-jmespath"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func modifyKubeconfigPrompt() bool {
	prompt := promptui.Select{
		Label: "Choose action",
		Items: []string{"Add cluster to the local kubeconfig", "Just print the cluster's kubeconfig"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		logger.Fatalf("Prompt failed %v\n", err)
	}
	return result == prompt.Items.([]string)[0]
}

type clusterInfo struct {
	name   string
	host   string
	port   string
	caCert string
}

func parseClusterResponse(decoded interface{}) (clusterInfo, error) {
	newClusterInfo := clusterInfo{}
	nameInterface, err := jmespath.Search("data.name", decoded)
	name, ok := nameInterface.(string)
	if err != nil || !ok {
		return clusterInfo{}, fmt.Errorf("error parsing cluster's name")
	}
	newClusterInfo.name = name
	caCertInterface, err := jmespath.Search("data.credentials.ca_cert", decoded)
	caCert, ok := caCertInterface.(string)
	if err != nil || !ok {
		return clusterInfo{}, fmt.Errorf("error parsing cluster's ca certificate")
	}
	newClusterInfo.caCert = caCert
	hostInterface, err := jmespath.Search("data.credentials.host", decoded)
	host, ok := hostInterface.(string)
	if err != nil || !ok {
		return clusterInfo{}, fmt.Errorf("error parsing cluster's host")
	}
	newClusterInfo.host = host
	portInterface, err := jmespath.Search("data.credentials.port", decoded)
	port, ok := portInterface.(string)
	if err != nil || !ok {
		return clusterInfo{}, fmt.Errorf("error parsing cluster's port")
	}
	newClusterInfo.port = port
	return newClusterInfo, nil
}

func updateKubeconfig(kubeconfig *api.Config, newClusterInfo clusterInfo) error {
	block, _ := pem.Decode([]byte(newClusterInfo.caCert))
	if block == nil || block.Type != "CERTIFICATE" {
		logger.Fatal("Failed to decode PEM block containing certificate")
	}
	pem.EncodeToMemory(block)
	newCluster := api.Cluster{Server: newClusterInfo.host + ":" + newClusterInfo.port, CertificateAuthorityData: pem.EncodeToMemory(block)}
	newContext := api.Context{AuthInfo: newClusterInfo.name, Cluster: newClusterInfo.name}
	ex, err := os.Executable()
	if err != nil {
		logger.Fatal(err)
	}
	mistCLIPath := path.Dir(ex) + "/" + os.Args[0]
	newAuthinfo := api.AuthInfo{Exec: &api.ExecConfig{Command: mistCLIPath, Args: []string{
		"kubeconfig", "get-cluster-creds", newClusterInfo.name}, APIVersion: "client.authentication.k8s.io/v1alpha1"}}
	if kubeconfig.Clusters == nil {
		kubeconfig.Clusters = make(map[string]*api.Cluster)
	}
	kubeconfig.Clusters[newClusterInfo.name] = &newCluster
	if kubeconfig.Contexts == nil {
		kubeconfig.Contexts = make(map[string]*api.Context)
	}
	kubeconfig.Contexts[newClusterInfo.name] = &newContext
	if kubeconfig.AuthInfos == nil {
		kubeconfig.AuthInfos = make(map[string]*api.AuthInfo)
	}
	kubeconfig.AuthInfos[newClusterInfo.name] = &newAuthinfo
	kubeconfig.CurrentContext = newClusterInfo.name
	return nil
}

func setClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-cluster",
		Short: "Sets a cluster entry in the local kubeconfig or prints the cluster's kubeconfig",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				params := viper.New()
				params.Set("only", "name")
				var decoded interface{}
				_, decoded, _, err := MistApiV2ListClusters(params)
				if err != nil {
					logger.Fatalf("Error calling operation: %s", err.Error())
				}
				data, _ := jmespath.Search("data[].name", decoded)
				j, _ := json.Marshal(data)
				str := strings.Replace(strings.Replace(strings.Replace(string(j[:]), "[", "", -1), "]", "", -1), " ", "\\ ", -1)
				return strings.Split(str, ","), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			params := viper.New()
			params.Set("credentials", true)
			_, decoded, _, err := MistApiV2GetCluster(args[0], params)
			if err != nil {
				logger.Fatalf("Error calling operation: %s", err.Error())
			}
			newClusterInfo, err := parseClusterResponse(decoded)
			if err != nil {
				logger.Fatal("Failed to parse cluster: %s", err.Error())
			}
			modifyKubeconfig := modifyKubeconfigPrompt()
			home := homedir.HomeDir()
			kubeconfig := &api.Config{}
			if modifyKubeconfig && home != "" {
				kubeconfig = clientcmd.GetConfigFromFileOrDie(filepath.Join(home, ".kube", "config"))
			} else if modifyKubeconfig && home == "" {
				fmt.Println("Could not find home directory, printing the kubeconfig instead...")
				modifyKubeconfig = false
			}
			err = updateKubeconfig(kubeconfig, newClusterInfo)
			if err != nil {
				logger.Fatal("Failed to update kubeconfig: %s", err.Error())
			}
			if modifyKubeconfig {
				clientcmd.WriteToFile(*kubeconfig, filepath.Join(home, ".kube", "config"))
				return
			}
			// Convert the kubeconfig struct to json first
			// and then to yaml in order to overcome
			// limitations of the yaml Marshall function.
			jsonBody, err := json.Marshal(kubeconfig)
			if err != nil {
				logger.Fatal(err)
			}
			jsonMap := make(map[string]interface{})
			err = json.Unmarshal(jsonBody, &jsonMap)
			if err != nil {
				logger.Fatal(err)
			}
			kubeconfigYaml, err := yaml.Marshal(jsonMap)
			if err != nil {
				logger.Fatal(err)
			}
			fmt.Printf("%s", string(kubeconfigYaml))
		},
	}
	return cmd
}

func getClusterCredsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-cluster-creds",
		Short: "Gets kubectl compatible cluster creds",
		Args:  cobra.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				params := viper.New()
				params.Set("only", "name")
				var decoded interface{}
				_, decoded, _, err := MistApiV2ListClusters(params)
				if err != nil {
					logger.Fatalf("Error calling operation: %s", err.Error())
				}
				data, _ := jmespath.Search("data[].name", decoded)
				j, _ := json.Marshal(data)
				str := strings.Replace(strings.Replace(strings.Replace(string(j[:]), "[", "", -1), "]", "", -1), " ", "\\ ", -1)
				return strings.Split(str, ","), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			params := viper.New()
			params.Set("credentials", true)
			_, decoded, _, err := MistApiV2GetCluster(args[0], params)
			if err != nil {
				logger.Fatalf("Error calling operation: %s", err.Error())
			}
			tokenInterface, err := jmespath.Search("data.credentials.token", decoded)
			token, ok := tokenInterface.(string)
			if err != nil || !ok {
				logger.Fatalf("Error parsing cluster credentials: %s", err.Error())
			}
			template := `{"kind": "ExecCredential", "apiVersion": "client.authentication.k8s.io/v1alpha1", "spec": {}, "status": {"token": "%s"}}`
			fmt.Printf(template+"\n", token)
		},
	}
	return cmd
}

func kubeconfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Modify kubeconfig",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.AddCommand(setClusterCmd())
	cmd.AddCommand(getClusterCredsCmd())
	cmd.SetErr(os.Stderr)
	return cmd
}
