#/bin/bash

caddyVersion=2.1.1
sirtunnelVersion=0.1.0

caddyGz=caddy_${caddyVersion}_linux_amd64.tar.gz
curl -O -L https://github.com/caddyserver/caddy/releases/download/v${caddyVersion}/${caddyGz}
tar xvf ${caddyGz}
rm ${caddyGz}
rm LICENSE
rm README.md

sudo setcap 'cap_net_bind_service=+ep' caddy

curl -O -L https://github.com/anderspitman/SirTunnel/releases/download/${sirtunnelVersion}/sirtunnel
chmod +x sirtunnel
