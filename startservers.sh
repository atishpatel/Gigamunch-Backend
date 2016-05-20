#!/bin/bash

################################################################################
# start redis
################################################################################
# echo 'starting redis'
# redis-server > ./redislog.tmp &
################################################################################
# setup mysql
################################################################################
/usr/local/opt/mysql56/bin/mysql.server start
# create gigamunch database
cat misc/setup.sql | mysql -uroot
################################################################################
# goapp
################################################################################
echo 'starting goapp'
if [ $1 == "muncher" ]; then
  echo "using muncher and server"
  dev_appserver.py --datastore_path ./.datastore endpoint-gigamuncher/app.yaml server/app-dev.yaml
else
  echo "using chef and server"
  dev_appserver.py --datastore_path ./.datastore endpoint-gigachef/app.yaml server/app-dev.yaml
fi
################################################################################
# clean up
################################################################################
# stop mysql
/usr/local/opt/mysql56/bin/mysql.server stop

# echo 'stopping redis server'
# redis-cli shutdown
wait
