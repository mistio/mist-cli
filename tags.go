package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/jmespath/go-jmespath"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.ops.mist.io/mistio/openapi-cli-generator/cli"
	"gopkg.in/h2non/gentleman.v2"
)

var taggableResources []string = []string{
	"cloud",
	"cluster",
	"key",
	"machine",
	"network",
	"rule",
	"schedule",
	"script",
	"secret",
	"volume",
	"zone"}

var resourceListControllersMap map[string]func(params *viper.Viper) (*gentleman.Response, map[string]interface{}, cli.CLIOutputOptions, error) = map[string]func(params *viper.Viper) (*gentleman.Response, map[string]interface{}, cli.CLIOutputOptions, error){
	"cloud":    MistApiV2ListClouds,
	"cluster":  MistApiV2ListClusters,
	"key":      MistApiV2ListKeys,
	"machine":  MistApiV2ListMachines,
	"network":  MistApiV2ListNetworks,
	"rule":     MistApiV2ListRules,
	"schedule": MistApiV2ListSchedules,
	"script":   MistApiV2ListScripts,
	"secret":   MistApiV2ListSecrets,
	"volume":   MistApiV2ListVolumes,
	"zone":     MistApiV2ListZones,
}

var resourceGetControllersMap map[string]func(param string, params *viper.Viper) (*gentleman.Response, map[string]interface{}, cli.CLIOutputOptions, error) = map[string]func(param string, params *viper.Viper) (*gentleman.Response, map[string]interface{}, cli.CLIOutputOptions, error){
	"cloud":    MistApiV2GetCloud,
	"cluster":  MistApiV2GetCluster,
	"key":      MistApiV2GetKey,
	"machine":  MistApiV2GetMachine,
	"network":  MistApiV2GetNetwork,
	"rule":     MistApiV2GetRule,
	"schedule": MistApiV2GetSchedule,
	"script":   MistApiV2GetScript,
	"secret":   MistApiV2GetSecret,
	"volume":   MistApiV2GetVolume,
	"zone":     MistApiV2GetZone,
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Resource struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

type Operation struct {
	Operation string         `json:"operation"`
	Tags      []KeyValuePair `json:"tags"`
	Resources []Resource     `json:"resources"`
}

type tagResourceBody struct {
	Operations []Operation `json:"operations"`
}

func tagValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return taggableResources, cobra.ShellCompDirectiveNoFileComp
	}
	if len(args) == 1 {
		params := viper.New()
		params.Set("only", "name")
		var decoded interface{}
		_, decoded, _, err := resourceListControllersMap[args[0]](params)
		if err != nil {
			logger.Fatalf("Error calling operation: %s", err.Error())
		}
		data, _ := jmespath.Search("data[].name", decoded)
		j, _ := json.Marshal(data)
		str := strings.Replace(strings.Replace(strings.Replace(string(j[:]), "[", "", -1), "]", "", -1), " ", "\\ ", -1)
		return strings.Split(str, ","), cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func tagRun(cmd *cobra.Command, args []string, params *viper.Viper, tagOperation string) {
	resourceType := args[0]
	resourceName := args[1]
	stringTags := args[2]
	_, decodedResource, _, err := resourceGetControllersMap[resourceType](resourceName, params)
	rawResourceID, _ := jmespath.Search("data.id", decodedResource)
	resourceID, ok := rawResourceID.(string)
	if !ok {
		logger.Fatalf("Error parsing resource: %s", err.Error())
	}
	tags := []KeyValuePair{}
	for _, stringTag := range strings.Split(stringTags, ",") {
		splittedTag := strings.Split(stringTag, "=")
		kv := KeyValuePair{}
		kv.Key = splittedTag[0]
		if len(splittedTag) > 1 {
			kv.Value = splittedTag[1]
		}
		tags = append(tags, kv)
	}
	resources := []Resource{{ResourceType: resourceType + "s", ResourceID: resourceID}}
	operations := []Operation{{Operation: tagOperation, Tags: tags, Resources: resources}}
	body := tagResourceBody{Operations: operations}
	rawBody, err := json.Marshal(body)
	if err != nil {
		logger.Fatalf("Error marshalling tags: %s", err.Error())
	}
	_, decodedTag, outputOptions, err := MistApiV2TagResources(params, string(rawBody))
	if err != nil {
		logger.Fatalf("Error calling operation: %s", err.Error())
	}

	if err := cli.Formatter.Format(decodedTag, params, outputOptions); err != nil {
		logger.Fatalf("Formatting failed: %s", err.Error())
	}

}

func tagCmd() *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:               "tag RESOURCE_TYPE RESOURCE TAGS",
		Short:             "Tag resource",
		Args:              cobra.MinimumNArgs(3),
		ValidArgsFunction: tagValidArgsFunction,
		Run: func(cmd *cobra.Command, args []string) {
			tagRun(cmd, args, params, "add")
		},
	}
	cmd.SetErr(os.Stderr)
	return cmd
}

func untagCmd() *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:               "untag RESOURCE_TYPE RESOURCE TAGS",
		Short:             "Untag resource",
		Args:              cobra.MinimumNArgs(3),
		ValidArgsFunction: tagValidArgsFunction,
		Run: func(cmd *cobra.Command, args []string) {
			tagRun(cmd, args, params, "remove")
		},
	}
	cmd.SetErr(os.Stderr)
	return cmd
}

// MistApiV2TagResources Tag Resources
func MistApiV2TagResources(params *viper.Viper, body string) (*gentleman.Response, interface{}, cli.CLIOutputOptions, error) {
	handlerPath := "tag-resources"
	if mistApiV2Subcommand {
		handlerPath = "Mist CLI " + handlerPath
	}

	err := setContext()
	if err != nil {
		return nil, nil, cli.CLIOutputOptions{}, err
	}

	server, err := getServer()
	if err != nil {
		return nil, nil, cli.CLIOutputOptions{}, err
	}

	url := server + "/api/v2/tags"

	req := cli.Client.Post().URL(url)

	if body != "" {
		req = req.AddHeader("Content-Type", "application/json").BodyString(body)
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return resp, nil, cli.CLIOutputOptions{}, errors.Wrap(err, "Request failed")
	}

	var decoded interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return resp, nil, cli.CLIOutputOptions{}, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return resp, nil, cli.CLIOutputOptions{}, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after
	}

	return resp, decoded, cli.CLIOutputOptions{[]string{}, []string{}, []string{}, []string{}, map[string]string{}}, nil
}
