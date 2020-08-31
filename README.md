# {{ .Project.Name }}

![help wanted](https://img.shields.io/badge/-help%20wanted-success)
[![godoc reference](https://godoc.org/{{ .PackageUrl }}?status.svg)](https://godoc.org/{{ .PackageUrl }})
[![Go Report Card](https://goreportcard.com/badge/{{ .PackageUrl }})](https://goreportcard.com/report/{{ .PackageUrl }})
[![license MIT](https://img.shields.io/badge/license-MIT-green)](https://{{ .PackageUrl }}/blob/master/README.md)
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

**english** [русский](https://{{ .PackageUrl }}/blob/master/doc/ru_RU/README.md)

{{ .Title.Text }}

<p align="center">
<img src="{{ .Title.ImageUrl }}"/>
</p>

## Getting started

TODO

## How to use

![preview]({{ .PreviewUrl }})

**Code examples are [here](https://{{ .PackageUrl }}/blob/master/examples)**

### Simple How-To

{{ .AdditionalHowto }}

{{ .SimpleFAQ }}

## Contributing

Please read [contributing guide](https://{{ .PackageUrl }}/blob/master/doc/en_US/CONTRIBUTING.md) if you want to help. And the help is very necessary!

## TODO

{{ range $item := .TODO }}* {{ $item }}
{{ end }}
## Authors

{{ range $author := .Authors }}* **{{ $author.Name }}** — [{{ $author.Nick }}](https://github.com/{{ $author.Nick }})
{{ end }}
## License

This project is licensed under the MIT License - see the [LICENSE](https://{{ .PackageUrl }}/blob/master/doc/en_US/LICENSE.md) file for details
