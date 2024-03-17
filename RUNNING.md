# Configuring and Running `kel-agent`

`kel-agent` configuration is primarily done with a YAML file; the default location is
platform-dependent but can be listed with `kel-agent -h`. Many options can also be set with
command-line arguments.

## Websocket Server

Default configuration:

```yaml
websocket:
  address: localhost
  port: 8081
  allowedOrigins:
    - https://forester.radio
```

### Running on localhost

In the simplest case, `kel-agent` is running on the same computer as your radio programs and
browser. In this case, you can have `kel-agent` bind to `address: localhost` which will only allow
programs on the same computer to connect. This is straightforward, safe and the default.

```shell
$ kel-agent
7:19PM INF Serving websocket address=ws://localhost:8081/websocket
```

To use a different port, use `websocket.port` YAML config or the `host` program argument.

```yaml
websocket:
  port: 9988
```

```shell
$ kel-agent
7:19PM INF Serving websocket address=ws://localhost:9988/websocket
```

### Running on another machine

If you want to run your radio programs and `kel-agent` on one computer and your browser on another,
this is possible. There are a couple of approaches. Neither is super easy, which I hope to fix.

NOTE: I do _not_ recommend serving this in a way that's exposed to the internet because there is
_no_ authentication. If exposed to the internet, anyone could potentially initiate transmissions
with your radio.

#### SSH port forwarding

This method is relatively simple and quick to execute, but is more brittle than serving secure
websockets because there is some setup each time you want to use the agent remotely, and
conceptually a little harder. Your remote machine must be running an SSH server for this to work.

On the remote machine with your radio software, run `kel-agent` normally. It can be bound to
`localhost`.

```shell
$ kel-agent
7:19PM INF Serving websocket address=ws://localhost:8081/websocket
```

On the machine with your browser, start a command line and establish an SSH tunnel with port
forwarding:

```shell
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

#### Secure Websocket with TLS

This method needs a little more set up ahead of time, but is easier to use once it's set up.

First, you'll need `kel-agent` to bind to `0.0.0.0` to allow connections from other computers.
Second, due to the mixed content policy which is standard in web browsers, you'll need to specify a
TLS certificate and private key for `kel-agent` to use. The easy way to do this is to use the
[Let's Encrypt](https://letsencrypt.org/) free public service to generate the private key and
certificate for you, signed by LE's certificate authority and recognized by almost all browsers.
Using LE usually assumes that there's a
[web server exposed to the internet](https://letsencrypt.org/docs/challenge-types/#http-01-challenge)
(again, I _don't_ recommend this with `kel-agent`). There's also a
[`dns-01` challenge](https://letsencrypt.org/docs/challenge-types/) if you have a domain name that
the remote computer can be addressed by, even if it's not accessible on the internet.

If Let's Encrypt is not an option for you, you'll need to follow
https://stackoverflow.com/a/60516812/587091 to manually generate your private key and certificate.
In short,

1. generate a CA key and root certificate, then
2. a server key and certificate signing request with the server's hostname,
3. sign the request to generate the server certificate, then finally
4. install the root certificate in your browser's trusted authorities.

Yeah, I really need to make this easier.

```yaml
websocket:
  address: 0.0.0.0
  port: 8081
  cert: /home/joe/.config/kel-agent/fullchain.pem
  key: /home/joe/.config/kel-agent/privkey.pem
```

```shell
$ kel-agent
7:19PM INF Serving websocket address=wss://radio-pi.myhome.net:8081/websocket
```

Notice that the log message doesn't just say `ws://` but `wss://` which means "secure websocket."
Once running this way, your web application can be configured to connect directly to the remote
computer.

### Allowed origins

As part of the same-origin policy which is standard in web browsers, `kel-agent` will only accept
browser connections from certain origins (basically, websites). By default, only the website
`https://forester.radio` plus some local developer addresses are allowed to connect to `kel-agent`,
but this can be customized if others develop web applications that use `kel-agent`. I'm happy to
accept pull requests to expand the default list!

```yaml
websocket:
  allowedOrigins:
    - https://forester.radio
    - https://someother.nifty.app
```

```shell
$ kel-agent
7:19PM INF allowed origins origins=["https://forester.radio","https://someother.nifty.app"]
```

## WSJT-X

`kel-agent` can be used with WSJT-X to automate the process of logging contacts. WSJT-X will attempt
to connect to something listening on UDP port 2237 by default; `kel-agent` listens there and will
pass the contact information to the web application.

```yaml
wsjtx:
  enabled: true
  address: 224.0.0.1
  port: 2237
```

```shell
$ kel-agent
7:19PM INF Listening to WSJT-X on UDP address=224.0.0.1:2237
```

Note that 224.0.0.1 is the multicast address that WSJT-X uses by default on Linux and Mac. On
Windows, `kel-agent` listens by default instead on 127.0.0.1. This matches WSJT-X's behavior.
