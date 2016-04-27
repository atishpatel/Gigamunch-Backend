#!/bin/bash

################################################################################
# start redis
################################################################################
echo 'starting redis'
redis-server > /dev/null &
################################################################################
# setup mysql
################################################################################
/usr/local/opt/mysql56/bin/mysql.server start
# create gigamunch database
cat misc/create_gigamunch_datbase.sql | mysql -uroot
# create live_posts table
mysql -uroot gigamunch < misc/create_live_posts_table.sql
# TODO create user for get and create mysql and one of cron delete
################################################################################
# goapp
################################################################################
echo 'starting goapp'
dev_appserver.py --datastore_path ./.datastore endpoint-gigachef/app.yaml server/app-dev.yaml # endpoint-gigamuncher/app.yaml
################################################################################
# clean up
################################################################################
echo 'stopping mysql'
/usr/local/opt/mysql56/bin/mysql.server stop

echo 'stopping redis server'
redis-cli shutdown
wait
