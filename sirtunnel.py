#!/usr/bin/env python3

import sys
import json
import time
from urllib import request
from urllib.error import HTTPError
import argparse
import logging

def delete(tunnel_id, caddy_api):
    delete_url = '{}/id/{}'.format(caddy_api, tunnel_id)
    req = request.Request(method='DELETE', url=delete_url)
    request.urlopen(req)

if __name__ == '__main__':
    parser = argparse.ArgumentParser("SirTunnel", description="An easy way to securely expose a webserver running on one computer via a public URL")
    parser.add_argument("--replace", action="store_true", help="Replace the domain if already part of the routes")
    parser.add_argument("--no-duplicates", action="store_true", help="Don't allow duplicate; either abort or replace depending on --replace")
    parser.add_argument("--check-availability", action="store_true", help="Checks that port or domaine might already be in use")
    parser.add_argument("--debug", action="store_true", help="Additional logs")
    parser.add_argument("--check", action="store_true", help="Check each second whether entry still exists")
    parser.add_argument("--caddy-api", default="http://127.0.0.1:2019", help="Caddy's admin api")
    parser.add_argument("domain", help="External domain name")
    parser.add_argument("tunnel_port", type=int, help="The tunnel port between your computer and the server per the command ssh -tr <tunnel_port>:localhost:8000")

    args = parser.parse_args()
    host = args.domain
    port = str(args.tunnel_port)
    debug = args.debug

    caddy_api = args.caddy_api

    logging.basicConfig(level="DEBUG" if debug else "INFO")

    LOGGER = logging.getLogger(__name__)
    LOGGER.debug("Log level set to debug")

    tunnel_id = host + '-' + port

    LOGGER.debug("Tunnel id build %s", tunnel_id)

    headers = {
        'Content-Type': 'application/json'
    }
    if not args.check_availability:
        LOGGER.warn("Skipping checks that domain might already be in use")
    else:
        LOGGER.info("Checking domain and ports availability")
        req = request.Request(method='GET', url="{}/config/apps/http/servers/sirtunnel/routes".format(caddy_api), headers=headers)
        outcome = request.urlopen(req).read().decode('utf-8')
        routes = json.loads(outcome)
        for route in routes:
            domain, _, used_port = route['@id'].partition('-')
            if domain == host or used_port == port:
                LOGGER.error("Host or port already in use on route: %s:%s", host, used_port)
                if args.no_duplicates:
                    if not args.replace:
                        LOGGER.critical("Duplicate entry not allowed. Aborting")
                        sys.exit(1)
                    else:
                        delete(route['@id'], caddy_api)
                        LOGGER.warning("Entry %s has been deleted", route['@id'])
                else:
                    LOGGER.warning("Adding entry despite duplicate. This is likely to cause problems")

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
    
    create_url = '{}/config/apps/http/servers/sirtunnel/routes'.format(caddy_api)
    req = request.Request(method='POST', url=create_url, headers=headers)
    request.urlopen(req, body)

    print("Tunnel created successfully")

    while True:
        try:
            time.sleep(1)
            # Quick check that the tunnel still exists
            if args.check:
                req = request.Request(method='GET', url=caddy_api + '/id/' + tunnel_id)
                response = request.urlopen(req)
        except KeyboardInterrupt:
            print("Cleaning up tunnel")
            delete_url = caddy_api + '/id/' + tunnel_id
            req = request.Request(method='DELETE', url=delete_url)
            request.urlopen(req)
            break
        except HTTPError as ex:
            LOGGER.debug(str(ex))
            LOGGER.critical("Domain entry does not exist anymore. Aborting")
            break
