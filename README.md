# Gigamunch!

# Setup
The following programs need to be installed:
  - gcloud
  - golang 1.8
  - `brew install mysql@5.6`
  - protoc 
    - https://github.com/google/protobuf/releases 
    - `go get -u github.com/golang/protobuf/protoc-gen-go`
    - `go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger`
    - `go get github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis`
  - `brew install yarn`
  - `yarn global add gulp-cli`

Do the following in your .bashprofile or .zshrc file:
  - add GOPATH
    - recommended ~/Development/go
  - add to PATH
    - `.` 
    - `/usr/local/mysql/support-files/`
    - `$GOPATH/bin`
    - `~/Development/protoc/bin`

Setting up for web development:
  - `app build proto`
  - `go get ./server`
  - `go get ./cookapi`
  - run `yarn install`
    - in ./
    - in ./admin/app
    - in ./driver/app
  - setup mysql servers by running ./misc/setup.sql

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
