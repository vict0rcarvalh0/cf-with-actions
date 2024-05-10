#!/bin/bash
release_type=$1
export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
export AWS_SESSION_TOKEN=""
export AWS_DEFAULT_REGION=""

AWS_ECR_PASSWORD=$(aws ecr get-login-password)

sudo docker stop (docker ps -aq)

sudo docker system prune -f

# Login no Docker com credenciais do ECR
echo $AWS_ECR_PASSWORD | sudo docker login --username AWS --password-stdin 713621535342.dkr.ecr.us-east-1.amazonaws.com
sudo docker pull 713621535342.dkr.ecr.us-east-1.amazonaws.com/builds:$release_type

sudo docker run -d --env-file .env -p 8080:8080 713621535342.dkr.ecr.us-east-1.amazonaws.com/builds:$release_type