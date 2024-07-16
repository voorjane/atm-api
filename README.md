## Тестовое задание

### Задача

[Click](https://docs.google.com/document/d/117-46922KU_HaKepb3v1zirv9S2m7Ue69jgXfoRGjTg/edit?pli=1#heading=h.pf7how9ocifx)

### Как запускать

- Создать файл .env
Необходимые параметры: \
`DATABASE_HOST` \
`DATABASE_PORT` \
`POSTGRES_USER` \
`POSTGRES_PASSWORD` \
`POSTGRES_DB`


- `make`

- После запуска контейнеров api готово к работе
- Порт 8888

### Запросы

`POST`
- /accounts/{id}/deposit - пополнение баланса
- /accounts/{id}/withdraw - снятие средств
- /accounts - создание нового аккаунта
Все POST запросы принимают JSON с балансом: \
`{ "balance": 123 }`

`GET`
- /accounts/balance - проверка баланса
