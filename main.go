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

	"github.com/gorilla/websocket"
	"github.com/jmespath/go-jmespath"
	"github.com/mistio/cobra"
	"github.com/pkg/errors"
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

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Open a shell to a machine",
	Args:  cobra.ExactValidArgs(1),
	Group: "machines",
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

		oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer terminal.Restore(int(os.Stdin.Fd()), oldState)
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

var resourceTypes = []string{"cloud", "machine", "volume", "network", "zone", "image", "key"}

func isValidResourceType(arg string) bool {
	for _, r := range resourceTypes {
		if strings.HasPrefix(r+"s", arg) {
			return true
		}
	}
	return false
}

func getResourceType(arg string) string {
	for _, r := range resourceTypes {
		if strings.HasPrefix(r+"s", arg) {
			return r
		}
	}
	return ""
}

func getResourceTypes(toComplete string) []string {
	return resourceTypes
}

func getResourcesFromBackend(resourceType string, toComplete string) []string {
	params := viper.New()
	params.Set("only", "name")
	var decoded interface{}
	switch resourceType {
	case "cloud":
		_, decoded, _, _ = MistApiV2ListClouds(params)
	case "machine":
		_, decoded, _, _ = MistApiV2ListMachines(params)
	case "volume":
		_, decoded, _, _ = MistApiV2ListVolumes(params)
	case "network":
		_, decoded, _, _ = MistApiV2ListNetworks(params)
	// case "zone":
	// 	_, decoded, _, _ = MistApiV2ListZones(params)
	case "image":
		_, decoded, _, _ = MistApiV2ListImages(params)
	case "key":
		_, decoded, _, _ = MistApiV2ListKeys(params)
	}
	data, _ := jmespath.Search("data[].name", decoded)
	j, _ := json.Marshal(data)
	str := strings.Replace(strings.Replace(strings.Replace(string(j[:]), "[", "", -1), "]", "", -1), " ", "\\ ", -1)
	return strings.Split(str, ",")
}

func getCmd() *cobra.Command {
	params := viper.New()
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return getResourceTypes(toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			if len(args) == 1 {
				resourceType := getResourceType(args[0])
				return getResourcesFromBackend(resourceType, toComplete), cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp

		}, Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a valid resource type")
			}
			if isValidResourceType(args[0]) {

				return nil
			}
			return fmt.Errorf("invalid resource type specified: %s", args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			var decoded map[string]interface{}
			var outputOptions cli.CLIOutputOptions
			var err error

			if strings.HasPrefix("clouds", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetCloud(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListClouds(params)
				}
			} else if strings.HasPrefix("machines", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetMachine(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListMachines(params)
				}
			} else if strings.HasPrefix("volumes", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetVolume(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListVolumes(params)
				}
			} else if strings.HasPrefix("networks", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetNetwork(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListNetworks(params)
				}
			} else if strings.HasPrefix("images", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetImage(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListImages(params)
				}
			} else if strings.HasPrefix("keys", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetKey(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListKeys(params)
				}
			} else if strings.HasPrefix("rules", args[0]) {
				if len(args) == 2 {
					_, decoded, outputOptions, err = MistApiV2GetRule(args[1], params)
				} else if len(args) == 1 {
					_, decoded, outputOptions, err = MistApiV2ListRules(params)
				}
			}

			if err != nil {
				logger.Fatalf("Error calling operation: %s", err.Error())
			}

			if err := cli.Formatter.Format(decoded, outputOptions); err != nil {
				logger.Fatalf("Formatting failed: %s", err.Error())
			}
		},
	}

	cmd.Flags().String("search", "", "Only return results matching search filter")
	cmd.Flags().String("only", "", "Only return these fields")
	cmd.Flags().String("deref", "", "Dereference foreign keys")

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
		Version:   "1.0.0",
	})

	// Initialize the API key authentication.
	apikey.Init("Authorization", apikey.LocationHeader)

	// Add command groups
	cli.Root.AddGroup(&cobra.Group{Group: "clouds", Title: "  # CLOUDS"})
	cli.Root.AddGroup(&cobra.Group{Group: "machines", Title: "  # MACHINES"})
	cli.Root.AddGroup(&cobra.Group{Group: "volumes", Title: "  # VOLUMES"})
	cli.Root.AddGroup(&cobra.Group{Group: "networks", Title: "  # NETWORKS"})
	cli.Root.AddGroup(&cobra.Group{Group: "zones", Title: "  # ZONES"})
	cli.Root.AddGroup(&cobra.Group{Group: "keys", Title: "  # KEYS"})
	cli.Root.AddGroup(&cobra.Group{Group: "images", Title: "  # IMAGES"})
	cli.Root.AddGroup(&cobra.Group{Group: "scripts", Title: "  # SCRIPTS"})
	cli.Root.AddGroup(&cobra.Group{Group: "templates", Title: "  # TEMPLATES"})
	cli.Root.AddGroup(&cobra.Group{Group: "rules", Title: "  # RULES"})
	cli.Root.AddGroup(&cobra.Group{Group: "teams", Title: "  # TEAMS"})

	// Add completion command
	cli.Root.AddCommand(completionCmd)

	// Register auto-generated commands
	mistApiV2Register(false)

	// Add ssh command
	cli.Root.AddCommand(sshCmd)

	// Add get commend
	cli.Root.AddCommand(getCmd())

	cli.Root.Execute()
}
