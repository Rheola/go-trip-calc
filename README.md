# go-trip-calc

``` 
 docker network create trip-api
```

``` 
    docker-compose up -d
```


```curl
curl --location --request POST 'http://localhost:8080/routes' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": {
        "lat": 45.057010,
        "lon": 38.993252
    },
    "to": {
        "lat": 40.040582,
        "lon": 39.030845
    }
}'
```


curl --location --request GET 'http://localhost:8080/routes/1'




