[![Go Report Card](https://goreportcard.com/badge/github.com/k0swe/kel-agent)](https://goreportcard.com/report/github.com/k0swe/kel-agent)
[![Release](https://github.com/k0swe/kel-agent/workflows/Release/badge.svg)](https://github.com/k0swe/kel-agent/releases/latest)
[![Release version](https://img.shields.io/github/v/release/k0swe/kel-agent)](https://github.com/k0swe/kel-agent/releases/latest)

# <img src="https://raw.githubusercontent.com/k0swe/kel-agent/main/assets/radio.k0swe.Kel_Agent.svg" width="100px" alt="kel-agent logo"> kel-agent

`kel-agent` is a protocol bridge/translator between amateur radio software and WebSocket clients;
it does **not** use AI. Despite the name, "agent" here refers to a background service, not an
AI agent.

This bridge was built to support https://github.com/k0swe/forester but can be used by any web
application that needs to communicate with amateur radio installed programs.

```mermaid
flowchart LR
    Internet([Internet]) --- Application[Web Application]
    subgraph Computer
        subgraph Browser[Web Browser]
            Application
        end
        agent[kel-agent]
        WSJTX[WSJT-X]
        rigctld[rigctl/Hamlib]
        HRD[HRD]
        etc[...]
    end
    Application ---|websocket| agent
    agent ---|UDP| WSJTX
    agent ---|lib| rigctld
    agent -.->|UDP| HRD
    agent -.-> etc
    style agent fill:#8ecfff,stroke:#1f2937,stroke-width:2px
```

This currently supports communication with WSJT-X, and now officially supports rig control via
`rigctl`/Hamlib. Ham Radio Deluxe support remains planned.

To get started using `kel-agent`, download an appropriate executable from the
[latest release](https://github.com/k0swe/kel-agent/releases/latest). Windows, Mac, Debian/Ubuntu
Linux and Raspberry Pi installers are available.

See the [Running documentation](RUNNING.md) for how to configure, execute and serve `kel-agent`.

## Acknowledgements

The wire logo for `kel-agent` was created by [Freepik](https://www.flaticon.com/authors/freepik) on
[Flaticon](https://www.flaticon.com).
