// Code generated by openapi-cli-generator. DO NOT EDIT.
// See https://github.com/danielgtaylor/openapi-cli-generator

package main

import (
	"fmt"
	"strings"

	"github.com/danielgtaylor/openapi-cli-generator/cli"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/h2non/gentleman.v2"
)

var mistApiV2Subcommand bool

func mistApiV2Servers() []map[string]string {
	return []map[string]string{

		map[string]string{
			"description": "dogfood",
			"url":         "https://dogfood.ops.mist.io/",
		},
	}
}

// MistApiV2AddCloud Add cloud
func MistApiV2AddCloud(params *viper.Viper, body string) (*gentleman.Response, map[string]interface{}, error) {
	handlerPath := "add-cloud"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/clouds"

	req := cli.Client.Post().URL(url)

	if body != "" {
		req = req.AddHeader("Content-Type", "application/json").BodyString(body)
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded map[string]interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after.(map[string]interface{})
	}

	return resp, decoded, nil
}

// MistApiV2ListClouds List clouds
func MistApiV2ListClouds(params *viper.Viper) (*gentleman.Response, map[string]interface{}, error) {
	handlerPath := "list-clouds"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/clouds"

	req := cli.Client.Get().URL(url)

	paramSearch := params.GetString("search")
	if paramSearch != "" {
		req = req.AddQuery("search", fmt.Sprintf("%v", paramSearch))
	}
	paramSort := params.GetString("sort")
	if paramSort != "" {
		req = req.AddQuery("sort", fmt.Sprintf("%v", paramSort))
	}
	paramStart := params.GetString("start")
	if paramStart != "" {
		req = req.AddQuery("start", fmt.Sprintf("%v", paramStart))
	}
	paramLimit := params.GetInt64("limit")
	if paramLimit != 0 {
		req = req.AddQuery("limit", fmt.Sprintf("%v", paramLimit))
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded map[string]interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after.(map[string]interface{})
	}

	return resp, decoded, nil
}

// MistApiV2DeleteCloud Delete cloud
func MistApiV2DeleteCloud(paramCloud string, params *viper.Viper) (*gentleman.Response, interface{}, error) {
	handlerPath := "delete-cloud"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/clouds/{cloud}"
	url = strings.Replace(url, "{cloud}", paramCloud, 1)

	req := cli.Client.Delete().URL(url)

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after
	}

	return resp, decoded, nil
}

// MistApiV2GetCloud Get cloud
func MistApiV2GetCloud(paramCloud string, params *viper.Viper) (*gentleman.Response, interface{}, error) {
	handlerPath := "get-cloud"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/clouds/{cloud}"
	url = strings.Replace(url, "{cloud}", paramCloud, 1)

	req := cli.Client.Get().URL(url)

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after
	}

	return resp, decoded, nil
}

// MistApiV2ListRules Get rules
func MistApiV2ListRules(params *viper.Viper) (*gentleman.Response, map[string]interface{}, error) {
	handlerPath := "list-rules"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/rules"

	req := cli.Client.Get().URL(url)

	paramFilter := params.GetString("filter")
	if paramFilter != "" {
		req = req.AddQuery("filter", fmt.Sprintf("%v", paramFilter))
	}
	paramSort := params.GetString("sort")
	if paramSort != "" {
		req = req.AddQuery("sort", fmt.Sprintf("%v", paramSort))
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded map[string]interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after.(map[string]interface{})
	}

	return resp, decoded, nil
}

// MistApiV2AddRule Add rule
func MistApiV2AddRule(paramQueries string, paramWindow string, paramFrequency string, paramTriggerAfter string, paramActions string, paramSelectors string, params *viper.Viper, body string) (*gentleman.Response, map[string]interface{}, error) {
	handlerPath := "add-rule"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/rules"

	req := cli.Client.Post().URL(url)

	req = req.AddQuery("queries", paramQueries)

	req = req.AddQuery("window", paramWindow)

	req = req.AddQuery("frequency", paramFrequency)

	req = req.AddQuery("trigger_after", paramTriggerAfter)

	req = req.AddQuery("actions", paramActions)

	req = req.AddQuery("selectors", paramSelectors)

	if body != "" {
		req = req.AddHeader("Content-Type", "").BodyString(body)
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded map[string]interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after.(map[string]interface{})
	}

	return resp, decoded, nil
}

// MistApiV2ToggleRule Toggle rule
func MistApiV2ToggleRule(paramRule string, paramAction string, params *viper.Viper, body string) (*gentleman.Response, interface{}, error) {
	handlerPath := "toggle-rule"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/rules/{rule}"
	url = strings.Replace(url, "{rule}", paramRule, 1)

	req := cli.Client.Put().URL(url)

	req = req.AddQuery("action", paramAction)

	if body != "" {
		req = req.AddHeader("Content-Type", "").BodyString(body)
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after
	}

	return resp, decoded, nil
}

// MistApiV2DeleteRule Delete rule
func MistApiV2DeleteRule(paramRule string, params *viper.Viper) (*gentleman.Response, interface{}, error) {
	handlerPath := "delete-rule"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/rules/{rule}"
	url = strings.Replace(url, "{rule}", paramRule, 1)

	req := cli.Client.Delete().URL(url)

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after
	}

	return resp, decoded, nil
}

// MistApiV2RenameRule Rename rule
func MistApiV2RenameRule(paramRule string, paramAction string, params *viper.Viper, body string) (*gentleman.Response, interface{}, error) {
	handlerPath := "rename-rule"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/rules/{rule}"
	url = strings.Replace(url, "{rule}", paramRule, 1)

	req := cli.Client.Patch().URL(url)

	req = req.AddQuery("action", paramAction)

	if body != "" {
		req = req.AddHeader("Content-Type", "").BodyString(body)
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after
	}

	return resp, decoded, nil
}

// MistApiV2UpdateRule Update rule
func MistApiV2UpdateRule(paramRule string, params *viper.Viper, body string) (*gentleman.Response, map[string]interface{}, error) {
	handlerPath := "update-rule"
	if mistApiV2Subcommand {
		handlerPath = "mist-api-v2 " + handlerPath
	}

	server := viper.GetString("server")
	if server == "" {
		server = mistApiV2Servers()[viper.GetInt("server-index")]["url"]
	}

	url := server + "/api/v2/rules/{rule}"
	url = strings.Replace(url, "{rule}", paramRule, 1)

	req := cli.Client.Post().URL(url)

	paramQueries := params.GetString("queries")
	if paramQueries != "" {
		req = req.AddQuery("queries", fmt.Sprintf("%v", paramQueries))
	}
	paramWindow := params.GetString("window")
	if paramWindow != "" {
		req = req.AddQuery("window", fmt.Sprintf("%v", paramWindow))
	}
	paramFrequency := params.GetString("frequency")
	if paramFrequency != "" {
		req = req.AddQuery("frequency", fmt.Sprintf("%v", paramFrequency))
	}
	paramTriggerAfter := params.GetString("trigger-after")
	if paramTriggerAfter != "" {
		req = req.AddQuery("trigger_after", fmt.Sprintf("%v", paramTriggerAfter))
	}
	paramActions := params.GetString("actions")
	if paramActions != "" {
		req = req.AddQuery("actions", fmt.Sprintf("%v", paramActions))
	}
	paramSelectors := params.GetString("selectors")
	if paramSelectors != "" {
		req = req.AddQuery("selectors", fmt.Sprintf("%v", paramSelectors))
	}

	if body != "" {
		req = req.AddHeader("Content-Type", "").BodyString(body)
	}

	cli.HandleBefore(handlerPath, params, req)

	resp, err := req.Do()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Request failed")
	}

	var decoded map[string]interface{}

	if resp.StatusCode < 400 {
		if err := cli.UnmarshalResponse(resp, &decoded); err != nil {
			return nil, nil, errors.Wrap(err, "Unmarshalling response failed")
		}
	} else {
		return nil, nil, errors.Errorf("HTTP %d: %s", resp.StatusCode, resp.String())
	}

	after := cli.HandleAfter(handlerPath, params, resp, decoded)
	if after != nil {
		decoded = after.(map[string]interface{})
	}

	return resp, decoded, nil
}

func mistApiV2Register(subcommand bool) {
	root := cli.Root

	if subcommand {
		root = &cobra.Command{
			Use:   "mist-api-v2",
			Short: "Mist API",
			Long:  cli.Markdown(""),
		}
		mistApiV2Subcommand = true
	} else {
		cli.Root.Short = "Mist API"
		cli.Root.Long = cli.Markdown("")
	}

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "add-cloud",
			Short:   "Add cloud",
			Long:    cli.Markdown("Adds a new cloud and returns the cloud's id. ADD permission required on cloud.\n## Request Schema (application/json)\n\nproperties:\n  credentials:\n    $ref: '#/components/schemas/CloudCredentials'\n  features:\n    $ref: '#/components/schemas/CloudFeatures'\n  provider:\n    description: The provider of the cloud\n    enum:\n    - amazon\n    - digitalocean\n    - google\n    - openstack\n    - packet\n    - vsphere\n    type: string\n  title:\n    description: The name of the cloud to add\n    type: string\nrequired:\n- title\n- provider\n- credentials\ntype: object\n"),
			Example: examples,
			Group:   "clouds",
			Args:    cobra.MinimumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				body, err := cli.GetBody("application/json", args[0:])
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to get body")
				}

				_, decoded, err := MistApiV2AddCloud(params, body)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "list-clouds",
			Short:   "List clouds",
			Long:    cli.Markdown("List clouds owned by the active org. READ permission required on cloud."),
			Example: examples,
			Group:   "clouds",
			Args:    cobra.MinimumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {

				_, decoded, err := MistApiV2ListClouds(params)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cmd.Flags().String("search", "", "Only return results matching search filter")
		cmd.Flags().String("sort", "", "Order results by")
		cmd.Flags().String("start", "", "Start results from index or id")
		cmd.Flags().Int64("limit", 0, "Limit number of results, 1000 max")

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "delete-cloud cloud",
			Short:   "Delete cloud",
			Long:    cli.Markdown("Delete target cloud"),
			Example: examples,
			Group:   "clouds",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {

				_, decoded, err := MistApiV2DeleteCloud(args[0], params)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "get-cloud cloud",
			Short:   "Get cloud",
			Long:    cli.Markdown("Get details about target cloud"),
			Example: examples,
			Group:   "clouds",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {

				_, decoded, err := MistApiV2GetCloud(args[0], params)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "list-rules",
			Short:   "Get rules",
			Long:    cli.Markdown("Return a filtered list of rules"),
			Example: examples,
			Group:   "rules",
			Args:    cobra.MinimumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {

				_, decoded, err := MistApiV2ListRules(params)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cmd.Flags().String("filter", "", "Only return results matching filter")
		cmd.Flags().String("sort", "", "Order results by")

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "add-rule queries window frequency trigger-after actions selectors",
			Short:   "Add rule",
			Long:    cli.Markdown("Add a new rule, READ permission required on target resource, ADD permission required on Rule"),
			Example: examples,
			Group:   "rules",
			Args:    cobra.MinimumNArgs(6),
			Run: func(cmd *cobra.Command, args []string) {
				body, err := cli.GetBody("", args[6:])
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to get body")
				}

				_, decoded, err := MistApiV2AddRule(args[0], args[1], args[2], args[3], args[4], args[5], params, body)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "toggle-rule rule action",
			Short:   "Toggle rule",
			Long:    cli.Markdown("Enable or disable a rule"),
			Example: examples,
			Group:   "rules",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				body, err := cli.GetBody("", args[2:])
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to get body")
				}

				_, decoded, err := MistApiV2ToggleRule(args[0], args[1], params, body)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "delete-rule rule",
			Short:   "Delete rule",
			Long:    cli.Markdown("Delete a rule given its UUID."),
			Example: examples,
			Group:   "rules",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {

				_, decoded, err := MistApiV2DeleteRule(args[0], params)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "rename-rule rule action",
			Short:   "Rename rule",
			Long:    cli.Markdown("Rename a rule"),
			Example: examples,
			Group:   "rules",
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				body, err := cli.GetBody("", args[2:])
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to get body")
				}

				_, decoded, err := MistApiV2RenameRule(args[0], args[1], params, body)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

	func() {
		params := viper.New()

		var examples string

		cmd := &cobra.Command{
			Use:     "update-rule rule",
			Short:   "Update rule",
			Long:    cli.Markdown("Update a rule given its UUID, EDIT permission required on rule"),
			Example: examples,
			Group:   "rules",
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				body, err := cli.GetBody("", args[1:])
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to get body")
				}

				_, decoded, err := MistApiV2UpdateRule(args[0], params, body)
				if err != nil {
					log.Fatal().Err(err).Msg("Error calling operation")
				}

				if err := cli.Formatter.Format(decoded); err != nil {
					log.Fatal().Err(err).Msg("Formatting failed")
				}

			},
		}
		root.AddCommand(cmd)

		cmd.Flags().String("queries", "", "")
		cmd.Flags().String("window", "", "")
		cmd.Flags().String("frequency", "", "")
		cmd.Flags().String("trigger-after", "", "")
		cmd.Flags().String("actions", "", "")
		cmd.Flags().String("selectors", "", "")

		cli.SetCustomFlags(cmd)

		if cmd.Flags().HasFlags() {
			params.BindPFlags(cmd.Flags())
		}

	}()

}
