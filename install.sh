#/bin/bash

caddyVersion=2.1.1
sirtunnelVersion=0.1.0

echo Download Caddy
caddyGz=caddy_${caddyVersion}_linux_amd64.tar.gz
curl -s -O -L https://github.com/caddyserver/caddy/releases/download/v${caddyVersion}/${caddyGz}
tar xf ${caddyGz}

echo Clean up extra Caddy files
rm ${caddyGz}
rm LICENSE
rm README.md

echo Enable Caddy to bind low ports
sudo setcap 'cap_net_bind_service=+ep' caddy

echo Download sirtunnel binary
curl -s -O -L https://github.com/anderspitman/SirTunnel/releases/download/${sirtunnelVersion}/sirtunnel
chmod +x sirtunnel
