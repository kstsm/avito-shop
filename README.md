# Avito Shop

## Установка и запуск проекта

### Клонирование репозитория
Для начала работы склонируйте репозиторий в удобную Вам директорию:
```bash
git clone https://github.com/kstsm/avito-shop
```
### Настройка переменных окружения
Создайте `.env` файл, скопировав в него значения из `.env.example`, и укажите необходимые параметры.
### Запуск проекта
Выполните команду:
```bash
docker-compose up -d --build
```
После запуска сервер будет доступен по адресу: http://localhost:8080

## Тестирование

### Юнит-тесты
Проект покрыт тестами на 80.5%:
```bash
go test -v ./internal/handler -cover
```

![image](https://github.com/user-attachments/assets/9f473111-57c2-46bb-9295-e5c83f16e41c)

### Интеграционные тесты
Интеграционные тесты находятся в отдельной директории `internal/tests`.<br> Реализованы оба сценария из задания:
покупка мерча, передача монет

## Линтер
Описан конфигурационный линтер `.golangci.yam`

## Проблемы
База данных, предназначенная для тестов, по какой-то причине не видела мой конфиг. В результате пришлось реализовать подключение через fmt.Sprintf как временное решение

## Проведенно нагрузочное тестирование

![auth](https://github.com/user-attachments/assets/71c3bb4d-50bc-4d15-8ed4-395310e4a608)
![info](https://github.com/user-attachments/assets/8d248ff9-fb38-4644-868a-1fa2fa96042e)
![coins](https://github.com/user-attachments/assets/f3147110-93f5-49e7-85dd-b9428763cda6)
![item](https://github.com/user-attachments/assets/d9418abc-c5ce-4cd6-93db-f0063fa4adff)



