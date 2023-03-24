# Менеджер паролей GophKeeper

## Описание проекта
GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

Проект разработан на базе инструментария Fyne https://fyne.io и представляет собой полноценное десктопное кроссплатформенное приложение

<p  style="text-align: center"><img src="img/1.jpg"></p>
<p  style="text-align: center"><img src="img/2.jpg"></p>
<p  style="text-align: center"><img src="img/3.jpg"></p>

## Сервер

### БД

Сервер использует для хранения базу postgresql

Настройки для docker контейнера расположены в папке `.docker`

Миграцию структуры таблиц можно запустить командой через Makefile
`make migrate-up`

### SSL сертификаты

Сертификаты можно сгенерировать командой через Makefile `make cert`

### Запуск сервера

Запустить можно командой

`
go run cmd/server/main.go -a localhost:8080 -d postgres://root:12345@localhost:5432/gophkeeper -j jwt_secret_key_here
`

## Клиент - приложение

Запустить можно командой go run `go run cmd/client/main.go`
Дополнительные параметры запуска можно посмотреть в Makefile

## Endpoints

- `POST /sign-up`
    - Обработчик регистрации пользователя
    - Запрос: `{"login":testuser","password":"testpassword"}`
    - Ответ: `{"access_token": "accesstoken","refresh_token": "refreshtoken"}`

- `POST /sign-in`
    - Обработчик авторизации пользователя
    - Запрос: `{"login":testuser","password":"testpassword"}`
    - Ответ: `{"access_token": "accesstoken","refresh_token": "refreshtoken"}`

- `POST /refresh-token`
    - Обработчик обновление токенов пользователя
- 
- `POST /sign-key`
    - Обработчик получение ключа подписи для данных

### Карты

Требуется авторизация   `Authorization: Bearer access_token`

- `POST /store/card`
    - Обработчик добавления данных карт
- `DELETE /store/card`
  - Обработчик удаления данных карты
- `GET /store/card`
  - Обработчик просмотра данных карты
- `GET /store/card/list`
  - Обработчик просмотра списка карт

### Логин / Пароль

Требуется авторизация `Authorization: Bearer access_token`

- `POST /store/cred`
  - Обработчик просмотра записи логин/пароль
- `DELETE /store/cred`
  - Обработчик удаления записи логин/пароль
- `GET /store/cred`
  - Обработчик просмотра записи логин/пароль
- `GET /store/cred/list`
  - Обработчик просмотра списка логин/пароль
### Текстовые данные

Требуется авторизация `Authorization: Bearer access_token`

- `POST /store/text`
  - Обработчик добавления текстовых данных
- `DELETE /store/text`
  - Обработчик удаления текстовых данных
- `GET /store/text`
  - Обработчик просмотра текстовых данных
- `GET /store/text/list`
  - Обработчик просмотра списка текстовых данных

### Файлы

Требуется авторизация `Authorization: Bearer access_token`

- `POST /store/file`
    - Обработчик добавления файла
- `DELETE /store/file`
    - Обработчик удаления файла
- `GET /store/file`
    - Обработчик просмотра данных файла
- `GET /store/file/list`
    - Обработчик просмотра списка файлов