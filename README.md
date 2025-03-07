# Book Server

# Clone the Repository

``` bash
$ git clone <your-repo-url>
$ cd book-server
```

# Build and Run the Go Server Locally

``` bash
$ go build -o book-server
$ ./book-server
```

#  Run with Docker

``` bash
$ docker build -t book-server:latest .
$ docker run -p 3000:3000 book-server:latest
```
# Deploy to Kubernetes with Kind

``` bash
$ kind create cluster
$ kind load docker-image book-server:latest
$ kubectl apply -f deployment.yaml
$ kubectl apply -f service.yaml
```

# Access the Server

``` bash
$ kubectl port-forward service/book-server 3000:8080
```
