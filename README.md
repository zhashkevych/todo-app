# REST API Для Создания TODO Списков на Go

## <a href="https://www.youtube.com/playlist?list=PLbTTxxr-hMmyFAvyn7DeOgNRN8BQdjFm8">Видеокурс на YouTube</a>

## В курсе разобранны следующие концепции:
- Разработка Веб-Приложений на Go, следуя дизайну REST API.
- Работа с фреймворком <a href="https://github.com/gin-gonic/gin">gin-gonic/gin</a>.
- Подход Чистой Архитектуры в построении структуры приложения. Техника внедрения зависимости.
- Работа с БД Postgres. Запуск из Docker. Генерация файлов миграций. 
- Конфигурация приложения с помощью библиотеки <a href="https://github.com/spf13/viper">spf13/viper</a>. Работа с переменными окружения.
- Работа с БД, используя библиотеку <a href="https://github.com/jmoiron/sqlx">sqlx</a>.
- Регистрация и аутентификация. Работа с JWT. Middleware.
- Написание SQL запросов.
- Graceful Shutdown

### Для запуска БД:

```
docker pull postgres
docker run --name=todo-db -e POSTGRES_PASSWORD='qwerty' -p=5432:5432 -d --rm postgres
```

### Для запуска приложения:

```
go run cmd/main.go
```
