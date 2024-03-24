# What is it?

If you have a webserver running on one computer (say your development laptop),
and you want to expose it securely (ie HTTPS) via a public URL, SirTunnel
allows you to easily do that.

# How do you use it?

If you have:

* A SirTunnel [server instance](#running-the-server) listening on port 443 of
  `example.com`.
* A copy of the sirtunnel.py script available on the PATH of the server.
* An SSH server running on port 22 of `example.com`.
* A webserver running on port 8080 of your laptop.

And you run the following command on your laptop:

```bash
ssh -tR 9001:localhost:8080 example.com sirtunnel.py sub1.example.com 9001
```

Now any requests to `https://sub1.example.com` will be proxied to your local
webserver.


# How does it work?

The command above does 2 things:

1. It starts a standard [remote SSH tunnel][2] from the server port 9001 to
   local port 8080.
2. It runs the command `sirtunnel.py sub1.example.com 9001` on the server.
   The python script parses `sub1.example.com 9001` and uses the Caddy API to
   reverse proxy `sub1.example.com` to port 9001 on the server. Caddy
   automatically retrieves an HTTPS cert for `sub1.example.com`.

**Note:** The `-t` is necessary so that doing CTRL-C on your laptop stops the
`sirtunnel.py` command on the server, which allows it to clean up the tunnel
on Caddy. Otherwise it would leave `sirtunnel.py` running and just kill your
SSH tunnel locally.


# How is it different?

There are a lot of solutions to this problem. In fact, I've made something of
a hobby of maintaining [a list][0] of the ones I've found so far.

The main advantages of SirTunnel are:

* Minimal. It leverages [Caddy][1] and whatever SSH server you already have
  running on your server. Other than that, it consists of a 50-line Python
  script on the server.  That's it. Any time you spend learning to customize
  and configure it will be time well spent because you're learning Caddy and
  your SSH server.
* 0-configuration. There is no configuration on the server side.  Not even CLI
  arguments.
* Essentially stateless. The only state is the certs (which is handled entirely
  by Caddy) and the tunnel mappings, which are ephemeral and controlled by the
  clients.
* Automatic HTTPS certificate management. Some other solutions do this as well,
  so it's important but not unique.
* No special client is required. You can use any standard SSH client that
  supports remote tunnels. Again, this is not a unique feature.


# Running the server

Assuming you already have an ssh server running, getting the SirTunnel server
going consists of simply downloading a copy of Caddy and running it with the
provided config. Take a look at [`install.sh`](./install.sh) and
[`run_server.sh`](./run_server.sh) for details.

**Note:** Caddy needs to bind to port 443, either by running as root (not
recommended), setting the `CAP_NET_BIND_SERVICE` capability on the Caddy binary
(what the `install.sh` script does), or changing `caddy_config.json` to bind
to a different port (say 9000) and using something like iptables to forward
to that port.

# Future Features

SirTunnel is intended to be a minimal tool. As such, I'm unlikely to add many
features moving forward. However, the simplicity makes it easier to modify
for your needs. If you find a feature missing, maybe one of the forks below
has what you need or gives you some ideas:

* https://github.com/matiboy/SirTunnel adds functionality to help multiple
  users avoid overwriting each others' tunnels.
* https://github.com/daps94/SirTunnel implements tunnel checking and auto
  cleanup to address issue with stale tunnels.


[0]: https://github.com/anderspitman/awesome-tunneling

[1]: https://caddyserver.com/

[2]: https://www.ssh.com/ssh/tunneling/example#remote-forwarding
