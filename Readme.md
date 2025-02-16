# Магазин мерча

## Установка

```
git clone https://github.com/ваш-репозиторий/avitoMerch.git
cd avitoMerch
go mod tidy
```

## Настройка

Приложение использует несколько переменных окружения:

- ```DATABASE_HOST```: Хост базы данных. По умолчанию используется db.  
- ```DATABASE_PORT```: Порт базы данных. По умолчанию используется 5432.  
- ```DATABASE_USER```: Имя пользователя для подключения к базе данных. По умолчанию используется postgres.  
- ```DATABASE_PASSWORD```: Пароль для подключения к базе данных. Установите его на значение, которое вы используете (например, password).  
- ```DATABASE_NAME```: Имя базы данных. По умолчанию используется shop.  
- ```SERVER_PORT```: Порт, на котором будет работать сервер. По умолчанию используется порт 8080.  
- ```JWT_SECRET_KEY```: Секретный ключ для аутентификации JWT. Установите его на значение, которое вы хотите использовать (например, your-secret-key).  
- ```BCRYPT_COST```: Стоимость хеширования для Bcrypt. По умолчанию используется значение 12.  

## Структура проекта

- ```cmd/main.go```: Главная точка входа приложения, настраивает зависимости и запускает HTTP-сервер.  
- ```internal/app```: Содержит основную логику приложения, включая подключение к базе данных, маршрутизацию и middleware.
- ```internal/config```: Обрабатывает загрузку конфигурации приложения из переменных окружения.
- ```internal/models```: Определяет основные модели данных (Пользователь, Товар, Транзакция).
- ```internal/handler```: HTTP обработчики для различных API эндпоинтов (аутентификация, покупка товаров и т.д.).
- ```internal/middleware```: Пользовательское middleware для аутентификации (JWT), логирования и других функциональностей.
- ```internal/repository```: Уровень доступа к данным с интерфейсами и реализациями для взаимодействия с базой данных.
- ```internal/service```: Уровень бизнес-логики, организующий взаимодействие между репозиториями и обработчиками.
- ```tests/```: Содержит различные тестовые наборы, включая юнит-тесты, интеграционные тесты и E2E тесты для обеспечения качества кода и функциональности.

### Контейнер

Для запуска приложения в контейнере необходимо выполнить команду:
- ```sudo docker compose up```  
Для предварительного тестрирования раскомментируйте строку ```CMD ["go", "test", "./tests/..."]``` и выполните команду  
- ```sudo docker compose up```  
Перед запуском каждый из команд не забудьте удалить соответствующий container, image и volume (если они были созданы ранее)