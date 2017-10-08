# Gigamunch!

# Setup
The following programs need to be installed:
  - gcloud
  - mysql 5.6 (just use brew for OS X)
  - `brew install yarn`
  - `yarn add global gulp-cli`

Do the following in your .bashprofile or .zsh file:
  - add /usr/local/mysql/support-files/ to PATH

Setting up for web development:
  - `yarn install`
  - give 'app' executable permission
    - `chmod 755 app`

To run local development:
  - `app serve (admin | server)`

# App Engine Architecture
Here are the modules:
  - default:
    - In the 'server' folder.
    - This module serves landing page.
  - admin:
    - In the 'admin' folder.
  - driver:
    - In the 'driver' folder.
  - subscriber:
    - In the 'subscriber' folder.

# Other notes
  - When added a page to the website, the app.yaml and app-shell.html page must be edited or the page will keep reloading.
