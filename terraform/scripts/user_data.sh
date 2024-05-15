#!/bin/bash -xe
LOG_FILE="/var/log/userdata_execution.log"
{
  # Define the working directory
  WORK_DIR="${work_dir}"

  # Update the system and install packages
  sudo yum update -y
  sudo yum install -y golang git

  # Clone the repository
  echo "Cloning the repository..." &>> $LOG_FILE
  git clone ${repo_url} $WORK_DIR &>> $LOG_FILE

  # Check if the directory exists
  if [ ! -d "$WORK_DIR" ]; then
    echo "Error: The repository was not cloned correctly to $WORK_DIR." &>> $LOG_FILE
    exit 1
  fi

  echo "Listing the cloned directory content:" &>> $LOG_FILE
  ls -l $WORK_DIR &>> $LOG_FILE

  # Move to the project directory
  echo "Changing to the project directory..." &>> $LOG_FILE
  cd $WORK_DIR/sample-app-go &>> $LOG_FILE

  # Check if the directory exists
  if [ $? -ne 0 ]; then
    echo "Error: The project directory does not exist in $WORK_DIR/sample-app-go." &>> $LOG_FILE
    exit 1
  fi

  # Initialize and run the Go application
  echo "Initializing and running the Go application..." &>> $LOG_FILE
  export GOCACHE="$HOME/.cache/go-build"
  go mod init sample-app-go &>> $LOG_FILE
  go mod tidy &>> $LOG_FILE
  go build &>> $LOG_FILE
  ./sample-app-go &>> $LOG_FILE
} 2>&1 | tee -a $LOG_FILE
