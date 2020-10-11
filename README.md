# kel-agent

An agent program for translating between various amateur radio installed programs and WebSockets.
This will allow the creation of cloud-based amateur radio applications while using integration
points only available through local processes. 

![Architecture](architecture.svg)

At first this will support receiving status and log messages from WSJT-X. Planned support includes
`rigctld` and Ham Radio Deluxe for transceiver remote control.

## Running

In the simplest case, `kel-agent` is running on the same computer as your radio programs and
browser. In this case, you can have `kel-agent` bind to `localhost` which will only allow programs
on the same computer to connect. This is straightforward, safe and the default.

```
$ kel-agent
2020/10/10 18:50:32 kel-agent ready to serve at ws://localhost:8081
```

If you want to run your radio programs and `kel-agent` on one computer and your browser on another,
this is possible. First, you'll need to bind to `0.0.0.0` to allow connections from other computers.
Second, you'll need to specify a TLS certificate and private key. Eventually I hope to make this
easy, but for now you'll need to follow https://stackoverflow.com/a/60516812/587091. In short,

1. generate a CA key and root certificate, then
2. a server key and certificate signing request,
3. sign the request to generate the server certificate, then finally
4. install the root certificate in your browser's trusted authorities.
 
Yeah, I really need to make this easier. 

```
$ kel-agent -host 0.0.0.0:8081 -key server.key -cert server.crt
2020/10/10 19:05:39 kel-agent ready to serve at wss://0.0.0.0:8081
```
