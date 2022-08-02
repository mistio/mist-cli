package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/containerd/console"
	"github.com/gorilla/websocket"
	"github.com/jmespath/go-jmespath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	trie "github.com/v-pap/trie"
	"gitlab.ops.mist.io/mistio/openapi-cli-generator/apikey"
	"gitlab.ops.mist.io/mistio/openapi-cli-generator/cli"
	terminal "golang.org/x/term"
)

var logger = log.New(os.Stdout, "", 0)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(mist completion bash)

# To load completions for each session, execute once:
Linux:
  $ mist completion bash > /etc/bash_completion.d/mist
MacOS:
  $ mist completion bash > /usr/local/etc/bash_completion.d/mist

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ mist completion zsh > "${fpath[1]}/_mist"

# You will need to start a new shell for this setup to take effect.

Fish:

$ mist completion fish | source

# To load completions for each session, execute once:
$ mist completion fish > ~/.config/fish/completions/mist.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

var customUsageTpl = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var customUsageSubCommandTpl = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}

API Docs:
  $API_DOCS{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func customUsageFunc(c *cobra.Command) error {
	err := setContext()
	if err != nil {
		return err
	}
	server, err := getServer()
	if err != nil {
		return err
	}
	apiDocs := server + c.UsageTemplate()
	c.SetUsageTemplate(strings.ReplaceAll(customUsageSubCommandTpl, "$API_DOCS", apiDocs))
	c.SetUsageFunc(nil)
	c.Usage()
	return nil
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get CLI & API version",
		Args:  cobra.ExactValidArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			type version struct {
				Sha      string `json:"sha"`
				Name     string `json:"name"`
				Repo     string `json:"repo"`
				Modified bool   `json:"modified"`
			}
			type versionResp struct {
				Version version `json:"version"`
			}
			err := setContext()
			if err != nil {
				fmt.Println(err)
				return
			}
			server, err := getServer()
			if err != nil {
				fmt.Println(err)
				return
			}
			url := server + "/version"
			req := cli.Client.Get().URL(url)
			resp, err := req.Do()
			if err != nil {
				fmt.Println(err)
				return
			}
			var ver versionResp
			err = resp.JSON(&ver)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("CLI version: $CLI_VERSION")
			fmt.Printf("Server version: %s - %s#%s", ver.Version.Name, ver.Version.Repo, ver.Version.Sha)
			if ver.Version.Modified {
				fmt.Printf(" modified")
			}
			fmt.Println("")
		},
	}
	return cmd
}

func sshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Open a shell to a machine",
		Args:  cobra.ExactValidArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			if len(args) == 0 {
				params := viper.New()
				params.Set("only", "name")
				params.Set("search", "key_associations:true AND state:running")
				var decoded interface{}
				_, decoded, _, err := MistApiV2ListMachines(params)
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
			machine := args[0]
			// Time allowed to write a message to the peer.
			writeWait := 2 * time.Second

			// Time allowed to read the next pong message from the peer.
			pongWait := 10 * time.Second

			// Send pings to peer with this period. Must be less than pongWait.
			pingPeriod := (pongWait * 9) / 10

			err := setContext()
			if err != nil {
				logger.Fatalf("Could not set context %s", err)
			}
			server, err := getServer()
			if err != nil {
				logger.Fatal(err)
			}
			if !strings.HasSuffix(server, "/") {
				server = server + "/"
			}
			if !strings.HasPrefix(server, "http") {
				server = "http://" + server
			}
			path := server + "api/v2/machines/" + machine + "/actions/ssh"
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}}
			req, err := http.NewRequest("POST", path, nil)
			if err != nil {
				logger.Fatal(err)
			}
			token, err := getToken()
			if err != nil {
				logger.Fatal(err)
			}
			req.Header.Add("Authorization", token)
			if err != nil {
				logger.Fatal(err)
			}
			resp, err := client.Do(req)
			if err != nil {
				logger.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode/100 != 3 {
				logger.Fatalf("Could not SSH into machine: %s", resp.Status)
			}
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Fatal(err)
			}
			location := resp.Header.Get("location")
			c, resp, err := websocket.DefaultDialer.Dial(location, http.Header{"Authorization": []string{token}})
			if err != nil {
				logger.Fatal(err)
			}
			// Handle the case of redirections
			if resp != nil && resp.StatusCode == 302 {
				u, _ := resp.Location()
				c, resp, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{token}})
				if err != nil || resp.StatusCode/100 != 2 {
					logger.Fatal(err)
				}
			}
			defer c.Close()
			current := console.Current()
			if err := current.SetRaw(); err != nil {
				logger.Fatal(err)
			}
			terminal.NewTerminal(current, "")
			defer current.Reset()
			done := make(chan bool)

			var writeMutex sync.Mutex

			err = updateTerminalSize(c, &writeMutex, writeWait)
			if err != nil {
				logger.Fatal(err)
			}

			go handleTerminalResize(c, &done, &writeMutex, writeWait)
			go readFromRemoteStdout(c, &done, pongWait)
			go writeToRemoteStdin(c, &done, &writeMutex, writeWait)
			go sendPingMessages(c, &done, writeWait, pingPeriod)

			<-done
		},
	}
	cmd.SetErr(os.Stderr)
	return cmd
}

func streamingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream JOB_ID",
		Short: "Stream logs of a running script",
		Args:  cobra.ExactValidArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			job_id := args[0]
			// Time allowed to write a message to the peer.
			writeWait := 10 * 3600 * time.Second

			// Time allowed to read the next pong message from the peer.
			pongWait := 20 * time.Second

			// Send pings to peer with this period. Must be less than pongWait.
			pingPeriod := (10 * time.Second * 9) / 10

			server, err := getServer()
			if err != nil {
				logger.Println(err)
				return
			}
			if !strings.HasSuffix(server, "/") {
				server = server + "/"
			}
			if !strings.HasPrefix(server, "http") {
				server = "http://" + server
			}
			path := server + "api/v2/jobs/" + job_id
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}}
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				logger.Println(err)
				return
			}
			token, err := getToken()
			if err != nil {
				logger.Println(err)
				return
			}
			req.Header.Add("Authorization", token)
			if err != nil {
				logger.Println(err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				logger.Println(err)
				return
			}
			var r any
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&r)
			if err != nil {
				logger.Println(err)
				return
			}
			data, data_exists := r.(map[string]any)["data"]
			var location string
			if data_exists {
				_, ok := data.(map[string]any)["stream_uri"].(string)
				if !ok {
					logger.Fatal(errors.New("stream_uri not found in api response for given JOB_ID"))
					return
				}
				location = data.(map[string]any)["stream_uri"].(string)
			} else {
				logger.Fatal(errors.New("api response for given JOB_ID does not contain any data"))
			}
			defer resp.Body.Close()
			c, resp, err := websocket.DefaultDialer.Dial(location, http.Header{"Authorization": []string{token}})
			if err != nil {
				logger.Println(err)
				return
			}
			if resp != nil && resp.StatusCode == 302 {
				u, _ := resp.Location()
				c, resp, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{token}})
			}
			defer c.Close()
			if err != nil {
				logger.Println(err)
				return
			}

			current := console.Current()
			if err := current.SetRaw(); err != nil {
				panic(err)
			}
			terminal.NewTerminal(current, "")
			defer current.Reset()
			done := make(chan bool)

			var writeMutex sync.Mutex

			err = updateTerminalSize(c, &writeMutex, writeWait)
			if err != nil {
				logger.Println(err)
				return
			}
			go func() {
				cmd.InOrStdin()
				_, _, err := bufio.NewReader(cmd.InOrStdin()).ReadRune()
				if err != nil {
					logger.Println(err)
					return
				}
				done <- true
				os.Exit(0)
			}()
			go readFromRemoteStdout(c, &done, pongWait)
			go sendPingMessages(c, &done, writeWait, pingPeriod)
			<-done
		},
	}
	cmd.SetErr(os.Stderr)
	return cmd
}

func getResourceMeterCmdRun(params *viper.Viper, resource string, detailedName bool) {
	dtStart := params.GetString("start")
	if dtStart == "" {
		dtStart = fmt.Sprintf("%d", (time.Now()).Unix()-60*60)
	}
	dtEnd := params.GetString("end")
	if dtEnd == "" {
		dtEnd = fmt.Sprintf("%d", (time.Now()).Unix())
	}
	_, resourceMetricsStart, _ := getMeteringData(dtStart, dtEnd, params.GetString("search"), fmt.Sprintf("first_over_time({metering=\"true\",%s_id=~\".+\"}", resource)+"[%ds])")
	metricsSet, resourceMetricsEnd, resourceNames := getMeteringData(dtStart, dtEnd, params.GetString("search"), fmt.Sprintf("last_over_time({metering=\"true\",%s_id=~\".+\"}", resource)+"[%ds])")
	if detailedName {
		for resourceID, name := range resourceNames {
			resourceNames[resourceID] = resource + "/" + name
		}
	}
	machineMetricsGauges := calculateDiffs(resourceMetricsStart, resourceMetricsEnd, metricsSet)
	formatMeteringData(metricsSet, machineMetricsGauges, resourceNames)
}

func getResourceMeterCmd(resource string, aliasesMap map[string][]string, detailedName bool) *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:     resource,
		Aliases: aliasesMap[resource],
		Short:   fmt.Sprintf("Get metering data for %s resource", resource),
		Args:    cobra.ExactValidArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			getResourceMeterCmdRun(params, resource, false)
		},
	}
	cmd.Flags().String("start", "", "start <rfc3339 | unix_timestamp>")
	cmd.Flags().String("end", "", "end <rfc3339 | unix_timestamp>")
	cmd.Flags().String("search", "", "Only return results matching search filter")

	cli.SetCustomFlags(cmd)

	if cmd.Flags().HasFlags() {
		params.BindPFlags(cmd.Flags())
	}
	return cmd
}

func getAllResourcesMeterCmd(resources []string, aliasesMap map[string][]string) *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:     "all",
		Aliases: aliasesMap["all"],
		Short:   "Get metering data for all resources",
		Args:    cobra.ExactValidArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			for i, resource := range resources {
				getResourceMeterCmdRun(params, resource, true)
				if i != len(resources)-1 {
					fmt.Println("")
				}
			}
		},
	}
	cmd.Flags().String("start", "", "start <rfc3339 | unix_timestamp>")
	cmd.Flags().String("end", "", "end <rfc3339 | unix_timestamp>")
	cmd.Flags().String("search", "", "Only return results matching search filter")

	cli.SetCustomFlags(cmd)

	if cmd.Flags().HasFlags() {
		params.BindPFlags(cmd.Flags())
	}
	return cmd
}

func calculateAliases(command, suffix string) []string {
	if len(command) == 0 {
		return []string{}
	}
	prefix := strings.TrimSuffix(command, suffix)
	subString := ""
	aliases := []string{}
	for i := len(prefix); i < len(command)-1; i++ {
		subString += string(command[i])
		aliases = append(aliases, prefix+subString)
	}
	if len(aliases) > 0 && string(command[len(command)-1]) != "s" {
		aliases = append(aliases, command+"s")
	}
	return aliases
}

func meterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meter",
		Short: "Get metering data",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	resources := []string{"machine", "volume"}
	resourcesTrie := trie.New()
	aliasesMap := make(map[string][]string)
	for _, resource := range append(resources, []string{"all"}...) {
		resourcesTrie.Insert(resource)
	}
	for _, resource := range append(resources, []string{"all"}...) {
		suffix, ok := resourcesTrie.FindLongestUniqueSuffix(resource)
		if !ok {
			continue
		}
		aliasesMap[resource] = calculateAliases(resource, suffix)
	}
	for _, resource := range resources {
		cmd.AddCommand(getResourceMeterCmd(resource, aliasesMap, false))
	}
	cmd.AddCommand(getAllResourcesMeterCmd(resources, aliasesMap))
	cmd.SetErr(os.Stderr)
	return cmd
}

func main() {
	cli.Init(&cli.Config{
		AppName:   "mist",
		EnvPrefix: "MIST",
		Version:   "",
	})
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize the API key authentication.
	apikey.Init("Authorization", apikey.LocationHeader)

	// Add command groups
	/*cli.Root.AddGroup(&cobra.Group{Group: "clouds", Title: "  # CLOUDS"})
	cli.Root.AddGroup(&cobra.Group{Group: "machines", Title: "  # MACHINES"})
	cli.Root.AddGroup(&cobra.Group{Group: "volumes", Title: "  # VOLUMES"})
	cli.Root.AddGroup(&cobra.Group{Group: "networks", Title: "  # NETWORKS"})
	cli.Root.AddGroup(&cobra.Group{Group: "zones", Title: "  # ZONES"})
	cli.Root.AddGroup(&cobra.Group{Group: "keys", Title: "  # KEYS"})
	cli.Root.AddGroup(&cobra.Group{Group: "images", Title: "  # IMAGES"})
	cli.Root.AddGroup(&cobra.Group{Group: "scripts", Title: "  # SCRIPTS"})
	cli.Root.AddGroup(&cobra.Group{Group: "templates", Title: "  # TEMPLATES"})
	cli.Root.AddGroup(&cobra.Group{Group: "rules", Title: "  # RULES"})
	cli.Root.AddGroup(&cobra.Group{Group: "teams", Title: "  # TEAMS"})*/

	// Add completion command
	cli.Root.AddCommand(completionCmd)

	// Register auto-generated commands
	mistApiV2Register(false)

	// Add version command
	cli.Root.AddCommand(versionCmd())

	// Add ssh command
	cli.Root.AddCommand(sshCmd())

	cli.Root.AddCommand(streamingCmd())
	// Add metering command
	cli.Root.AddCommand(meterCmd())

	cli.Root.AddCommand(tagCmd())

	cli.Root.AddCommand(untagCmd())

	cli.Root.AddCommand(kubeconfigCmd())

	cli.Root.Execute()
}
