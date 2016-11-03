#!/bin/bash

################################################################################
# help
################################################################################
if [[ $1 == "help" ]] || [[ $1 == "" ]]; then
  echo "Here are the commands supported by the script:"
  echo -e "\tapp [help|serve|build|deploy]"
  echo -e "\tapp serve [eater|*]"
  echo -e "\tapp build [app|cook|proto]"
  echo -e "\tapp deploy [--prod|*] [eater|server|cook]"
fi 

################################################################################
# build 
################################################################################
if [[ $1 == "build" ]]; then
  if [[ $* == *app* ]]; then
    echo "Building server/app:"
    cd server/app
    polymer build
    rm -rf build/unbundled
    cd ../..
  fi
  if [[ $* == *cook* ]]; then
    echo "Building server/cook:"
    cd server/cook
    polymer build
    # remove unneccessary files
    rm -rf build/unbundled
    cd ../..
  fi
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
  sqlip="104.154.108.220"
  if [[ $* == *--prod* ]] || [[ $* == *-p* ]]; then
    project="gigamunch-omninexus"
    sqlip="104.154.236.200"
  fi
  echo "Deploying the following to $project" 
  if [[ $* == *eater* ]]; then
    echo "Deploying eater:"
    cat eaterapi/app.yaml.template | sed "s/PROJECT_ID/$project/g; s/SQL_IP/$sqlip/g" > eaterapi/app.yaml
    cd eaterapi
    aedeploy gcloud app deploy --project=$project --version=1
    cd ..
  fi
  if [[ $* == *cook* ]]; then
    echo "Deploying cook:"
    cat cookapi/app.yaml.template | sed "s/PROJECT_ID/$project/g" > cookapi/app.yaml
    goapp deploy cookapi/app.yaml
  fi
  if [[ $* == *server* ]]; then
    echo "Deploying server:"
    cat server/app.yaml.template | sed "s/PROJECT_ID/$project/g; s/_SERVEPATH_/\/build\/bundled/g; s/MODULE/default/g" > server/app.yaml
    goapp deploy server/app.yaml
  fi
  exit 0
fi


################################################################################
# serve
################################################################################
if [[ $1 == "serve" ]]; then
  # setup mysql
  /usr/local/opt/mysql56/bin/mysql.server start
  # create gigamunch database
  cat misc/setup.sql | mysql -uroot
  # start goapp serve
  project="gigamunch-omninexus-dev"
  sqlip="localhost"
  if [[ $2 == "eater" ]]; then
    echo "Starting eaterapi and server."
    cat server/app.yaml.template | sed "s/PROJECT_ID/$project/g; s/_SERVEPATH_//g; s/MODULE/server/g" > server/app.yaml
    cat eaterapi/app.yaml.template | sed "s/PROJECT_ID/$project/g; s/SQL_IP/$sqlip/g" > eaterapi/app.yaml
    dev_appserver.py --datastore_path ./.datastore eaterapi/app.yaml server/app.yaml
  else
    echo "Starting cookapi and server."
    cat cookapi/app.yaml.template | sed "s/PROJECT_ID/$project/g" > cookapi/app.yaml
    cat server/app.yaml.template | sed "s/PROJECT_ID/$project/g; s/_SERVEPATH_//g; s/MODULE/server/g" > server/app.yaml
    dev_appserver.py --datastore_path ./.datastore cookapi/app.yaml server/app.yaml
  fi
  # stop mysql
  /usr/local/opt/mysql56/bin/mysql.server stop
fi 

wait
