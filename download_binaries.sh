#!/usr/bin/bash

mkdir -p serverops/bin
curl -o serverops/bin/nginx https://file.cnc5.dev/nginx/1.27.5/nginx
curl -o serverops/bin/sing-box https://file.cnc5.dev/sing-box/1.12.0/sing-box
chmod +x serverops/bin/nginx serverops/bin/sing-box
