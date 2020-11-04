![Release](https://github.com/k0swe/kel-agent/workflows/Release/badge.svg)

# kel-agent

An agent program for translating between various amateur radio installed programs and WebSockets.
This will allow the creation of cloud-based amateur radio applications while using integration
points only available through local processes. 

![Architecture](architecture.svg)

At first this will support receiving status and log messages from WSJT-X. Planned support includes
`rigctld` and Ham Radio Deluxe for transceiver remote control.

## Running on localhost

In the simplest case, `kel-agent` is running on the same computer as your radio programs and
browser. In this case, you can have `kel-agent` bind to `localhost` which will only allow programs
on the same computer to connect. This is straightforward, safe and the default.

```
$ kel-agent
2020/10/10 18:50:32 kel-agent ready to serve at ws://localhost:8081
```

## Running on another machine

If you want to run your radio programs and `kel-agent` on one computer and your browser on another,
this is possible. There are a couple of approaches. Neither is super easy, which I hope to fix.

NOTE: I do *not* recommend serving this in a way that's exposed to the internet because there is
*no* authentication. If exposed to the internet, anyone could potentially initiate transmissions
with your radio.

### SSH port forwarding

This method is relatively simple and quick to execute, but is more brittle than serving secure
websockets because there is some setup each time you want to use the agent remotely, and
conceptually a little harder. Your remote machine must be running an SSH server for this to work.

On the remote machine with your radio software, run `kel-agent` normally. It can be bound to
`localhost`.

```
$ kel-agent
2020/10/10 18:50:32 kel-agent ready to serve at ws://localhost:8081
```

On the machine with your browser, start a command line and establish an SSH tunnel with port
forwarding:

```
$ ssh -N -L localhost:8081:localhost:8081 radio-pi
```

The first `localhost:8081` means "on this (browser) machine, bind to port 8081 and only expose to
`localhost` so other computers can't use it." The second `localhost:8081` means "once you log into 
the remote computer, start forwarding traffic to port 8081 on its (remote) `localhost`." Finally, 
`radio-pi` in my example is the remote hostname which is running `kel-agent` and the SSH server.

The command will look like it's not doing anything; just let it run, and the tunnel will stay open.

Now your web application can be configured to connect to `localhost`. Traffic bound for
`localhost:8081` will get securely forwarded to the remote machine. Both the browser and `kel-agent`
think they're talking to local processes, and you won't get mixed content warnings.

### Secure Websocket with TLS

This method needs a little more set up ahead of time, but is easier to use once it's set up.

First, you'll need `kel-agent` to bind to `0.0.0.0` to allow connections from other computers. 
Second, due to the mixed content policy which is standard in web browsers, you'll need to specify a
TLS certificate and private key. Eventually I hope to make this easy, but for now you'll need to
follow https://stackoverflow.com/a/60516812/587091. In short,

1. generate a CA key and root certificate, then
2. a server key and certificate signing request with the server's hostname,
3. sign the request to generate the server certificate, then finally
4. install the root certificate in your browser's trusted authorities.
 
Yeah, I really need to make this easier. 

```
$ kel-agent -host 0.0.0.0:8081 -key server.key -cert server.crt
2020/10/10 19:05:39 kel-agent ready to serve at wss://0.0.0.0:8081
```

Once running this way, your web application can be configured to connect directly to the remote
computer.

## Allowed origins

As part of the same-origin policy which is standard in web browsers, `kel-agent` will only accept
browser connections from certain origins (basically, websites). By default, only the website
`https://log.k0swe.radio` plus some local developer addresses are allowed to connect to `kel-agent`,
but this can be customized if others develop web applications that use `kel-agent`. I'm happy to 
accept pull requests to expand the default list!

```
$ kel-agent -origins "https://log.k0swe.radio,https://someother.nifty.app"
2020/10/10 19:18:52 Allowed origins are [https://log.k0swe.radio https://someother.nifty.app]
```
