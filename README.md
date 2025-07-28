# Reflector

A basic unit for reflecting xray, with some useful abstractions for clustering.

## Goals
- Easily usable xray that is configured according to best practices
- Abstractions for clustering and HA

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [License](#license)

## Features

- Automatic selfsteal with letsencrypt certificates
- User forwarding between connections
- Multiple outbound users
- Link type outbounds

## Requirements
- git
- curl
- python3.13

## Installation
> [!IMPORTANT]
> `download_binaries.sh` downloads binaries from my server: a sing-box that was most recently tested with reflector, and a static build of nginx that also was most recently tested to work with reflector. You can use your own binaries, by specifying `--nginx-bin`|`--xray-bin` or by putting them to `serverops/bin/{sing-box,nginx}`
```bash
git clone https://github.com/CNC5/reflector.git
cd reflector
bash download_binaries.sh
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
```

## Usage
```bash
python -m serverops --help
```
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

## Configuration
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
        - name: bob
          uuid: bb3bb3bb-bbbb-bbbb-bbbb-bb3bbbbbb3bb
          flow: ''
          short_id: b0b01d
      private_key: <privkey>
      camo:
        type: local
        template: NextJSExample
        fqdn: vacuums.lemao.xyz
        issuer:
          type: letsencrypt # can also be 'selfsigned' for testing, no email required
          email: an0n@mozmail.com

  outbounds:
    - name: vless-eu
      type: link
      link: "vless://..."
    - name: vless-out2
      type: vless
      server: 100.1.1.1
      server_port: 443
      server_name: unsus.eu
      fingerprint: chrome
      users:
        - name: bob
          uuid: b0b01df0-0b11-1111-1111-111111111111
          flow: ''
          short_id: 111111dc
        - name: alice
          uuid: a11ce01d-2222-2222-2222-222222222222
          flow: ''
          short_id: 222222ac
      public_key: <pubkey>
    - name: direct
      type: direct

  routes:
    - user: alice
      outbound: vless-eu
    - user: bob
      outbound: vless-out2

  metrics: # TODO
    port: 12345
    listen: localhost

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
