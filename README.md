# nutella-go

A simple Go app to cache search results from the free Nominatim API using Redis and Docker.

### Run

```
docker compose up --build
```

##### or

```
go run .
```

### Health Check

```
localhost:8080/health
```


### Search with Cache
```
localhost:8080/cache?q=istanbul
```
