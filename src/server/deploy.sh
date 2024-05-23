#!/bin/bash

# The commands to run on the remote server
COMMANDS="
    sudo apt-get update --fix-missing -y && sudo apt-get install build-essential make -y

    # Check if Docker is installed
    if ! type docker > /dev/null 2>&1; then
        # Update packages
        sudo apt update;
        # Install prerequisites
        sudo apt install -y apt-transport-https ca-certificates curl software-properties-common;
        # Add Docker's official GPG key
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -;
        # Add Docker's stable repository
        sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable\";
        # Update package database
        sudo apt update;
        # Install Docker
        sudo apt install -y docker-ce;
        # Add the current user to the docker group
        sudo usermod -aG docker \$(whoami);
        # Log out and log back in
        su - \$(whoami);
    fi;

    # Check if Node.js is installed
    if ! type node > /dev/null 2>&1; then
        # NVM installation
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

        # Load NVM
        export NVM_DIR=\"\$HOME/.nvm\"
        [ -s \"\$NVM_DIR/nvm.sh\" ] && \. \"\$NVM_DIR/nvm.sh\"
        [ -s \"\$NVM_DIR/bash_completion\" ] && \. \"\$NVM_DIR/bash_completion\"

        # Install Node.js
        nvm install --lts;
        nvm use --lts;
    fi;
    
eval \$(ssh-agent -s) && ssh-add ~/.ssh/id_ed25519

    # Clone or pull latest changes from repository
    if [ -d \"deploy-buddy\" ]; then
        cd deploy-buddy;
        git pull;
    else
        git clone git@github.com:Inteli-College/2024-1B-T03-ES10-G02.git deploy-buddy;
        cd deploy-buddy;
    fi;

    # Check out development branch and start application
    # git checkout development;
    git checkout infra/docker

    # cd into the project directory
    cd src/server
    
    make prod;
"

# SSH into the server and run the commands
ssh -i $SSH_KEY ubuntu@$IP "$COMMANDS"