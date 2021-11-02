 
# Okteto go sdk
## Deploy preview sample

In this example we are going to use okteto CLI as sdk to deploy a preview environment.

This script will create a new preview environment by getting the options from the following environment variables:
- OKTETO_URL: This URL is the url for Okteto Enterprise 
- OKTETO_TOKEN: This is the okteto token to connect to Okteto Enterprise
- REPOSITORY: The repository to deploy (defaults to the current repository)
- BRANCH: The branch to deploy (defaults to the current branch)
- SCOPE: The scope of preview environment to create. Accepted values are ['personal', 'global']
- SOURCE_URL: The URL of the original pull/merge request.
- FILENAME: Relative path within the repository to the manifest file (default to okteto-pipeline.yaml or .okteto/okteto-pipeline.yaml)
- VARIABLES: Set variable. Format is key1=value1;key2=value2

For example if we want to deploy [movies](https://github.com/okteto/movies) example into a preview environment on [Okteto Cloud](https://cloud.okteto.com):

`OKTETO_URL=https://cloud.okteto.com OKTETO_TOKEN=<YOUR-API-TOKEN> REPOSITORY=https://github.com/okteto/movies BRANCH=master SCOPE=personal go run my-preview-<YOUR-OKTETO-NAME>`
