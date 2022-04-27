# airspace-monitoring
Exam project for SRS course

## Development environment setup
### Pre commit hooks
For security and quality reasons, this project makes use of Git pre commit hooks to avoid committing secrets and low quality code files.

The current setup requires the use of [`pre-commit`](https://pre-commit.com/).

Once this repository has been cloned, run `pre-commit install` in the directory. The first commit after the installation of the hooks requires the [Golang](https://go.dev) compiler to be installed.

### Test of Github Actions Workflows
A useful tool to locally test workflows is [`act`](https://github.com/nektos/act).

`act` shows different behaviours than Github Actions when using matrixes and checks if matrixes are empty and for the values of `${{ github.event.before }}` for new branches.
