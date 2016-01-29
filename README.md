# Gigamunch!


# Setup
In order to do local development, there are some config files you need to download from
the project.

In your .bashprofile, export GIGAMUNCH_PRIVATE_DIR to point to the private directory

# App Engine Architecture
There are currently two modules.
default:
  - In the 'main' folder.
  - This module serves all front-end related request such as template based page rendering.
gigachefendpoint:
  - In the 'endpoints/gigachef' folder.
gigamuncherendpoint:
  - In the 'endpoint/gigamuncher' folder.
