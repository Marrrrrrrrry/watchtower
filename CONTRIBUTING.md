## Prerequisites
To contribute code changes to this project you will need the following development kits.
 * [Go](https://golang.org/doc/install)
 * [Docker](https://docs.docker.com/engine/installation/)
 
As watchtower utilizes go modules for vendor locking, you'll need at least Go 1.26.
You can check your current version of the go language as follows:
```bash
  ~ $ go version
  go version go1.26.0 darwin/amd64
```


## Checking out the code
Do not place your code in the go source path.
```bash
git clone git@github.com:<yourfork>/watchtower.git
cd watchtower
```

## Building and testing
watchtower is a go application and is built with go commands. The following commands assume that you are at the root level of your repo.
```bash
go build                               # compiles and packages an executable binary, watchtower
go test ./... -v                       # runs tests with verbose output
./watchtower                           # runs the application (outside of a container)
```

If you dont have it enabled, you'll either have to prefix each command with `GO111MODULE=on` or run `export GO111MODULE=on` before running the commands. [You can read more about modules here.](https://github.com/golang/go/wiki/Modules)

To build a Watchtower image of your own, use the self-contained Dockerfiles. As the main Dockerfile, they can be found in `dockerfiles/`:
- `dockerfiles/Dockerfile.dev-self-contained` will build an image based on your current local Watchtower files.
- `dockerfiles/Dockerfile.self-contained` will build an image based on current Watchtower's repository on GitHub.

e.g.:
```bash
sudo docker build . -f dockerfiles/Dockerfile.dev-self-contained -t marrrrrrrrry/watchtower # to build an image from local files
```

## Dependency changes and security
Dependency updates should arrive through Dependabot or be prepared by a maintainer. Pull requests from third parties that claim to fix a vulnerability deserve heightened scrutiny.

When reviewing a pull request that modifies `go.mod` or `go.sum`:

- Accept version changes of module paths that are already in use. A pull request that adds or rewrites a module path, for example by replacing an official module with a similarly named one, must be rejected.
- Verify reported vulnerabilities against the official Go vulnerability database at [pkg.go.dev/vuln](https://pkg.go.dev/vuln) before acting on them. Claims made in issue reports or pull request descriptions are not authoritative.
- Note that GitHub Actions workflows triggered by `pull_request_target` run with access to repository secrets. Such workflows must never check out code from the pull request.