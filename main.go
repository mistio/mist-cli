package main

import (
	"encoding/json"
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
			err := setMistContext()
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

			err := setMistContext()
			if err != nil {
				fmt.Printf("Cannot set context %s\n", err)
				return
			}
			server, err := getServer()
			if err != nil {
				fmt.Println(err)
				return
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
			token, err := getToken()
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Authorization", token)
			if err != nil {
				fmt.Println(err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			location := resp.Header.Get("location")
			c, resp, err := websocket.DefaultDialer.Dial(location, http.Header{"Authorization": []string{token}})
			if resp != nil && resp.StatusCode == 302 {
				u, _ := resp.Location()
				c, resp, err = websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{token}})
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			defer c.Close()
			if err != nil {
				fmt.Println(err)
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
				fmt.Println(err)
				return
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

func meterCmd() *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:   "meter",
		Short: "Get metering data",
		Args:  cobra.ExactValidArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			dtStart := params.GetString("start")
			if dtStart == "" {
				dtStart = fmt.Sprintf("%d", (time.Now()).Unix()-60*60)
			}
			dtEnd := params.GetString("end")
			if dtEnd == "" {
				dtEnd = fmt.Sprintf("%d", (time.Now()).Unix())
			}
			_, machineMetricsStart, _ := getMeteringData(dtStart, dtEnd, params.GetString("search"), "first_over_time({metering=\"true\"}[%ds])")
			metricsSet, machineMetricsEnd, machineNames := getMeteringData(dtStart, dtEnd, params.GetString("search"), "last_over_time({metering=\"true\"}[%ds])")
			machineMetricsGauges := calculateDiffs(machineMetricsStart, machineMetricsEnd, metricsSet)
			formatMeteringData(metricsSet, machineMetricsGauges, machineNames)
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

func main() {
	cli.Init(&cli.Config{
		AppName:   "mist",
		EnvPrefix: "MIST",
		Version:   "",
	})

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

	// Add metering command
	cli.Root.AddCommand(meterCmd())

	cli.Root.Execute()
}
