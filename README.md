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

Лучше ставить из snap репозитория, там будет самая свежая версия:

```bash
sudo snap install go
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
go build -o kunstkammer ./cmd/kunstkammer/*.go 
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
        "type": "delivery",
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

- *type*: Тип задачи (например, delivery). Задачи типа __discovery__ после исправления методики не доступны для создания.

- *size*: Размер задачи в часах.

- *title*: Название задачи. __Номер и код задачи будет сгенерирован автоматически.__

## Создание и настройка конфигурации

Файл конфигурации должен быть в формате JSON. Пример файла config.json:

```json
{
  "token": "your-kaiten-api-token",
  "base_url": "https://kaiten-host-name/api/latest",
  "log_level": "debug",
  "board_id": 192,
  "column_id": 776,
  "lane_id": 1275,
  "tags": ["ГГИС", "C++"]
}
```

Поля:

- *token*: Токен для доступа к API Kaiten.

- *base_url*: Базовый URL API Kaiten.

- *log_level*: Уровень логирования (debug, info, error).

- *"board_id"*: Номер доски в Kaiten, по-умолчанию 192.

- *"column_id"*: Номер колонки в Kaiten, по-умолчанию 776.

- *"lane_id"*: Номер дорожки в Kaiten, по-умолчанию 1275.

- *"tags"*: Список тэгов, которыми помечается карточка, по-умолчанию "ГГИС", "C++".

## Запуск модульных тестов

Для запуска модульных тестов в проекте используется стандартный инструмент go test. Ниже приведены инструкции по запуску тестов для каждого модуля.

### Запуск всех тестов

Чтобы запустить все модульные тесты в проекте, выполните команду в корневой директории проекта:

```bash
go test -v ./...
```

Эта команда запустит тесты для всех модулей, включая `internal/api`, `internal/models`, и другие.

### Запуск тестов для конкретного модуля

Если вы хотите запустить тесты только для определенного модуля, укажите путь к этому модулю.

#### Тесты для модуля `internal/api`

```bash
go test -v ./internal/api
```

#### Тесты для модуля `internal/models`

```bash
go test -v ./internal/models
```

#### Тесты для модуля internal/utils

```bash
go test -v ./internal/utils
```

### Запуск тестов с покрытием

Чтобы проверить покрытие кода тестами, используйте флаг -cover:

```bash
go test -v -cover ./...
```

Эта команда покажет процент покрытия кода тестами для каждого модуля.

### Генерация отчета о покрытии

Для более детального анализа покрытия кода можно сгенерировать HTML-отчет:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

После выполнения этих команд откроется файл coverage.html в браузере, где вы сможете увидеть, какие строки кода покрыты тестами, а какие нет.

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
