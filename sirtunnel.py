#!/usr/bin/env python3

import sys
import json
import time
from urllib import request, error


def check_tunnel_health(tunnel_id):
    check_url = 'http://127.0.0.1:2019/id/' + tunnel_id
    try:
        request.urlopen(check_url)
    except error.HTTPError as e:
        if e.code == 404:  # Not Found
            return False
    except error.URLError:
        return False
    return True


def delete_tunnel(tunnel_id):
    print("Cleaning up tunnel")
    delete_url = 'http://127.0.0.1:2019/id/' + tunnel_id
    req = request.Request(method='DELETE', url=delete_url)
    request.urlopen(req)


if __name__ == '__main__':

    host = sys.argv[1]
    port = sys.argv[2]
    tunnel_id = host + '-' + port

    caddy_add_route_request = {
        "@id": tunnel_id,
        "match": [{
            "host": [host],
        }],
        "handle": [{
            "handler": "reverse_proxy",
            "upstreams":[{
                "dial": ':' + port
            }]
        }]
    }

    body = json.dumps(caddy_add_route_request).encode('utf-8')
    headers = {
        'Content-Type': 'application/json'
    }
    create_url = 'http://127.0.0.1:2019/config/apps/http/servers/sirtunnel/routes'
    req = request.Request(method='POST', url=create_url, headers=headers, data=body)
    request.urlopen(req)

    print("Tunnel created successfully")

    while True:
        try:
            if not check_tunnel_health(tunnel_id):
                print("Tunnel disconnected unexpectedly")
                delete_tunnel(tunnel_id)
                break
            time.sleep(1)
        except KeyboardInterrupt:
            delete_tunnel(tunnel_id)
            break
