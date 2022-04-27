# airspace-monitoring
Exam project for SRS course

## Development environment setup
### Pre commit hooks
For security and quality reasons, this project makes use of Git pre commit hooks to avoid committing secrets and low quality code files.

The current setup requires the use of [`pre-commit`](https://pre-commit.com/).

Once this repository has been cloned, run `pre-commit install` in the directory. The first commit after the installation of the hooks requires the [Golang](https://go.dev) compiler to be installed.

### Container scan
Currently, scanning containers on the CI pipeline uses [Trivy](https://github.com/aquasecurity/trivy) and [Dockle](https://github.com/goodwithtech/dockle), throught the [official Azure action](https://github.com/Azure/container-scan).

To locally test whether or not the containers will pass the CI test Trivy and Dockle can be run locally. For installation, refer to [this for Trivy](https://aquasecurity.github.io/trivy/v0.27.1/getting-started/installation/) and [this for Dockle](https://github.com/goodwithtech/dockle#installation).

Execution of the tools is simply
```
trivy image image_name:label
```
and
```
TODO
```

