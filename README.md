# airspace-monitoring
Exam project for SRS course

## Development environment setup
### Pre commit hooks
For security and quality reasons, this project makes use of Git pre commit hooks to avoid committing secrets and low quality code files.

The current setup requires the use of [`pre-commit`](https://pre-commit.com/).

Once this repository has been cloned, run `pre-commit install` in the directory. The first commit after the installation of the hooks requires the [Golang](https://go.dev) compiler to be installed.
