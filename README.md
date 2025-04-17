# REST API Server

Простой REST API сервер для управления пользователями с базой данных MongoDB. Реализованы основные CRUD операции.

Особенности:
- Создание, чтение, обновление и удаление пользователей
- Настройка через конфигурационный файл (YAML)
- Поддержка двух типов подключения: TCP-порт или Unix-сокет
- Логирование операций
- Обработка ошибок с middleware
- Интеграция с MongoDB

Технологии:
- Golang
- HTTP Router: [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- Конфигурация: [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- MongoDB Driver: [official mongo-driver](https://go.mongodb.org/mongo-driver)
- Логирование:[logrus](https://github.com/sirupsen/logrus)
