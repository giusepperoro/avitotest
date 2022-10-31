# avitotest1022
___
# Микросервис для работы с балансом пользователей
___
Этот сервис API, состоящий из базы данных PostgreSQL, в которую записываются данные работы балансом и сервис golang API, который обрабатывает POST и GET запросы.

# Описание API
### Query: Header POST: ```URL/create ```
#### Body:
```
{
"initial_balance": 10000
}
```
#### Response on success:
```
{
    "client_id": 1234
}
```
#### Body:
```
{
    "initial_balance": -10000
}
```
#### Response on error: 
```
{
    "error": "unable to create user"
}
```
___

### Query: Header GET: ```URL/balance ```
#### Body:
```
{
    "client_id": 3480
}
```
#### Response on success:
```
{
    "balance": 10000
}
```
#### Response on error:
```
{
    "balance": 0,
    "error": "unable to create user"
}
```
___

### Query: Header POST: ```URL/refill ```
#### Body:
```
{
    "client_id": 3480,
    "amount": 3000
}
```
#### Response on success:
```
{
    "approved": true
}
```
#### Response on error:
```
{
    "approved": false,
    "error": "unable to refill balance"
}
```
___
### Query: Header POST: ```URL/withdrawal ```
#### Body:
```
{
    "client_id": 3480,
    "service_id": 256,
    "order_id": 500,
    "amount": 4500
}
```
#### Response on success:
```
{
    "success": true
}
```
#### Response on error:
```
{
    "success": false,
    "error": "unable to create withdrawal"
}
```
___

### Query: Header POST: ```URL/processWithdrawal ```
#### Body:
```
{
    "client_id": 3480,
    "service_id": 256,
    "order_id": 500,
    "amount": 4500
}
```
#### Response on success:
```
{
    "success": true
}
```
#### Response on error:
```
{
    "success": false,
    "error": "unable to process withdrawal"
}
```
___

### Query: Header POST: ```URL/cancelWithdrawal ```
#### Body:
```
{
    "client_id": 3480,
    "service_id": 256,
    "order_id": 500,
    "amount": 4500
}
```
#### Response on success:
```
{
    "success": true
}
```
#### Response on error:
```
{
    "success": false,
    "error": "unable to cansel withdrawal"
}
```
___
# Запуск сервиса локально

#### 1. git clone https://github.com/giusepperoro/avitotest.git

#### 2. Используйте docker-compose для создания контейнеров с помощью PostgreSQL и сервиса API
```
docker compose up
```

#### 3. Откройте браузер или приложение для тестирования API, введите localhost:80/create и сделайте несколько запросов!

___


