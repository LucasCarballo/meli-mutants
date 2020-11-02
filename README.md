# Requirements

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [Redis](https://redis.io/topics/quickstart)

# Running tests

To run tests, browse project folder and run tests in project:

```bash
    cd ~/meli-mutants
    go test ./...
```

# Running the project

To run the test you will need to have running a redis container.

```bash 
    docker run --name container-redis -d redis
```

Then, from the project folder, just run the `main.go` file

```bash
    cd ~/meli-mutants
    go run main.go
```
