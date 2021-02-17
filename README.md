### Сокращатель ссылок

#### POST запрос /shortify/(ссылка) 

Создает короткий url и возвращает его

#### GET запрос /(короткая ссылка)

Перенаправляет на оригинальную ссылку

### Запуск

docker-compose build

docker-compose up

### Тесты

curl -X POST "http://localhost:9000/shortify/mail.ru"

curl -X POST "http://localhost:9000/shortify/yandex.ru"

curl -X GET "http://localhost:9000/(ссылка полученная из метода post)"
