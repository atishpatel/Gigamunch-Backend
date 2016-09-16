# Gigamunch!

# Setup
The following programs need to be installed:
  - golang app engine sdk
  - mysql 5.6 (just use brew for OS X)
  - nmp
  - bower

Do the following in your .bashprofile file:
  - export GIGAMUNCH_PRIVATE_DIR = point to the private directory with config files
  - add /usr/local/mysql/support-files/ to PATH

Setting up for web development:
  - npm install (/)
  - bower install (/server/app, /server/gigachef)

Run local development:
  - sh app.sh

# App Engine Architecture
There are currently three modules.
default:
  - In the 'main' folder.
  - This module serves all front-end related request such as template based page rendering.
endpoint-gigachef:
  - In the 'endpoints/gigachef' folder.
endpoint-gigamuncher:
  - In the 'endpoint/gigamuncher' folder.

# Other notes
  - When added a page to the website, the app.yaml and app-shell.html page must be edited or the page will keep reloading.
