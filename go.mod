module github.com/k0swe/kel-agent

// This version needs to track golang in Debian stable backports (currently bookworm)
// https://packages.debian.org/bookworm-backports/golang
go 1.21

require (
	dario.cat/mergo v1.0.1
	github.com/adrg/xdg v0.5.3
	github.com/gorilla/websocket v1.5.3
	github.com/invopop/jsonschema v0.13.0
	github.com/k0swe/wsjtx-go/v4 v4.1.3
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.10.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/leemcloughlin/jdn v0.0.0-20201102080031-6f88db6a6bf2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mazznoer/csscolorparser v0.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
