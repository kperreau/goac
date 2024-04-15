# Go Affected Cache (GOAC)

## About
Go Affected Cache (GOAC) is a tool tailored to optimize the build process for binaries and Docker images in monorepo setups through caching.

It generates and caches a hash based on the list of dependencies and files imported by each project within the monorepo. This functionality prevents unnecessary builds by detecting changes accurately, thus conserving time and resources.

Originally developed to address specific needs in my own project's architecture, GOAC might not directly align with every user's project structure and requirements, making its utility somewhat experimental.

## Features
    Monorepo Efficiency: Specifically designed for monorepos, ensuring efficient handling of multiple interconnected projects.
    Intelligent Change Detection: Utilizes hashes of dependencies and files to determine the necessity of builds.
    Docker Integration: Optimizes Docker image creation by avoiding unnecessary builds and pushes, preventing redundant deployments.

## Installation
Follow these steps to install GOAC in your environment. Adjust as necessary for your specific setup.

```bash
go install github.com/kperreau/goac 
```

## Usage
GOAC offers commands such as affected and list to manage your monorepo effectively.

```
Usage:
  goac [command]

Available Commands:
  affected    List affected projects
  help        Help about any command
  list        List projects

Flags:
  -c, --concurrency int   Max Concurrency (default 4)
      --debug string      Debug files loaded/hashed
  -h, --help              help for goac
  -p, --projects string   Filter by projects name
```

### Checking / Building Affected Projects
```
Usage:
  goac affected [flags]

Examples:
goac affected -t build

Flags:
      --binarycheck     Affected if binary is missing
      --dockerignore    Read docker ignore (default true)
      --dryrun          Dry & run
  -f, --force           Force build
  -h, --help            help for affected
      --stdout          Print stdout of exec command
  -t, --target string   Target

Global Flags:
  -c, --concurrency int   Max Concurrency (default 4)
      --debug string      Debug files loaded/hashed
  -p, --projects string   Filter by projects name

```
Exemples:
```bash
goac affected -t build # build binary of affected project
goac affected -t build -p auth-service,docs # build binaries for auth-service and docs
goac affected -t build --force # build all binaries without checking affected projects
goac affected -t build --debug=name,hashed -p docs # build project docs with debug to display project name and hashed files
```

### Listing Projects
To list all projects configured in your monorepo based on the `.goacproject.yaml`:

```
Usage:
  goac list [flags]

Examples:
goac list

Flags:
  -h, --help   help for list

Global Flags:
  -c, --concurrency int   Max Concurrency (default 4)
      --debug string      Debug files loaded/hashed
  -p, --projects string   Filter by projects name
```

Exemples:
```bash
goac list
```

### Common Options
    --debug [types]: Controls the verbosity of command output, useful for debugging.
                     Available types: name,includes,excludes,local,dependencies

## Contribution
Contributions are welcome! If you'd like to contribute, please follow these steps:

    Fork the project
    Create your feature branch (git checkout -b feature/AmazingFeature)
    Commit your changes (git commit -m 'Add some AmazingFeature')
    Push to the branch (git push origin feature/AmazingFeature)
    Open a Pull Request

## Authors
- [@kperreau](https://www.github.com/kperreau)

## License
Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.