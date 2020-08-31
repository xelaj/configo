package simplest

import (
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/xelaj/configo"
	"github.com/xelaj/go-dry"
)

type AppConfig struct {
	FilesPath string `param:"template_path" validate:"required"`
	Styles    string `param:"styles"        validate:"required"`
}

var (
	Config *AppConfig = &AppConfig{}
)

func TryIt() {
	godotenv.Load("./examples/simplest/session.env")
	err := configo.InitConfigWithExplicitConfigPaths("simplest", "./examples/simplest/system", "./examples/simplest/user", Config)
	dry.PanicIfErr(err)
	pp.Println(Config)
}
