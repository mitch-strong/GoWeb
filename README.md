# GoWeb
Simple Golang Web Server

## How To Install and Run
```
1. go get -u github.com/mitch-strong/GoWeb/Web
2. cd $GOPATH/go/src/github.com/mitch-strong/GoWeb
3. docker build -t goweb ./
```
## Docker 
```
4. docker run -d --name=GoWeb -p 8080:8080 -v $GOPATH/go/src/github.com/mitch-strong/GoWeb/Web  $(docker images -q goweb)
```
## Minikube
```
4. minikube start
5. kubectl config use-context minikube
6. minikube dashboard
7. kubectl run goweb --image=goweb:latest --port=8080 
8. kubectl expose deployment goweb --type=LoadBalancer
9. minikube service goweb
```
rm -rf ~/.minikube  - Run this if minikube starts with error
docker login --username=yourhubusername

NOTE:  Ended up logging into docker and pushing image onto repo, then using that image instead.  Hosts no problem
docker login --username=yourhubusername
docker push mitchellstrong/goweb:latest
```
4. minikube start
5. kubectl config use-context minikube
6. minikube dashboard
7. kubectl run goweb --image=docker.io/mitchellstrong/goweb:latest --port=8080 
8. kubectl expose deployment goweb --type=LoadBalancer
9. minikube service goweb
```
