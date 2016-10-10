#!/bin/bash

################################################################################
# build 
################################################################################
if [[ $1 == "build" ]]; then
  # build protobuf and grpc
  echo "Building Gigamunch-Proto eater api."
  protoc -I Gigamunch-Proto/common/ -I Gigamunch-Proto/eater/ Gigamunch-Proto/common/*.proto Gigamunch-Proto/eater/*.proto --go_out=plugins=grpc:Gigamunch-Proto/eater
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