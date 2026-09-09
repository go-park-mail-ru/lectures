# Шпаргалка: команды Go

На лекции почти всегда достаточно `go run`. `go test` понадобится уже в первом ДЗ.

## Проверить установку

```bash
go version
go env GOPATH GOROOT
```

`GOROOT` — сам компилятор, `GOPATH` — кэш и модули. В повседневной работе почти не трогаем.

## Запустить код

```bash
go run 01_vars_1.go
go run .
```

Аргументы после имени файла — это аргументы вашей программы:

```bash
go run uniq.go input.txt
go run uniq.go -c input.txt
cat input.txt | go run uniq.go
```

## Собрать бинарник

```bash
go build
go build -o uniq
./uniq input.txt
```

Кросс-компиляция (собрать `.exe` с мака):

```bash
GOOS=windows GOARCH=amd64 go build -o uniq.exe
```

## Модуль и зависимости

В корне проекта, один раз:

```bash
go mod init github.com/USERNAME/homework1
go mod tidy
```

`tidy` подтянет нужные пакеты и выкинет лишние.

## Форматирование

```bash
gofmt -w .
go fmt ./...
```

Стиль в Go не обсуждают. Лучше включить format on save в IDE.

## Тесты

```bash
go test
go test -v
go test -cover
go test ./...
go test -v -run TestOK
```

`-v` — какие тесты прошли, `-cover` — покрытие, `./...` — все пакеты ниже текущей папки.

## Документация

```bash
go doc fmt.Println
go doc flag
go help
go help run
go help test
```

Сайт пакетов: https://pkg.go.dev
