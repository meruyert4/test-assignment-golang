# Go Package Manager (Test Assignment)

## Описание

Это CLI-пакетный менеджер на Go, который умеет:
- Упаковывать файлы в архив по маскам, указанным в конфиге (JSON)
- Загружать архивы на сервер по SSH
- Скачивать архивы по SSH и распаковывать их

## Используемые технологии
- Go (стандартная библиотека)
- [github.com/pkg/sftp](https://github.com/pkg/sftp) — работа с SFTP по SSH
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) — SSH-клиент
- [github.com/joho/godotenv](https://github.com/joho/godotenv) — загрузка переменных окружения из .env

## Форматы конфигов

### Пример файла для упаковки (packet.json)
```json
{
  "name": "packet-1",
  "ver": "1.10",
  "targets": [
    "./archive_this1/*.txt",
    { "path": "./archive_this2/*", "exclude": ["*.tmp"] }
  ],
  "packets": [
    { "name": "packet-3", "ver": "<=2.0" }
  ]
}
```

- `targets` — массив строк (маска) или объектов `{path, exclude}`
- `packets` — зависимости (опционально)

### Пример файла для распаковки (packages.json)
```json
{
  "packages": [
    { "name": "packet-1", "ver": ">=1.10" },
    { "name": "packet-2" },
    { "name": "packet-3", "ver": "<=1.10" }
  ]
}
```
- `ver` может быть с условием: `>=`, `<=`, `>`, `<`, `=`

## Переменные окружения
Создайте файл `.env` в корне проекта:
```
SSH_HOST=localhost
SSH_PORT=22
SSH_USER=ваш_логин_на_сервере
SSH_PASSWORD=ваш_пароль
```

## Как пользоваться

```bash
go build -o pm cmd/main.go
```

### 1. Упаковка и загрузка архива
```sh
./pm create ./packet.json
```
- Собирает файлы по маскам, архивирует, отправляет архив на сервер по SSH.

### 2. Скачивание и распаковка архивов
```sh
./pm update ./packages.json
```
- Скачивает архивы с сервера по условиям версий, распаковывает их в текущую папку.# test-assignment-golang
