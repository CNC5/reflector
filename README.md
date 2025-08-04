# Reflector

A basic unit for reflecting xray, with some useful abstractions for clustering.

## Goals
- Easily usable vless that is configured according to best practices
- Abstractions for clustering and HA

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Preparation](#preparation)
- [Configuration](#configuration)
- [License](#license)

## Features
- Automatic selfsteal with letsencrypt certificates
- User forwarding between connections
- Multiple outbound users
- Link type outbounds

## Preparation
It is assumed you have a VPS and a domain that has an A record to this server IP address

Perform ONLY ONE of docker or bare-metal

#### Docker
```bash
git clone https://github.com/CNC5/reflector.git
cd reflector
bash download_camo_templates.sh
```

You can proceed to [Configuration](#configuration)

#### Bare-metal
> [!IMPORTANT]
> `download_binaries.sh` downloads binaries: a sing-box that was most recently tested with reflector, and a static build of nginx that also was most recently tested to work with reflector. You can use your own binaries, by specifying `--nginx-bin`|`--xray-bin` or by putting them to `serverops/bin/{sing-box,nginx}`.
###### Requirements
Have these packages
- git (install)
- curl (install)
- python3.13 (interpreter)
- certbot (for letsencrypt certs)

For Ubuntu install with `apt install`

For Alpine install with `apk add`

```bash
git clone https://github.com/CNC5/reflector.git
cd reflector
bash download_binaries.sh
bash download_camo_templates.sh
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
```

Test that everything works
```bash
python -m serverops --help
```
Example output
```
usage: __main__.py [-h] [--tmp TMP] [-c CONFIG] [--pid-file PID_FILE] [--nginx-bin NGINX_BIN] [--xray-bin XRAY_BIN]
                   [--camo-dir CAMO_DIR] [-d] [-s SIGNAL]

options:
  -h, --help            show this help message and exit
  --tmp TMP             directory for tmp storage
  -c, --config CONFIG   config file, in cwd
  --pid-file PID_FILE   pid file, in the tmp directory
  --nginx-bin NGINX_BIN
                        nginx binary
  --xray-bin XRAY_BIN   xray binary
  --camo-dir CAMO_DIR   camo templates dir
  -d, --debug
  -s, --signal SIGNAL   send a signal to the operator
```

You can proceed to [Configuration](#configuration)

## Configuration
Copy the config below and paste it into `config.yaml`
```yaml
apiVersion: v1
kind: Reflector
spec:
  inbounds:
    - name: vless-in
      type: vless
      listen: 0.0.0.0
      listen_port: 443
      users:
        - name: alice
          uuid: a4aaaa4a-aaaa-aaaa-aaaa-a4aaaaaaaa4a
          flow: ''
          short_id: a11ce01d
      private_key: <privkey>
      camo:
        type: local
        template: NextJSExample
        fqdn: vacuums.lemao.xyz
        issuer:
          type: letsencrypt # can also be 'selfsigned' for testing, no email required
          email: an0n@mozmail.com

  outbounds:
    - name: direct
      type: direct

#    - name: vless-eu
#      type: link
#      link: "vless://..."
#    - name: vless-out2
#      type: vless
#      server: 100.1.1.1
#      server_port: 443
#      server_name: unsus.eu
#      fingerprint: chrome
#      users:
#        - name: bob
#          uuid: b0b01df0-0b11-1111-1111-111111111111
#          flow: ''
#          short_id: 111111dc
#        - name: alice
#          uuid: a11ce01d-2222-2222-2222-222222222222
#          flow: ''
#          short_id: 222222ac
#      public_key: <pubkey>

  routes:
    - user: alice
      outbound: direct

  metrics: # TODO
    port: 12345
    listen: localhost

```

Edit it to yout liking:
- leave desired inbounds
- generate new user ids with uuidgen or with `docker run -it --rm --entrypoint /app/serverops/bin/sing-box ghcr.io/cnc5/reflector:latest-alpine generate uuid` (docker) or with `serverops/bin/sing-box generate uuid` (bare-metal). Note both keys, you will need them for client configuration.
- generate new private key and public key pairs with `docker run -it --rm --entrypoint /app/serverops/bin/sing-box ghcr.io/cnc5/reflector:latest-alpine generate reality-keypair` (docker) or with `serverops/bin/sing-box generate reality-keypair` (bare-metal).
- for camo select the template name you want to be set up (NextJSExample is available as a preset), change fqdn to a domain name that you pointed to the server.
- set issuer type to letsencrypt for production and selfsigned for testing.
- set issuer email to the email you want be presented to letsencrypt.
- leave desired outbounds, you likely want only the direct if setting up an exit server.
- add routes for users you set in inbounds, use their `name`s.

Now you can run
```bash
docker compose up -d
```
OR
```bash
python -m serverops
```

## License
```
Copyright (C) 2025 by CNC5 <z1xs4xg62@mozmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

In addition, no derivative work may use the name or imply association
with this application without prior consent.
```
