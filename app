#!/bin/bash

# Define a timestamp function
timestamp() {
  date +"%T"
}
timestamp

################################################################################
# build
################################################################################
if [[ $1 == "build" ]]; then
  if [[ $* == *admin* ]]; then
    echo "Building admin/app:"
    cd admin/app
    gulp build
    cd ../..
  fi
  if [[ $* == *server* ]]; then
    echo "Building server:"
    cd server
    gulp build
    cd ..
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
    echo "Building Gigamunch-Proto APIs."
    # Eater
    protoc -I Gigamunch-Proto/common/ -I Gigamunch-Proto/eater/ Gigamunch-Proto/common/*.proto Gigamunch-Proto/eater/*.proto --go_out=plugins=grpc:Gigamunch-Proto/eater
    # Shared
    protoc -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I Gigamunch-Proto/shared/ Gigamunch-Proto/shared/*.proto --go_out=plugins=grpc:Gigamunch-Proto/shared
    mv Gigamunch-Proto/shared/github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/shared/*.go Gigamunch-Proto/shared/
    rm -fR Gigamunch-Proto/shared/github.com
    # Admin
    protoc -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I Gigamunch-Proto/admin/ -I Gigamunch-Proto/shared/ Gigamunch-Proto/admin/*.proto --go_out=plugins=grpc:Gigamunch-Proto/admin --swagger_out=logtostderr=true:admin/app
    sed -i 's/"http",//g; s/"2.0",/"2.0",\n"securityDefinitions": {"auth-token": {"type": "apiKey","in": "header","name": "auth-token"}},"security": [{"auth-token": []}],/g' admin/app/AdminAPI.swagger.json
    # Server
    protoc -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I Gigamunch-Proto/server/ -I Gigamunch-Proto/shared/ Gigamunch-Proto/server/*.proto --go_out=plugins=grpc:Gigamunch-Proto/server --swagger_out=logtostderr=true:server
    # Typescript
    gulp build
    # Copy Typescript definitions to folder
    cp Gigamunch-Proto/admin/*.d.ts admin/app/ts/prototypes
    cp Gigamunch-Proto/shared/*.d.ts admin/app/ts/prototypes
  fi
  timestamp
  exit 0
fi

################################################################################
# deploy
################################################################################
if [[ $1 == "deploy" ]]; then
  project="gigamunch-omninexus-dev"
  sqlip="104.154.108.220"
  domain="gigamunch-omninexus-dev.appspot"
  if [[ $* == *--prod* ]] || [[ $* == *-p* ]]; then
    project="gigamunch-omninexus"
    sqlip="104.154.236.200"
    domain="eatgigamunch"
  fi
  echo "Deploying the following to $project"
  if [[ $* == *cook* ]]; then
    echo "Deploying cook:"
    cat cookapi/app.template.yaml | sed "s/PROJECTID/$project/g; s/SQL_IP/$sqlip/g; s/_DOMAIN_/$domain/g" > cookapi/app.yaml
    gcloud app deploy cookapi/app.yaml --project=$project --version=1 --quiet
  fi
  if [[ $* == *admin* ]]; then
    echo "Deploying admin:"
    cat admin/app.template.yaml | sed "s/PROJECTID/$project/g; s/SQL_IP/$sqlip/g; s/_DOMAIN_/$domain/g" > admin/app.yaml
    gcloud app deploy admin/app.yaml --project=$project --version=1 --quiet
  fi
  if [[ $* == *server* ]]; then
    echo "Deploying server:"
    cat server/app.template.yaml | sed "s/PROJECTID/$project/g; s/SQL_IP/$sqlip/g; s/_SERVEPATH_/\/build\/default/g; s/MODULE/default/g; s/_DOMAIN_/$domain/g" > server/app.yaml
    gcloud app deploy server/app.yaml  --project=$project --version=2 --quiet
  fi
  if [[ $* == *queue* ]]; then
    echo "Deploying queue.yaml:"
    gcloud app deploy admin/queue.yaml  --project=$project
  fi
  if [[ $* == *cron* ]]; then
    echo "Deploying cron.yaml:"
    gcloud app deploy admin/cron.yaml  --project=$project
  fi
  timestamp
  exit 0
fi

################################################################################
# serve
################################################################################
if [[ $1 == "serve" ]]; then
  # setup mysql
  if [[ $OSTYPE == "linux-gnu" ]]; then
    service mysql start&
  else
    /usr/local/opt/mysql@5.6/bin/mysql.server start
  fi
  # start goapp serve
  project="gigamunch-omninexus-dev"
  sqlip="104.154.108.220"
  if [[ $2 == "admin" ]]; then
    echo "Starting admin:"
    cat admin/app.template.yaml | sed "s/PROJECTID/$project/g; s/SQL_IP/$sqlip/g; s/_DOMAIN_/$domain/g" > admin/app.yaml
    dev_appserver.py --datastore_path ./.datastore admin/app.yaml&
    cd admin/app
    gulp build&
    gulp watch
    cd ../..
  fi
  if [[ $2 == "server" ]]; then
    echo "Starting server:"
    cat server/app.template.yaml | sed "s/PROJECTID/$project/g; s/_SERVEPATH_//g; s/MODULE/server/g; " > server/app.yaml
    dev_appserver.py --datastore_path ./.datastore server/app.yaml&
    cd server
    gulp build&
    gulp watch
    cd ..
  fi
  # stop mysql

  if [[ $OSTYPE == "linux-gnu" ]]; then
    service mysql stop&
  else
    /usr/local/opt/mysql@5.6/bin/mysql.server stop
  fi
  # kill background processes
  trap 'kill $(jobs -p)' EXIT
  timestamp
  exit 0
fi

################################################################################
# help
################################################################################
if [[ $1 == "help" ]] || [[ $1 == "" ]]; then
  echo "Here are the commands supported by the script:"
  echo -e "\tapp [help|serve|build|deploy]"
  echo -e "\tapp serve [server|admin]"
  echo -e "\tapp build [server|admin|cook|proto]"
  echo -e "\tapp deploy [--prod|-p] [admin|cook|server|queue|cron]"
  timestamp
  exit 0
fi

wait
