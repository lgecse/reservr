# Deploying to AWS

User has to have an active session to the banzai-dev-aws account, preferably targeted for the eu-west-1 region.

### Deploy base infra CF template
```
aws cloudformation create-stack \                                                                                                                                                                                                                         ─╯
  --stack-name reservr \
  --template-body file://infra/infra.yaml --capabilities CAPABILITY_NAMED_IAM
```
or 
```
aws cloudformation update-stack \                                                                                                                                                                                                                         ─╯
  --stack-name reservr \
  --template-body file://infra/infra.yaml --capabilities CAPABILITY_NAMED_IAM
```

### Build and push docker image
Image version should be bumped on deploy. Application lambda needs to be redeployed with a new image tag, otherwise it won't update.
```
make IMG=reservr-lambda:0.0.4  docker-build docker-push
```

### Deploy application layer
```
 aws cloudformation create-stack \                                                                                                                                                                                                                         ─╯
   --stack-name reservr-app \
  --template-body file://infra/lambda.yaml --capabilities CAPABILITY_NAMED_IAM
```
or
```
aws cloudformation update-stack \                                                                                                                                                                                                                         ─╯
  --stack-name reservr-app \
  --template-body file://infra/lambda.yaml --capabilities CAPABILITY_NAMED_IAM
```

### Run webhook maintainer script, using the Webex App bearer token
Webex app token in stored in AWS secret store, and can be regenerated at developer.webex.com
The API GW URL is the output of the reservr-app CF stack.

```
cd webhook
go run main.go --token=<app token> --api-endpoint=<AWS API gateway URL>
```