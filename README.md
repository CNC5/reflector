# reflector
A management solution for xray that does boring steps for you
- high level abstractions, simplified config syntax
- automatic decoy server setup, content pullable from local directories or docker/oci images

> [!WARNING]
> This tool is for users who know __exactly__ what they are doing,
> many edge cases will result in a non-working setup
> without any errors due to freedoms allowed in configuration <br>
> If you are new to xray I strongly advise against using this tool,
> and instead using a web panel or pure xray-core to gain understanding

## config
Refer to [example.config.yaml](https://github.com/CNC5/reflector/blob/main/example.config.yaml)

## run
```
Usage:
  reflector [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  load        Load a management module
  run         Start the reflector

Flags:
  -d, --debug                     Enable debugging
  -h, --help                      help for reflector
  -r, --reflector-config string   Reflector config location (default "./config.yaml")
```

## build
```
git clone https://github.com/CNC5/reflector
cd reflector
go build
```
