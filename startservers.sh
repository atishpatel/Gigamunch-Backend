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
# create live_meals table
mysql -uroot gigamunch < misc/create_live_meals_table.sql
# TODO create user for get and create mysql and one of cron delete
################################################################################
# goapp
################################################################################
echo 'starting goapp'
# goapp serve server/app.yaml endpoint-gigachef/app.yaml endpoint-gigamuncher/app.yaml
goapp serve endpoint-gigamuncher/app.yaml
################################################################################
# clean up
################################################################################
echo 'stopping mysql'
/usr/local/opt/mysql56/bin/mysql.server stop
# kill all subprocesses
trap 'kill $(jobs -p)' EXIT
