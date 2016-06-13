# Gigamunch!

# Setup
The following programs need to be installed:
  - golang app engine sdk
  - mysql 5.6 (just use brew for OS X)
  - redis 3.2+

In order to do local development, the following config files are needed:
  - session_config.json

Do the following in your .bashprofile file:
  - export GIGAMUNCH_PRIVATE_DIR = point to the private directory with config files
  - add /usr/local/mysql/support-files/ to PATH

Setting up for development:
  - npm install
  - Install the following packages for Atom:
    - Go-plus
    - linter-eslint
    - laguage-babel

# App Engine Architecture
There are currently two modules.
default:
  - In the 'main' folder.
  - This module serves all front-end related request such as template based page rendering.
gigachefendpoint:
  - In the 'endpoints/gigachef' folder.
gigamuncherendpoint:
  - In the 'endpoint/gigamuncher' folder.

# Styling notes
  - gigachef and gigamuncher variables should be chef and muncher respectively to
    confusion with the packages gigachef and gigamuncher

# Other notes
  - When added a page to the website, the app.yaml and app-shell.html page must be edited or the page will keep reloading.
