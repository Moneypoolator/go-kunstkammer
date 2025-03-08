# [go-kunstkammer][1]

Это приложение позволяет автоматически создавать задачи в Kaiten на основе файла с описанием задач и конфигурации.

Проект является альтернативной реализацией на Go проекта [ggis-panopticon][2]

## Установка инструментов Go в Ubuntu

Для работы с приложением необходимо установить Go (Golang) на вашу систему.

Обновите пакеты:

```bash
sudo apt update
```

Установите Go:

```bash
sudo apt install golang
```

Проверьте установку:

```bash
go version
```

Настройте переменную окружения GOPATH (если не настроена):

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Данная переменная полезна при установке собранного приложения командрй `go install`

## Запуск приложения из исходного кода

Клонируйте репозиторий:

```bash
git clone https://github.com/Moneypoolator/go-kunstkammer.git
cd go-kunstkammer
```

Установите зависимости:

```bash
go mod tidy
```

Запустите приложение:

```bash
go run cmd/kunstkammer/main.go -tasks tasks.json -config config.json
```

## Сборка исполняемого файла

Перейдите в корневую директорию проекта:

```bash
cd go-kunstkammer
```

Соберите исполняемый файл:

```bash
go build -o kunstkammer
```

Исполняемый файл будет создан в текущей директории с именем `kunstkammer`.

## Установка собранного исполняемого файла в системе

Переместите исполняемый файл в директорию, доступную в PATH:

```bash
sudo mv kunstkammer /usr/local/bin/
```

Проверьте установку:

```bash
kunstkammer --help
```

## Использование приложения

Приложение создает задачи в Kaiten на основе файла с описанием задач и конфигурации.

Пример команды:

```bash
kunstkammer -tasks tasks.json -config config.json
```

Параметры:

- *tasks*: Путь к файлу с описанием задач (например, tasks.json).

- *config*: Путь к файлу конфигурации (например, config.json).

## Создание и формат файла с описанием задач

Файл с описанием задач должен быть в формате JSON. Пример файла tasks.json:

```json
{
  "schedule": {
    "parent": "123",
    "responsible": "user@example.com",
    "tasks": [
      {
        "type": "delivery",
        "size": 8,
        "title": "Task 1"
      },
      {
        "type": "discovery",
        "size": 16,
        "title": "Task 2"
      }
    ]
  }
}
```

Поля:

- *parent*: ID родительской карточки.

- *responsible*: Email ответственного за задачи.

- *tasks*: Список задач.

- *type*: Тип задачи (например, delivery или discovery).

- *size*: Размер задачи в часах.

- *title*: Название задачи. __Номер и код задачи будет сгенерирован автоматически.__

## Создание и настройка конфигурации

Файл конфигурации должен быть в формате JSON. Пример файла config.json:

```json
{
  "token": "your-kaiten-api-token",
  "base_url": "https://kaiten-host-name/api/latest",
  "log_level": "debug"
}
```

Поля:

- *token*: Токен для доступа к API Kaiten.

- *base_url*: Базовый URL API Kaiten.

- *log_level*: Уровень логирования (debug, info, error).

## План улучшений

1. Поддержка переменных окружения: Добавить возможность настройки конфигурации через переменные окружения.

1. Тестирование: Написать unit-тесты для всех модулей.

1. Логирование: Добавить более детальное логирование для отладки.

1. Поддержка других форматов: Добавить поддержку YAML и TOML для конфигурации и файлов задач.

1. Интеграция с CI/CD: Настроить автоматическую сборку и тестирование через GitHub Actions.

1. Интерфейс командной строки: Улучшить CLI с поддержкой подкоманд и справки.

1. Обработка ошибок: Улучшить обработку ошибок и вывод полезных сообщений для пользователя.

1. Поддержка других систем: Добавить возможность интеграции с другими системами управления задачами (например, Jira, Trello).

## Лицензия

Этот проект распространяется под лицензией MIT.

[1]: https://github.com/Moneypoolator/go-kunstkammer "go-kunstkammer на Github"
[2]: https://github.com/Greggot/ggis-panopticon "ggis-panopticon на Github"
