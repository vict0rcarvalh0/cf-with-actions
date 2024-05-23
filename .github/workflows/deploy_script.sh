#!/bin/bash

release_type=$1
AWS_ACCESS_KEY_ID=$2
AWS_SECRET_ACCESS_KEY=$3
AWS_SESSION_TOKEN=$4

export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN
export AWS_DEFAULT_REGION="us-east-1"

AWS_ECR_PASSWORD=$(aws ecr get-login-password)

# Stop all running Docker containers
sudo docker stop $(docker ps -aq)

# Remove all unused containers, networks, images, and optionally, volumes
sudo docker system prune -f

# Login to Docker with ECR credentials
echo $AWS_ECR_PASSWORD | sudo docker login --username AWS --password-stdin 730335199760.dkr.ecr.us-east-1.amazonaws.com

# Pull the Docker image from ECR
sudo docker pull 730335199760.dkr.ecr.us-east-1.amazonaws.com/builds:$release_type

# Run the Docker container
sudo docker run -d --env-file .env -p 8080:8080 730335199760.dkr.ecr.us-east-1.amazonaws.com/builds:$release_type
