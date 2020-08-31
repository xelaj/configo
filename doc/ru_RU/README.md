# Configo

![help wanted](https://img.shields.io/badge/-help%20wanted-success)
[![godoc reference](https://godoc.org/github.com/xelaj/configo?status.svg)](https://godoc.org/github.com/xelaj/configo)
[![Go Report Card](https://goreportcard.com/badge/github.com/xelaj/configo)](https://goreportcard.com/report/github.com/xelaj/configo)
[![license MIT](https://img.shields.io/badge/license-MIT-green)](https://github.com/xelaj/configo/blob/master/README.md)
[![chat telegram](https://img.shields.io/badge/chat-telegram-0088cc)](https://bit.ly/2xlsVsQ)
![version v1.0.0](https://img.shields.io/badge/version-v0.1.0-red)
![unstable](https://img.shields.io/badge/stability-unstable-yellow)
<!--
code quality
golangci
contributors
go version
gitlab pipelines
-->


[english](https://github.com/xelaj/configo/blob/master/doc/en_US/README.md) **русский**

По факту, Самый удобный способ сконфигурировать ваше приложение.

<p align="center">
<img src="{{ .Title.ImageUrl }}"/>
</p>

## Как установить

В первую очередь, необходимо установить сам пакет

``` sh
go get -v github.com/xelaj/configo
```

Далее вам нужно сформировать структуру конфигурации так, как вам нравится, и всего лишь проинициализировать конфигурацию в объект структуры, например

``` go
package main

import (
    "fmt"
    "github.com/xelaj/configo"
)

type AppConfig struct {
    Host string
    Port string
    Users []string
}

func main() {
    config := new(AppConfig)
    err := configo.Init("myapp")
    if err != nil {
        panic(err)
    }

    fmt.Println(config)
}
```

Все! Это самое базовое решение, которое вы можете сделать, у configo есть так же к

### Тег param

``` hcl
host = dev.xelaj.org
port = 65534
users = [
    "Jackie"
    "Alice"
    "Michael"
]
```

``` go
type AppConfig struct {
    AppHost    string   `param:"host"`
    ServerPort uint16   `param:"port"`
    Users      []string `param:"users"`
}

func main() {
    fmt.Println(config)
}

// &AppConfig{
//     AppHost: "dev.xelaj.org",
//     ServerPort: 0x10,
//     Users: []string{
//         "Jackie",
//         "Alice",
//         "Michael",
//     },
// }
```

### Тег default

``` go
type AppConfig struct {
    AppHost    string   `param:"host"`
    ServerPort uint16   `param:"port"     default:"80"`
    Users      []string `param:"users"`
    TLSKeys    []string `param:"tls_keys" default:"[\"a\", \"b\"]"` // even json!
}

func main() {
    fmt.Println(config)
}

// &AppConfig{
//     AppHost: "dev.xelaj.org",
//     ServerPort: 0x10,
//     Users: []string{
//         "Jackie",
//         "Alice",
//         "Michael",
//     },
//     TLSKeys: [
//         "a",
//         "b",
//     ],
// }
```

### Тег validate

позволяет короче проверить значение перед тем как продолжить


### Вывод ошибок

выводит список либо переменных, которые не могут впихнуться в конфиг. ну или валидацию не проходят

### Usage

выводит список параметров которые могут быть использованы в конфигах

выводит или env переменными или списком из jsonpath

## Как использовать

![preview]({{ .PreviewUrl }})

**Примеры кода [здесь](https://github.com/xelaj/configo/blob/master/examples)**

### Simple How-To

{{ .AdditionalHowto }}

{{ .SimpleFAQ }}

## Вклад в проект

пожалуйста, прочитайте [информацию о помощи]https://github.com/xelaj/configo/blob/master/doc/ru_RU/CONTRIBUTING.md), если хотите помочь. А помощь очень нужна!

## TODO

{{ range $item := .TODO }}* {{ $item }}
{{ end }}
## Авторы

{{ range $author := .Authors }}* **{{ $author.Name }}** — [{{ $author.Nick }}](https://github.com/{{ $author.Nick }})
{{ end }}
## Лицензия

This project is licensed under the MIT License - see the [LICENSE](https://github.com/xelaj/configo/blob/master/doc/ru_RU/LICENSE.md) file for details
