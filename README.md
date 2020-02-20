# Gigamunch!

# Setup
The following programs need to be installed:
  - gcloud
  - golang 1.11
  - protoc 
    - https://github.com/google/protobuf/releases 
    - `go get -u github.com/golang/protobuf/protoc-gen-go`
    - `go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger`
    - `go get -u google.golang.org/grpc`
  - install node LTS: https://nodejs.org/en/
  - `brew install yarn`
  - `yarn global add gulp-cli`
  - `yarn global add @vue/cli`

Do the following in your .bashprofile or .zshrc file:
  - add GOPATH
    - recommended ~/Development/go
  - add to PATH
    - `.` 
    - `/usr/local/mysql/support-files/`
    - `$GOPATH/bin`
    - `~/Development/protoc/bin`

Setting up for web development:
  - run `yarn install`
    - in ./
    - in ./admin/app
    - in ./subserver/web
  - `go get ./cookapi`
  - `go get ./server`
  - `go get ./admin`
  - `go get ./subserver`
  - `app build proto`
  - add private folder
  - setup mysql servers 
    - run `sudo mysql_secure_installation`
    - login to mysql with `sudo mysql -uroot`
    - run following in mysql
      - `uninstall plugin validate_password;`
      - `CREATE USER 'server'@'localhost' IDENTIFIED BY 'gigamunch';`
      - `GRANT ALL PRIVILEGES ON *.* To 'server'@'localhost';`
      - copy, paste, and run ./misc/setup.sql

To run local development:
  - `app serve (admin | server | sub)`

# App Engine Architecture
Here are the modules:
  - default:
    - In the 'server' folder.
    - This module serves landing page.
  - admin:
    - In the 'admin' folder.
  - sub:
    - In the 'subserver' folder.
  - cookapi:
    - In the 'cookapi' folder.