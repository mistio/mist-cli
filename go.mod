module gitlab.ops.mist.io/mistio/mist-cli

go 1.15

require (
	github.com/containerd/console v1.0.1
	github.com/gorilla/websocket v1.4.2
	github.com/jmespath/go-jmespath v0.4.0
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/v-pap/trie v0.0.0-20220304164748-f2da6e8bb111
	gitlab.ops.mist.io/mistio/openapi-cli-generator v0.0.0-20220614120433-24df03ca6073
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/h2non/gentleman.v2 v2.0.5
)

replace github.com/krasun/trie => github.com/v-pap/trie v0.0.0-20220302152130-c7abc322b710
