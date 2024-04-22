# tammy-test-go

## Setup local (MacOS):

1. `brew upgrade golang`
2. `go get -u github.com/aws/aws-lambda-go`
3. `go get` the other imports in `main.go`

## Build and deploy to AWS Lambda function (MacOS):

1. Compile an executable named 'bootstrap' (it [must have that name](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-naming)):
   1. x86_64: `GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go`
   1. arm64: `GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go`
2. Send output `bootstrap` binary to a zip file: `zip myFunction.zip bootstrap`
3. Use output `myFunction.zip` to update a new/existing AWS Lambda function that uses Amazon Linux runtime.
   1. In AWS Console: Lambda function > Code > Code source > Upload from > .zip file.
4. Make sure Lambda function handler points to our code.
   1. In AWS Console: Lambda function > Code > Runtime settings > Handler should be `main`.
5. Invoke Lambda function with the KV `"name": "<some_name>"`. Check CloudWatch for any logged errors.

## Optional: Deploy with custom collector config

1. Save a copy of [SW collector config.yaml](https://github.com/solarwinds/opentelemetry-lambda/blob/swo/collector/config.yaml) in this directory.
2. Use the output `bootstrap` binary plus the collector config file to create a zip file: `zip myFunction.zip bootstrap config.yaml`
3. After uploading the zip file, set the function environment variable `OPENTELEMETRY_COLLECTOR_CONFIG_FILE: /var/task/config.yaml`.
   1. In AWS Console: Lambda function > Configuration > Environment variables > Edit button.
