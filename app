#!/bin/bash

################################################################################
# build 
################################################################################
if [[ $1 == "build" ]]; then

  if [[ $* == *proto* ]]; then
    # build protobuf and grpc
    echo "Building Gigamunch-Proto eater api."
    protoc -I Gigamunch-Proto/common/ -I Gigamunch-Proto/eater/ Gigamunch-Proto/common/*.proto Gigamunch-Proto/eater/*.proto --go_out=plugins=grpc:Gigamunch-Proto/eater
  fi
  exit 0
fi

################################################################################
# deploy 
################################################################################

if [[ $1 == "deploy" ]]; then 
  project="gigamunch-omninexus-dev"
  if [[ $* == *--prod* ]] || [[ $* == *-p* ]]; then
    project="gigamunch-omninexus"
  fi
  echo "deploying the following to $project" 
  if [[ $* == *eater* ]]; then
    echo "deploying eater"
    cd eaterapi
    aedeploy gcloud app deploy --project=$project --version=1
    cd ..
  fi
  if [[ $* == *server* ]]; then
    echo "deploying server"
    goapp deploy server/app-stage.yaml
  fi
  if [[ $* == *cook* ]]; then
    echo "deploying cook"
    goapp deploy cookapi/app.yaml
  fi
  exit 0
fi


################################################################################
# serve
################################################################################
if [[ $1 == "serve" ]] || [[ $1 == "" ]]; then
  # setup mysql
  /usr/local/opt/mysql56/bin/mysql.server start
  # create gigamunch database
  cat misc/setup.sql | mysql -uroot
  # start goapp serve
  if [[ $2 == "eater" ]]; then
    echo "Starting eaterapi and server."
    dev_appserver.py --datastore_path ./.datastore endpoint-gigamuncher/app.yaml server/app-dev.yaml
  else
    echo "Starting cookapi and server."
    dev_appserver.py --datastore_path ./.datastore cookapi/app.yaml server/app-dev.yaml
  fi
  # stop mysql
  /usr/local/opt/mysql56/bin/mysql.server stop
fi 

wait
