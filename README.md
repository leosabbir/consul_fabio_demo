# Installation of GoLang
The services I have written are based on golang. If you want to run it, you will need to install golang. 
Run following commands to install golang.

```
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go
```

### Set GOPATH
Set Golang Paths.
```
sudo mkdir /go
export GOPATH=/go
export PATH=$PATH:$GOPATH/bin
```

# Installation OF Consul
Service discovery using consul is the purpose of this demonstration. There will be consul agent running as server and agents running as client.

```
cd /usr/local/bin
sudo wget https://releases.hashicorp.com/consul/1.5.0/consul_1.5.0_linux_amd64.zip
sudo apt install unzip
sudo unzip consul_1.5.0_linux_amd64.zip
sudo rm consul_1.5.0_linux_amd64.zip
```

# Run Consul as Server
Consul server can be run with following:

```
cd ~
mkdir -p consul-config/server/
vim consul-config/server/config.json
```
Use following configuration
```
    {
        "bootstrap": true,
        "server":true,
        "log_level": "DEBUG",
        "enable_syslog": true,
        "enable_script_checks": true,
        "datacenter": "dc1",
        "addresses": {
                "http": "0.0.0.0"
        },
        "bind_addr": "172.31.17.96", // private IP
        "node_name": "ConsulServer",
        "data_dir": "~/consul-data",
        "ui_dir": "~/consul-ui",
        "acl_datacenter": "dc1",
        "acl_master_token": "123456789",
        "acl_default_policy": "allow",
        "encrypt": "pXoaLOJ816mO+da8y8zrsg=="
    }
```

Setup data directory and run.
```
cd ~
mkdir consul-data
sudo consul agent -config-dir ~/consul-config/server -ui
````

# Run Consul as Client
Set up data directory and run client:
```
cd ~
mkdir consul-data
sudo consul agent -data-dir=consul-data -bind=172.31.39.98 -join=172.31.17.96 -encrypt=pXoaLOJ816mO+da8y8zrsg==
```
  **- bind = host private IP**

  **- join = ip of host where consul server is running**

# Installation of Fabio
Installation details can be found on this Reference: https://github.com/fabiolb/fabio/wiki/Installation

```
go get github.com/fabiolb/fabio 
./fabio
```

# Run User Service
Run user service that will register to consul. You might need to run
```
go get
```
inside user_service to download dependencies.

When dependency is met, you can run the service by passing LOCAL_IP enviroment variable as
```
LOCAL_IP=<host_ip> go run main.go
```
inside user_service.

# Run Storage Service
Run storage service that will register to consul. You might need to run
```
go get
```
inside storage_service to download dependencies.

When dependency is met, you can run the service by passing LOCAL_IP environment variable as
```
LOCAL_IP=<host_ip> go run main.go
```
inside user_service.

# Access from gateway

Run gateway.
```
LOCAL_IP=<host_ip> go run main.go
```
inside gateway.


### User service and gateway need to have properly set *fabiolocation* to send request to fabio.
