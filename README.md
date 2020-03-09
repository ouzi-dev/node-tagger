# node-tagger

## Overview

node-tagger is a Kubernetes operator that applies specified tags to all aws nodes of the cluster. 

## Deployment
### Prerequisites

The controller requires AWS credentials to be set before deploying it. This is accomplished by creating a secret with name `aws-credentials` in the controller namespace with the following keys:
* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY

For example running the following will create an appropriate secret in the `node-tagger` namespace:
```
kubectl create secret generic aws-credentials --from-literal=AWS_ACCESS_KEY_ID=access_key --from-literal=AWS_SECRET_ACCESS_KEY=secret_access_key --namespace=node-tagger
```

### Required IAM permissions
The operator requires `ec2:CreateTags` and `ec2:DescribeInstance` permissions for the nodes that we are going to tag

### Deploy the operator

Deploy the operator dependencies:
```
kubectl apply -f deploy/service_account.yaml -n node-tagger
kubectl apply -f deploy/role.yaml -n node-tagger
kubectl apply -f deploy/role_binding.yaml -n node-tagger
```

Deploy the operator:
```
kubectl apply -f deploy/deployment.yaml -n node-tagger
```

### Deploying via helm chart

#### Without existing credentials secret
```
helm upgrade --install node-tagger https://github.com/ouzi-dev/node-tagger/releases/download/${VERSION}/node-tagger-${VERSION}.tgz \
    -n node-tagger \
    --set awsCredentials.create=true \
    --set awsCredentials.awsAccessKeyId=access_key \
    --set awsCredentials.awsSecretAccessKey=secret_access_key \
    --set awsCredentials.awsRegion=region \
    --set tagsToApply[0].name=tag1,tagsToApply[0].value=value1
```
Where ${VERSION} is the version you want to install

#### With existing credentials secret
```
helm upgrade --install node-tagger https://github.com/ouzi-dev/node-tagger/releases/download/${VERSION}/node-tagger-${VERSION}.tgz \
    -n node-tagger \
    --set awsCredentials.secretName=aws-credentials \
    --set tagsToApply[0].name=tag1,tagsToApply[0].value=value1
``` 
Where ${VERSION} is the version you want to install

#### Without using a secret (For example using a EKS service account)
```
helm upgrade --install node-tagger https://github.com/ouzi-dev/node-tagger/releases/download/${VERSION}/node-tagger-${VERSION}.tgz \
    -n node-tagger \
    --set awsCredentials.useSecret=false \
    --set tagsToApply[0].name=tag1,tagsToApply[0].value=value1 \
    --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"="arn:aws:iam::123456789:role/my-role-name"
``` 
Where ${VERSION} is the version you want to install
