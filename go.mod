module gitlab.ops.mist.io/mistio/mist-cli

go 1.15

require (
	github.com/containerd/console v1.0.1
	github.com/gorilla/websocket v1.4.2
	github.com/jmespath/go-jmespath v0.4.0
	github.com/manifoldco/promptui v0.9.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/v-pap/trie v0.0.0-20220304164748-f2da6e8bb111
	gitlab.ops.mist.io/mistio/openapi-cli-generator v0.0.0-20220708125651-795a9083b395
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	gopkg.in/h2non/gentleman.v2 v2.0.5
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/client-go v0.24.2
)

replace github.com/krasun/trie => github.com/v-pap/trie v0.0.0-20220302152130-c7abc322b710
