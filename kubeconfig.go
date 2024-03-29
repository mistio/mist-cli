package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmespath/go-jmespath"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.ops.mist.io/mistio/openapi-cli-generator/cli"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func modifyKubeconfigPrompt() bool {
	prompt := promptui.Prompt{
		Label:     "Modify local kubeconfig",
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err == nil
}

type clusterInfo struct {
	name   string
	host   string
	port   string
	caCert string
}

type clusterCreds struct {
	token  string
	expiry string
}

type cluster string

func (c cluster) getCredsFromCache() *clusterCreds {
	tokenKey := "contexts." + viper.GetString("context") + "." + string(c) + ".token"
	expiresKey := "contexts." + viper.GetString("context") + "." + string(c) + ".expires"
	cachedExpiry := cli.ClusterCache.GetTime(expiresKey)
	if cachedExpiry.IsZero() || cachedExpiry.Before(time.Now().Add(time.Minute).UTC()) {
		return nil
	}
	return &clusterCreds{token: cli.ClusterCache.GetString(tokenKey), expiry: cachedExpiry.Format(time.RFC3339)}
}

func (c cluster) setCredsToCache(creds *clusterCreds) error {
	tokenKey := "contexts." + viper.GetString("context") + "." + string(c) + ".token"
	expiresKey := "contexts." + viper.GetString("context") + "." + string(c) + ".expires"
	cli.ClusterCache.Set(tokenKey, creds.token)
	cli.ClusterCache.Set(expiresKey, creds.expiry)
	err := cli.ClusterCache.WriteConfig()
	return err
}

func (c cluster) getCredsFromMist() *clusterCreds {
	params := viper.New()
	params.Set("credentials", true)
	_, decoded, _, err := MistApiV2GetCluster(string(c), params)
	if err != nil {
		logger.Fatalf("Error calling operation: %s", err.Error())
	}
	tokenInterface, err := jmespath.Search("data.credentials.token", decoded)
	token, ok := tokenInterface.(string)
	if err != nil || !ok {
		logger.Fatalf("Error parsing cluster credentials: %s", err.Error())
	}
	tokenExpiryInterface, err := jmespath.Search("data.credentials.token_expiry", decoded)
	tokenExpiry, ok := tokenExpiryInterface.(string)
	if err != nil || !ok {
		logger.Fatalf("Error parsing cluster credentials: %s", err.Error())
	}
	creds := &clusterCreds{token: token, expiry: tokenExpiry}
	err = c.setCredsToCache(creds)
	if err != nil {
		logger.Printf("Could not save cluster credentials to cache: %s", err.Error())
	}
	return creds
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

func prepareAdress(host string, port string) string {
	httpPrefix := "https://"
	if port == "80" {
		httpPrefix = "http://"
	}
	address := host + ":" + port
	if !strings.HasPrefix(address, httpPrefix) {
		address = httpPrefix + address
	}
	return address
}

func updateKubeconfig(kubeconfig *api.Config, newClusterInfo clusterInfo) error {
	block, _ := pem.Decode([]byte(newClusterInfo.caCert))
	if block == nil || block.Type != "CERTIFICATE" {
		logger.Fatal("Failed to decode PEM block containing certificate")
	}
	pem.EncodeToMemory(block)
	newCluster := api.Cluster{Server: prepareAdress(newClusterInfo.host, newClusterInfo.port), CertificateAuthorityData: pem.EncodeToMemory(block)}
	newContext := api.Context{AuthInfo: newClusterInfo.name, Cluster: newClusterInfo.name}
	ex, err := os.Executable()
	if err != nil {
		logger.Fatal(err)
	}
	mistCLIPath := path.Dir(ex) + "/" + os.Args[0]
	newAuthinfo := api.AuthInfo{Exec: &api.ExecConfig{Command: mistCLIPath, InteractiveMode: "Never", ProvideClusterInfo: true, Args: []string{
		"kubeconfig", "get-cluster-creds", newClusterInfo.name, "--context=" + viper.GetString("context")}, APIVersion: "client.authentication.k8s.io/v1"}}
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

func kubeconfigAutocomplete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func kubeconfigGetCmd() *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:               "update",
		Short:             "Add or update context for cluster(s) in the local kubeconfig",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: kubeconfigAutocomplete,
		Run: func(cmd *cobra.Command, args []string) {
			modifyKubeconfig := params.GetBool("yes")
			if !modifyKubeconfig && !modifyKubeconfigPrompt() {
				fmt.Println("Aborting...")
				return
			}
			home := homedir.HomeDir()
			if home == "" {
				logger.Fatal("Could not find home directory")
			}
			kubeconfig := clientcmd.GetConfigFromFileOrDie(filepath.Join(home, ".kube", "config"))
			addedClusters := ""
			for _, cluster := range args {
				paramsGetCluster := viper.New()
				paramsGetCluster.Set("credentials", true)
				_, decoded, _, err := MistApiV2GetCluster(cluster, paramsGetCluster)
				if err != nil {
					logger.Fatalf("Error calling operation: %s", err.Error())
				}
				newClusterInfo, err := parseClusterResponse(decoded)
				if err != nil {
					logger.Fatalf("Failed to parse cluster: %s", err.Error())
				}
				err = updateKubeconfig(kubeconfig, newClusterInfo)
				if err != nil {
					logger.Fatalf("Failed to update kubeconfig: %s", err.Error())
				}
				addedClusters = addedClusters + "\"" + newClusterInfo.name + "\","
			}
			clientcmd.WriteToFile(*kubeconfig, filepath.Join(home, ".kube", "config"))
			fmt.Printf("Clusters %s added to the local kubeconfig\n", strings.TrimSuffix(addedClusters, ","))
		},
	}
	cmd.Flags().Bool("yes", false, "Override yes/no prompt")
	params.BindPFlags(cmd.Flags())
	return cmd
}

func kubeconfigShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "show",
		Short:             "Display kubeconfig for cluster(s)",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: kubeconfigAutocomplete,
		Run: func(cmd *cobra.Command, args []string) {
			kubeconfig := &api.Config{}
			for _, cluster := range args {
				params := viper.New()
				params.Set("credentials", true)
				_, decoded, _, err := MistApiV2GetCluster(cluster, params)
				if err != nil {
					logger.Fatalf("Error calling operation: %s", err.Error())
				}
				newClusterInfo, err := parseClusterResponse(decoded)
				if err != nil {
					logger.Fatalf("Failed to parse cluster: %s", err.Error())
				}
				err = updateKubeconfig(kubeconfig, newClusterInfo)
				if err != nil {
					logger.Fatalf("Failed to update kubeconfig: %s", err.Error())
				}
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

func kubeconfigCreds() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "get-cluster-creds",
		Short:  "Get kubectl compatible cluster creds",
		Args:   cobra.MinimumNArgs(1),
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			err := setContext()
			if err != nil {
				fmt.Printf("Could not set context %s\n", err)
				return
			}
			c := cluster(args[0])
			creds := c.getCredsFromCache()
			if creds == nil {
				creds = c.getCredsFromMist()
			}
			template := `{"kind": "ExecCredential", "apiVersion": "client.authentication.k8s.io/v1", "spec": {}, "status": {"expirationTimestamp": "%s", "token": "%s"}}`
			fmt.Printf(template+"\n", creds.expiry, creds.token)
		},
	}
	return cmd
}

func kubeconfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Configure kubectl access to cluster",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.AddCommand(kubeconfigGetCmd())
	cmd.AddCommand(kubeconfigShowCmd())
	cmd.AddCommand(kubeconfigCreds())
	cmd.SetErr(os.Stderr)
	return cmd
}
