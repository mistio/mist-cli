package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"io"
	"time"

	"github.com/danielgtaylor/openapi-cli-generator/apikey"
	"github.com/danielgtaylor/openapi-cli-generator/cli"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	terminal "golang.org/x/term"
)

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

func main() {
	cli.Init(&cli.Config{
		AppName:   "mist",
		EnvPrefix: "MIST",
		Version:   "1.0.0",
	})

	// Initialize the API key authentication.
	apikey.Init("Authorization", apikey.LocationHeader)

	// Add completion command
	cli.Root.AddCommand(completionCmd)

	cli.Root.AddCommand(&cobra.Command{
		Use:   "ssh",
		Short: "Open a shell to a machine",
		Args:  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			machine := args[0]
			// Time allowed to write a message to the peer.
			writeWait := 2 * time.Second

			// Time allowed to read the next pong message from the peer.
			pongWait := 10 * time.Second

			// Send pings to peer with this period. Must be less than pongWait.
			pingPeriod := (pongWait * 9) / 10

			err := setProfile()
			if err != nil {
				fmt.Println("Cannot set profile %v", err)
				return
			}
			server, err := getServer()
			if err != nil {
				fmt.Println(err)
				return
			}
			path := "/api/v1/machines/" + machine + "/ssh"
			token, err := getToken()
			if err != nil {
				fmt.Println(err)
				return
			}
			u := &url.URL{Scheme: "ws", Host: server, Path: path}
			c, resp, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Authorization": []string{token}})
			if resp != nil && resp.StatusCode == 302 {
				u, _ = resp.Location()
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

			go func() {
				defer func() { done <- true }()
				c.SetReadDeadline(time.Now().Add(pongWait))
				c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
				for {
					mt, r, err := c.NextReader()
					if websocket.IsCloseError(err,
						websocket.CloseNormalClosure,   // Normal.
						websocket.CloseAbnormalClosure, // OpenSSH killed proxy client.
					) {
						fmt.Println(err)
						return
					}
					if err != nil {
						fmt.Printf("nextreader: %v\n", err)
						return
					}
					if mt != websocket.BinaryMessage {
						fmt.Println("binary message \n")
						return
					}
					if _, err := io.Copy(os.Stdout, r); err != nil {
						fmt.Printf("Reading from websocket: %v\n", err)
						return
					}
				}
			}()

			go func() {
				defer func() { done <- true }()
				for {
					var input []byte = make([]byte, 1)
					os.Stdin.Read(input)

					c.SetWriteDeadline(time.Now().Add(writeWait))
					err = c.WriteMessage(websocket.BinaryMessage, input)
					if err != nil {
						fmt.Println("write:", err)
						return
					}
				}
			}()

			go func() {
				defer func() { done <- true }()
				ticker := time.NewTicker(pingPeriod)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
							fmt.Println("ping:", err)
						}
					}
				}
			}()

			<-done
		},
	})

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
	// TODO: Add register commands here.
	mistApiV2Register(false)
	cli.Root.Execute()
}
