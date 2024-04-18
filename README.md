# Go Affected Cache (GOAC)

## 📖 About
Go Affected Cache (GOAC) is a tool tailored to optimize the build process for binaries and Docker images in monorepo setups through caching.

It generates and caches a hash based on the list of dependencies and files imported by each project within the monorepo. This functionality prevents unnecessary builds by detecting changes accurately, thus conserving time and resources.

Originally developed to address specific needs in my own project's architecture, GOAC might not directly align with every user's project structure and requirements, making its utility somewhat experimental.

## 🌟 Features
    Monorepo Efficiency: Specifically designed for monorepos, ensuring efficient handling of multiple interconnected projects.
    Intelligent Change Detection: Utilizes hashes of dependencies and files to determine the necessity of builds.
    Docker Integration: Optimizes Docker image creation by avoiding unnecessary builds and pushes, preventing redundant deployments.

## 🛠 Installation
Follow these steps to install GOAC in your environment. Adjust as necessary for your specific setup.

```bash
go install github.com/kperreau/goac@latest
```

## ⚙️ Configuration
Configuring GOAC is straightforward. You need to create a `.goacproject.yaml` file and place it at the root of each project/directory that requires a build.

### File example:
```yaml
# .goacproject.yaml

version: 1.0      # Please do not modify this value; keep it set to 1.0
name: goac        # Specify the name of your project, service, or application here
target:           # GOAC currently supports two targets: 'build' and 'build-image'
  build:          # This target compiles the Go binary
    exec:
      cmd: go     # The command to execute for compilation; 'go' in this case
      params:     # Parameters to be added; the final command will be: go build -ldflags="-s -w" -o ./goac goac
        - build
        - -ldflags=-s -w
        - -o
        - "{{project-path}}/{{project-name}}"
        - "{{project-path}}"
  build-image:    # This target builds the Docker image
    envs:
      - key: PROJECT_PATH
        value: "{{project-path}}"
    exec:
      cmd: ./_scripts/build-image.sh # Shell script to execute for building the image
```

To see what the script that builds the image of this project looks like, take a look at this example: [build-image.sh](./_scripts/build-image.sh)

### Variables
The configuration file interprets variables that will automatically be replaced by their values at runtime, it works for env value and params.

| Variable           | Type     | Description             |
|:-------------------| :------- |:------------------------|
| `{{project-name}}` | `string` | The name of the project |
| `{{project-path}}` | `string` | The path of the project |

## 🚀 Usage
GOAC offers commands such as affected and list to manage your monorepo effectively.

**⚠️ Important: To works properly, GOAC must be executed from the root directory of your project, where the `go.mod` file is located.**

```
Usage:
  goac [command]

Available Commands:
  affected    List affected projects
  completion  Generate the autocompletion script for the specified shell
  discover    List discovered projects
  help        Help about any command
  list        List projects

Flags:
  -h, --help   help for goac
```

### Checking / Building Affected Projects
```
List projects affected by recent changes based on GOAC cache.

Usage:
  goac affected [flags]

Examples:
goac affected -t build

Flags:
      --binarycheck       Affected if binary is missing
  -c, --concurrency int   Max Concurrency (default 4)
      --debug string      Display some data to debug
      --dockerignore      Read docker ignore (default true)
      --dryrun            Dry & run
  -f, --force             Force build
  -h, --help              help for affected
  -p, --projects string   Filter by projects name
      --stdout            Print stdout of exec command
  -t, --target string     Target
```
#### Debug Options
```
--debug [types]: Controls the verbosity of command output, useful for debugging.
                 Available types: name,includes,excludes,local,dependencies
```

#### Exemples:
```bash
goac affected -t build # build binary of affected project
goac affected -t build -p auth,user # build binaries for auth and user service
goac affected -t build --force # build all binaries without checking affected projects
goac affected -t build --debug=name,hashed -p docs # build project docs with debug to display project name and hashed files
```

### Listing Projects
List all projects configured in your monorepo based on the `.goacproject.yaml`:

```
Use it to list all your projects configured with GOAC.

Usage:
  goac list [flags]

Examples:
goac list

Flags:
  -c, --concurrency int   Max Concurrency (default 4)
  -h, --help              help for list
  -p, --projects string   Filter by projects name
```
#### Exemples:
```bash
goac list
goac list -p goac
```


### Discover Projects
GOAC can explore your repository to identify potential projects and automatically generate a default `.goacproject.yaml` configuration file per project.

```
Use it to discovering projects and create default config files.

Usage:
  goac discover [flags]

Examples:
goac discover

Flags:
  -c, --create   Create project config files
  -f, --force    Force creation file if already exist
  -h, --help     help for discover
```

#### Exemples:
```bash
goac discover
goac discover -c
```

#### Automatic Project Naming and Configuration
GOAC automatically generates a project name in each configuration file based on the directory path.
Additionally, the default configuration executes the [build-image.sh](./_scripts/build-image.sh) script for the build-image target.
Please note that this script is not created by default in your repository. You will need to either modify the configuration or create your own script.

## 📘 Note
The `.dockerignore` and its interpretation are crucial for GOAC.
It allows excluding all unused files, especially those likely to be generated and impact the cache, thereby potentially affecting the project indefinitely.
Make sure to exclude files that are not necessary, such as the `.git`, `.idea`, `.vscode`, etc.

For reference, see this [.dockerignore](.dockerignore)


## 👨‍💻 Contribution
Contributions are welcome! If you'd like to contribute, please follow these steps:

    Fork the project
    Create your feature branch (git checkout -b feature/AmazingFeature)
    Commit your changes (git commit -m 'Add some AmazingFeature')
    Push to the branch (git push origin feature/AmazingFeature)
    Open a Pull Request

## 📝 Authors
- [@kperreau](https://www.github.com/kperreau)

## 📄 License
Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.