# nutella-go

A simple Go caching app with Redis and Docker. Using the free search api of Nominatim.

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
