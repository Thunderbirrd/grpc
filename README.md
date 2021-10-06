# grpc
gRPC обёртка над сайтом https://www.rusprofile.ru/

Для локального запуска:  
`go run main.go`  
`go run client.go`  
На порту 8081 появится HTTP эндпоинты, на 8080 - grpc

Для запуска в `Docker`: docker-compose up, на порту 8081 будет `swagger-ui`

Формат запроса:
`GET` `localhost:8081/inn/{inn}`
