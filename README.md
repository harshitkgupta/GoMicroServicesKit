### Setup Environment
* Set up GOPATH
echo $GOPATH
export GOPATH=~/go

* Setup Directory Structure
bin- executable
pkg- package objects made by libraries
src- all source code


### Dependecy Management
* brew install dep
* dep version
* dep init  
create Gopkg.toml and Gopkg.lock files
* dep ensure -add github.com/gorilla/mux


## Vendor
* go get github.com/kardianos/govendor

migrate from godep to vendor
*  govendor migrate


Heroku
* heroku login
* heroku create
