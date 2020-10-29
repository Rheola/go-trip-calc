# сервис  расчета поездки


## Локальный запуск
``` 
   docker-compose -f docker-compose-local.yml  up -d  
   cp dist.env app/.env
```
Отредактировать ```app/.env```
```
   cd app 
   go mod download
   go run main.go
```

## Запуск  в Docker-контейнере
``` 
   cp dist.env .env 
```
Отредактировать ```.env```
``` 
   docker-compose -f    up -d  
```
  
  
Использование  
## Запрос на расчет 

`` POST /routes``
``header 'Content-Type: application/json'``
```json
{
    "from": {
        "lat": 45.028738,
        "lon": 38.968064
    },
    "to": {
        "lat": 45.052641,  
        "lon": 38.958389
    }
}
```

### Ответ
#### OK
```json
{
    "code": 201,
    "message": "4"
}
```
#### Fail
Header 400 Bad Request
```json 
{
    "code": 400,
    "message": "wrong 'from' param: latitude must be a number between -90 and 90"
}
```
### Example 

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

## Получить результат

`` GET /routes/<id>``

```curl
  curl --location --request GET 'http://localhost:8000/routes/9'
```
### Ответ OK
```json
{
    "code": 200,
    "message": {
        "distance": 3606,
        "duration": 746
    }
}
```
``distance`` - метры
``duration`` - секунды

### Рано
```Header 425 Too early```

```json
{
    "code": 425,
    "message": "Waiting. Route not calking yet"
}
```

### Не найдено 
```Header 404```

```json
{
    "code": 404,
    "message": "Not found"
}
```
### Ошибки расчета
```Header 200```

```json
{
    "code": 500,
    "message": "Calc error"
}
```