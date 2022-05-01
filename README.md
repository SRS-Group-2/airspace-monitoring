# airspace-monitoring
Exam project for SRS course

## Development environment setup
### Pre commit hooks
For security and quality reasons, this project makes use of Git pre commit hooks to avoid committing secrets and low quality code files.

The current setup requires the use of [`pre-commit`](https://pre-commit.com/).

Once this repository has been cloned, run `pre-commit install` in the directory. The first commit after the installation of the hooks requires the [Golang](https://go.dev) compiler to be installed.

### Security and quality check of Scala code
```
mvn compile spotbugs:spotbugs
```
If a Sarif output is needed, `mvn compile spotbugs:spotbugs '-DsarifOutput=true'`

### states_sources
Java 17, Maven and Docker are required to run this:
```
mvn clean package
docker-compose up -d
```
#### OpenSky API
To use the OpenSky API in Java/Scala for development and compiling it is necessary to build and install them locally. Do:
```
git clone https://github.com/openskynetwork/opensky-api
cd opensky-api/java
mvn clean install
cd ../..
rm -rf opensky-api
```
