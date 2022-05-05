# airspace-monitoring
Exam project for SRS course

## Development environment setup
### Pre commit hooks
For security and quality reasons, this project makes use of Git pre commit hooks to avoid committing secrets and low quality code files.

The current setup requires the use of [`pre-commit`](https://pre-commit.com/).

Once this repository has been cloned, run `pre-commit install` in the directory. The first commit after the installation of the hooks requires the [Golang](https://go.dev) compiler to be installed.

### Code scan
Code scanning is done using CodeQL and [Semgrep](https://semgrep.dev/).

Semgrep can be installed with `pip install semgrep` (alternatively, [docker can be used](https://semgrep.dev/docs/getting-started/)), scanning is done with 
```
semgrep \
  --config p/security-audit \
  --config p/secrets \
  --config p/ci \
  --config p/javascript \
  --config p/scala \
  --config p/golang \
  --config p/owasp-top-ten \
  --config p/clientside-js \
  --config p/terraform \
  --metrics=off
```

For CodeQL, a guide for installation is [here](https://codeql.github.com/docs/codeql-cli/getting-started-with-the-codeql-cli/). 
To run the checks:
```
codeql database create ./codeql-database --language=<language-you-are-analyzing> --db-cluster
codeql database analyze ./codeql-database/<language-you-are-analyzing> --format=csv --output=./codeql.csv <queries>
```
Local scan with CodeQL requires a GITHUB_TOKEN environment variable to be configured.

Remember to not commit the CodeQL databases!

### Container scan
Currently, scanning containers on the CI pipeline uses [Trivy](https://github.com/aquasecurity/trivy) and [Dockle](https://github.com/goodwithtech/dockle), throught the [official Azure action](https://github.com/Azure/container-scan).

To locally test whether or not the containers will pass the CI test Trivy and Dockle can be run locally. For installation, refer to [this for Trivy](https://aquasecurity.github.io/trivy/v0.27.1/getting-started/installation/) and [this for Dockle](https://github.com/goodwithtech/dockle#installation).

Execution of the tools is simply
```
trivy image image_name:label
```
and
```
VERSION=$(
 curl --silent "https://api.github.com/repos/goodwithtech/dockle/releases/latest" | \
 grep '"tag_name":' | \
 sed -E 's/.*"v([^"]+)".*/\1/' \
) && docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  goodwithtech/dockle:v${VERSION} [YOUR_IMAGE_NAME]
```

### Test of Github Actions Workflows
A useful tool to locally test workflows is [`act`](https://github.com/nektos/act).

`act` shows different behaviours than Github Actions when using matrixes and checks if matrixes are empty and for the values of `${{ github.event.before }}` for new branches.
