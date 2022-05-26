# airspace-monitoring
Exam project for SRS course

## Development environment setup
### Pre commit hooks
For security and quality reasons, this project makes use of Git pre commit hooks to avoid committing secrets and low quality code files.

The current setup requires the use of [`pre-commit`](https://pre-commit.com/).

Once this repository has been cloned, do:
- run `pre-commit install` in the main directory

The first commit after the installation of the hooks requires the [Golang](https://go.dev) compiler to be installed.

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

### Terraform
Manual deploy requires `terraform` and `gcloud` installed.

To update the version of the services images that will be used by Terraform: `./scripts/update-tags.sh`

Pre-deploy operations:
- create a project on Google Cloud, memorize the id of the project
- login to `gcloud`
- run the `act-gcp-apis.sh` script to abilitate the necessary Google APIs: `./scripts/act-on-gcp-apis.sh enable <project_id>`
- run `./scripts/create-service-accounts.sh <project_id>` to create the following service accounts, with the relative permissions:
  - `aircraft-list` with role "Cloud Datastore Viewer" and role "Log Writer"
  - `airspace-daily-history` with role "Cloud Datastore Viewer" and role "Log Writer"
  - `airspace-monthly-history` with role "Cloud Datastore Viewer" and role "Log Writer"
  - `airspace-history-calculator` with role "Cloud Datastore User" and role "Log Writer"
  - `flink-sa` with role "Cloud Datastore User", role "Pub/Sub Publisher" and role "Log Writer"
- push the following Docker images to a Google Cloud Repository, in the same region as the where the system will be deployed, using `./scripts/build-everything.sh <project_id> <region> <docker_repo_name>`, after having set gcloud as the credential helper for docker:
  - `aircraft_info`
  - `aircraft_list`
  - `aircraft_positions` 
  - `airspace_daily_history` 
  - `airspace_monthly_history` 
  - `airspace_history_calculator` 
  - `websocket_endpoints`
  - `web_ui`
  - `states_source` 
- create a Google Storage bucket and write its name as the "bucket" value of the 'backend "gcs"' object into `terraform/main.tf`
- create the default project in Firestore with the `airspace` collection
- create inside the `terraform` directory a file `secrets.auto.tfvars` in which the following variables (described in `terraform/variables.tf`) are declared, one per line, with the `var_name = var_value` syntax:
  - `project_id`
- change the variables declared in `terraform/variables.auto.tfvars` as necessary


Once the pre-deploy operations are done, execute:
```
gcloud init # to log in into Google Cloud, should already have done this
gcloud auth application-default login # to set up credentials for terraform
terraform -chdir=terraform init # to initialize the terraform state, necessary
terraform -chdir=terraform fmt # to correctly format the terraform files, not strictly necessary
terraform -chdir=terraform validate # to check the syntax of the Terraform files
terraform -chdir=terraform plan # to elaborate what the apply command will do
terraform -chdir=terraform apply # to deploy the system
terraform -chdir=terraform destroy # to destroy every change done by the apply command to the cloud
```

To configure `kubectl` and explore the cluster state (see https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl#apt_1):
```
gcloud container clusters get-credentials gke-1 --zone "us-central1"
```

The account running the terraform commands requires the following roles:
- IAM Project Admin
- Service Account Admin
- Editor

### How to use Gatling
Download [gatling](https://gatling.io/open-source/), unzip the bundle in your directory of choice, copy the content of `gatling_simulations` into `simulations` in the unzipped directory, run 
```
JAVA_OPTS="-Dusers=10 -Dramp=3600 -Dbase_url=<our_base_url>" ./bin/gatling.sh
```
Change parameters as needed.

This requires Java to be installed.

Good tool for mocking apis: https://beeceptor.com/
