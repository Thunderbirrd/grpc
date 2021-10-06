# grpc
gRPC обёртка над сайтом https://www.rusprofile.ru/

Для локального запуска:  
`go run main.go` в папке проекта  
`go run client.go` в папке `client`   
На порту 8081 появится HTTP эндпоинты, на 8080 - grpc

Для запуска в `Docker`: `docker-compose up` в папке проекта, на порту 8081 будет `swagger-ui`

Формат запроса:
`GET` `localhost:8081/inn/{inn}`
