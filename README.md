## Служба поздравлений с днем рождения

### Обзор

Сервис поздравлений с днем рождения предназначен для удобной отправки поздравлений сотрудникам с днем рождения. Пользователи могут регистрироваться, входить в систему, подписываться на уведомления о дне рождения других пользователей и отписываться от них. Сервис также включает в себя ежедневное задание для проверки наличия дней рождения и отправки уведомлений по электронной почте.
Рассылка будет идти с вашего smtp почтового сервера, нужно указать данные в docker compose файле в разделе с переменными окружения.
### Поднять сервис

    ```bash
    docker compose up
    ```
В корневой папке проекта.

### Api эндпоинты
### Регистрация

#### POST /registration

    ```bash
    curl -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "birthday": "1990-01-01"
    }'
    ```

#### POST /login

    ```bash
    curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{
    "email": "test@example.com",
    "password": "password"
    }'
    ```

#### POST /subscribe

    ```bash
    curl -X POST http://localhost:8080/subscribe \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <jwt-token>" \
    -d '{
    "email": "anothertestuser@example.com"
    }'
    ```
Jwt токен будет ответом на запрос на логин.

#### POST /unsubscribe

    ```bash
    curl -X POST http://localhost:8080/unsubscribe \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <jwt-token>" \
    -d '{
    "email": "jane@example.com"
    }'
    ```

#### GET /users

    ```bash
    curl -X GET http://localhost:8080/users \
    -H "Authorization: Bearer <jwt-token>"
    ```

### Тесты

    ```bash
    go test ./...
    ```