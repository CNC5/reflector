#!/usr/bin/bash

mkdir -p serverops/camo/templates
curl --output - "https://file.cnc5.dev/templates/NextJSExample.tgz" | tar zxC serverops/camo/templates
