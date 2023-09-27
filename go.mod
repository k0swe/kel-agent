module github.com/k0swe/kel-agent

// This version needs to track golang in Debian stable backports (currently bullseye)
// https://packages.debian.org/bullseye-backports/golang
go 1.19

require (
	dario.cat/mergo v1.0.0
	github.com/adrg/xdg v0.4.0
	github.com/gorilla/websocket v1.5.0
	github.com/k0swe/wsjtx-go/v4 v4.0.4
	github.com/rs/zerolog v1.31.0
	github.com/stretchr/testify v1.8.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/leemcloughlin/jdn v0.0.0-20201102080031-6f88db6a6bf2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mazznoer/csscolorparser v0.1.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
