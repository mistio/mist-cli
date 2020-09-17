package main

import (
	"os"

	"github.com/danielgtaylor/openapi-cli-generator/apikey"
	"github.com/danielgtaylor/openapi-cli-generator/cli"
	"github.com/spf13/cobra"
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

	// Mock out a profile to be used in the request.
	// cli.Creds.Set("profiles.default.api_key", "e2beabc47d91ba16936c426aec03b83ef0389a9d80832c24662bc9ea3ab8575c")
	cli.Root.AddCommand(completionCmd)
	// TODO: Add register commands here.
	mistApiV2Register(false)
	cli.Root.Execute()
}
