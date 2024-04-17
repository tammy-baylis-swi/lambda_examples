# tammy-test-go

To build and deploy to AWS Lambda function using Mac:

1. `GOOS=linux go build main.go`
2. Send output `main` binary to a zip file.
3. Use output `main.zip` to update a new/existing AWS Lambda function that uses Amazon Linux runtime.
   1. In AWS Console: Lambda function > Code > Code source > Upload from > .zip file.
4. Make sure Lambda function handler points to our code.
   1. In AWS Console: Lambda function > Code > Runtime settings > Handler should be `main`.
5. Invoke Lambda function with the KV `"name": "<some_name>"`