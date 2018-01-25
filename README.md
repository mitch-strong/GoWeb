# GoWeb
Simple Golang Web Server

## How To Install and Run
```
1. go get -u github.com/mitch-strong/GoWeb/Web
2. cd $GOPATH/go/src/github.com/mitch-strong/GoWeb
3. docker build -t goweb ./
```

### Keycloak
```
docker run -p 8080:8080 -name=keycloak -d -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin jboss/keycloak-examples
```
1. Create a Client
2. Create a user in this realm
3. Edit allowed redirect URLS
4. Edit constants in main.go and rebuild docker image

## Docker 
```
docker run -d --name=GoWeb --link keycloak -p 3000:3000 -v $HOME/go/src/github.com/mitch-strong/GoWeb/Web  $(docker images -q goweb)
```

NOTE:  When connecting to keycloak the main.go file constants will have to be changed to match the client id and secret of the keycloak client created.  Keycloak must be hosted on port 8080

## Minikube
```
minikube start
kubectl config use-context minikube
minikube dashboard
kubectl run goweb --image=goweb:latest --port=3000 
kubectl expose deployment goweb --type=LoadBalancer
minikube service goweb
```
```
rm -rf ~/.minikube  //Run this if minikube starts with error
```

NOTE:  Ended up logging into docker and pushing image onto repo, then using that image instead.  Hosts no problem
```
docker login --username=yourhubusername
docker push mitchellstrong/goweb:latest
```
## Minikube Docker Repo
```
minikube start
kubectl config use-context minikube
minikube dashboard
kubectl run goweb --image=docker.io/mitchellstrong/goweb:latest --port=3000 
kubectl expose deployment goweb --type=LoadBalancer
minikube service goweb
kubectl scale deployments/goweb --replicas=4
```
Port Forwarding
```
ssh -L localhost:8080:192.168.99.100:30220 -N 127.0.0.1
```