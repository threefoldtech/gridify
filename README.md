# Gridify
A tool used to deploy projects using threefold grid 3 

## Build
Clone the repo and run the following command inside the repo directory:
```bash
go build .
```
## Requirements
- gridify uses [ginit](https://github.com/rawdaGastan/ginit) so `.Procfile` and `.env` must exist in root directory of your project
- the project github repository must be public

## Usage
First [Build](#build) gridify then move the binary to any of $PATH directories, for example:
```bash
mv gridify /usr/local/bin
```

Login using your mnemonics and specify which grid network to deploy on by running:
```bash
gridify login
```

Use `gridify` to deploy your project and specify the ports you want gridify to assign domains to:
```bash
gridify deploy --ports 80,8080
```

To destroy deployed projects run:
```bash
gridify destroy
```

## Demo
See [gridify-demo](https://github.com/AbdelrahmanElawady/gridify-demo)

## Supported Projects Languages and Tools
- go 1.18
- python 3.10.10
- node 16.17.1
- npm 8.10.0
- caddy