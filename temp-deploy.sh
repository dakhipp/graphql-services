#!/bin/bash

# login to ecs
aws ecr get-login --no-include-email --region us-west-2 | sh

# build and push main image based on dockerfile
docker build -f deployments/docker/graphql/prod.Dockerfile -t staging/graphql .
docker tag staging/graphql:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/staging/graphql:latest
docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/staging/graphql:latest

# build and push auth image based on dockerfile
docker build -f deployments/docker/auth/prod.Dockerfile -t staging/auth .
docker tag staging/auth:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/staging/auth:latest
docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/staging/auth:latest

# build and push migrations image based on dockerfile
docker build -f deployments/docker/migrations/prod.Dockerfile -t staging/migrations .
docker tag staging/migrations:latest 690303654955.dkr.ecr.us-west-2.amazonaws.com/staging/migrations:latest
docker push 690303654955.dkr.ecr.us-west-2.amazonaws.com/staging/migrations:latest

# use aws cli to update service
# 1. cluster will stay static
# 2. service will stay static
# 3. can update task definition number if needed
# 4. force-new-deployment is needed if you do not change the task version number
aws --region us-west-2 ecs \
  update-service \
  --cluster arn:aws:ecs:us-west-2:690303654955:cluster/staging-ecs-cluster \
  --service arn:aws:ecs:us-west-2:690303654955:service/staging-graphql-web \
  --task-definition arn:aws:ecs:us-west-2:690303654955:task-definition/staging_graphql \
  --force-new-deployment
