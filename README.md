# Reflector

A basic unit for reflecting xray, with some useful abstractions for clustering.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- Automatic selfsteal with letsencrypt certificates
- User forwarding between connections
- Multiple outbound users
- Link type outbounds

## Usage
```bash
python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
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

## Example configuration
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
