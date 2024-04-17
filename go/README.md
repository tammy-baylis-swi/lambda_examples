# tammy-test-go

## Setup local (MacOS):

1. `brew upgrade golang`
2. `go get -u github.com/aws/aws-lambda-go`

## Build and deploy to AWS Lambda function (MacOS):

1. Compile an executable named 'bootstrap' (it [must have that name](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-naming)):
   1. x86_64: `GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go`
   1. arm64: `GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go`
2. Send output `bootstrap` binary to a zip file: `zip myFunction.zip bootstrap`
3. Use output `myFunction.zip` to update a new/existing AWS Lambda function that uses Amazon Linux runtime.
   1. In AWS Console: Lambda function > Code > Code source > Upload from > .zip file.
4. Make sure Lambda function handler points to our code.
   1. In AWS Console: Lambda function > Code > Runtime settings > Handler should be `main`.
5. Invoke Lambda function with the KV `"name": "<some_name>"`