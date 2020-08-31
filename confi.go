package configo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/fatih/structtag"
	"github.com/go-playground/validator"
	"github.com/iancoleman/strcase"
	"github.com/imdario/mergo"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/ungerik/go-dry"
	"github.com/xelaj/errs"
	"github.com/xelaj/v"
)

var (
	ConfigScheme interface{}
	commonEnvs   = []string{
		"USER",
		"PATH",
		"LANGUAGE",
	}
)

func Init(appName string) (interface{}, error) {
	if appName == "" {
		appName = v.AppName
	}

	err := ParseConfig(appName, ConfigScheme)
	if err != nil {
		return nil, errors.Wrap(err, "parsing config")
	}

	return ConfigScheme, nil
}

func InitConfig(appName, userRunned string, into interface{}) error {
	if appName == "" {
		appName = v.AppName
	}

	u, _ := user.Lookup(userRunned)
	return initConfig(appName,
		filepath.Join("/etc", appName),
		filepath.Join(u.HomeDir, ".local", "etc", appName),
		into,
	)
}

// DEPRECATED: No, seriously, this is just for examples and testing. Use only InitConfig, it do all job for you
func InitConfigWithExplicitConfigPaths(appName, globalPath, localPath string, into interface{}) error {
	return initConfig(appName, globalPath, localPath, into)
}

func initConfig(appName, globalPath, localPath string, into interface{}) error {
	typ := reflect.TypeOf(into)
	if typ.Kind() != reflect.Ptr {
		panic("not a pointer")
	}

	global := reflect.New(typ.Elem()).Interface()
	dry.PanicIfErr(ParseDir(globalPath, global))
	personal := reflect.New(typ.Elem()).Interface()
	dry.PanicIfErr(ParseDir(localPath, personal))
	session := reflect.New(typ.Elem()).Interface()
	dry.PanicIfErr(ParseEnvFile("./configs/session.env", "simpleapp", session))
	dry.PanicIfErr(mergo.Merge(into, global, mergo.WithOverride))
	dry.PanicIfErr(mergo.Merge(into, personal, mergo.WithOverride))
	dry.PanicIfErr(mergo.Merge(into, session, mergo.WithOverride))
	return nil
}

func ParseConfig(appName string, cfg interface{}) error {
	err := envconfig.Process(appName, cfg)
	if err != nil {
		return errors.Wrap(err, "processing env")
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		splitted := strings.Split(err.Error(), "\n")
		return errs.MultipleAsString(splitted...)
	}

	return nil
}

func ParseEnvFile(path, prefix string, into interface{}) error {
	prefix = strcase.ToScreamingSnake(prefix) + "_"

	data, err := ioutil.ReadFile(path)
	dry.PanicIfErr(err)

	envs, err := godotenv.Unmarshal(string(data))
	dry.PanicIfErr(err)

	ival := reflect.ValueOf(into)
	ityp := ival.Type()
	if ityp.Kind() != reflect.Ptr {
		panic("not a pointer")
	}
	if ityp.Elem().Kind() != reflect.Struct {
		panic("not a struct")
	}
	ival = ival.Elem()
	ityp = ityp.Elem()

ForEachField:
	for i := 0; i < ityp.NumField(); i++ {
		fval := ival.Field(i)
		ftyp := ityp.Field(i).Type

		tags, err := structtag.Parse(string(ityp.Field(i).Tag))
		dry.PanicIfErr(err)
		tag, err := tags.Get("param")
		dry.PanicIfErr(err)
		name := strcase.ToScreamingSnake(tag.Name)

		possibleParams := make(map[string]string)
		for k, v := range envs {
			if strings.HasPrefix(k, prefix+name) {
				possibleParams[k] = v
			}
		}

		switch ftyp.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8,
			reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32,
			reflect.Float64, reflect.Complex64,
			reflect.Complex128, reflect.String:
			exactKey := prefix + name

			CastValue(reflect.ValueOf(possibleParams[exactKey]), fval)
			continue ForEachField
		}

		panic("don't understand how to parse this thing!")
	}

	return nil
}

func ParseDir(path string, into interface{}) error {
	ival := reflect.ValueOf(into)
	ityp := ival.Type()
	if ityp.Kind() != reflect.Ptr {
		panic("not a pointer")
	}
	if ityp.Elem().Kind() != reflect.Struct {
		panic("not a struct")
	}
	ival = ival.Elem()
	ityp = ityp.Elem()

	path, err := filepath.Abs(path)
	dry.PanicIfErr(err)

	stat, err := os.Stat(path)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			return errs.NotFound("directory", path)
		default:
			panic(err)
		}
	}

	if stat.Mode()&os.ModeSymlink > 0 {
		return errors.New("doesn't working with symlinks")
	}

	files, err := dry.ListDirFiles(path)
	dry.PanicIfErr(err)

	for _, file := range files {
		_, ext := dry.PathSplitExt(file)
		switch ext {
		case "json":
			in := map[string]interface{}{}
			data, err := ioutil.ReadFile(filepath.Join(path, file))
			dry.PanicIfErr(err)

			err = json.Unmarshal(data, &in)
			dry.PanicIfErr(err)

			err = mergo.Map(into, in)
			dry.PanicIfErr(err)
		case "go":
			continue
		default:
			panic("invalid extension: " + ext)
		}
	}

	dirs, err := dry.ListDirDirectories(path)
	dry.PanicIfErr(err)

	for _, param := range dirs {
		dstElement := ival.FieldByName(param)
		for i := 0; i < ityp.NumField(); i++ {
			if ityp.Field(i).Tag.Get("param") == param {
				dstElement = ival.Field(i)
			}
		}

		zeroValue := reflect.Value{}
		if dstElement == zeroValue {
			panic("unknown field: " + param)
		}
		if dstElement.IsNil() {
			switch dstElement.Type().Kind() {
			case reflect.Ptr:
				dstElement.Set(reflect.New(dstElement.Type()).Elem())

			case reflect.Map:
				dstElement.Set(reflect.MakeMap(dstElement.Type()))

			case reflect.Slice:
				dstElement.Set(reflect.MakeSlice(dstElement.Type(), 0, 0))
			}
		}
		err := parseDir(filepath.Join(path, param), dstElement.Addr().Interface())
		dry.PanicIfErr(err)
	}

	defaults.MustSet(into)

	pp.Println(into)
	return nil
}

func parseDir(path string, into interface{}) error {
	if into == nil {
		panic("into is nil")
	}

	ival := reflect.ValueOf(into)
	ityp := ival.Type()
	if ityp.Kind() != reflect.Ptr {
		panic("not a pointer")
	}
	ival = ival.Elem()
	ityp = ityp.Elem()

	switch ityp.Kind() {
	case reflect.Slice:
		files, err := dry.ListDirFiles(path)
		dry.PanicIfErr(err)
		for _, file := range files {
			childType := ityp.Elem()
			if childType.Kind() == reflect.Ptr {
				childType = childType.Elem()
			}
			item := reflect.New(childType).Interface()

			in := map[string]interface{}{}
			data, err := ioutil.ReadFile(filepath.Join(path, file))
			dry.PanicIfErr(err)

			_, ext := dry.PathSplitExt(file)
			switch ext {
			case "json":
				err = json.Unmarshal(data, &in)
				dry.PanicIfErr(err)

			case "go":
				continue
			default:
				panic("invalid extension: " + ext)
			}

			decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				TagName:          "param",
				WeaklyTypedInput: true,
				Result:           item,
			})
			err = decoder.Decode(in)
			dry.PanicIfErr(err)

			itemValue := reflect.ValueOf(item)
			if ityp.Elem().Kind() != reflect.Ptr && itemValue.Type().Kind() == reflect.Ptr {
				itemValue = itemValue.Elem()
			}
			ival.Set(reflect.Append(ival, itemValue))
		}
	case reflect.Map:
		files, err := dry.ListDirFiles(path)
		dry.PanicIfErr(err)
		switch ityp.Key().Kind() {
		case reflect.Bool, reflect.Int, reflect.Int8,
			reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32,
			reflect.Float64, reflect.Complex64,
			reflect.Complex128, reflect.String:
		default:
			panic("supports only string any number or boolean")
		}

		for _, file := range files {
			key, ext := dry.PathSplitExt(file)
			switch ext {
			case "json":
				childType := ityp.Elem()
				if childType.Kind() == reflect.Ptr {
					childType = childType.Elem()
				}
				item := reflect.New(childType).Interface()

				in := map[string]interface{}{}
				data, err := ioutil.ReadFile(filepath.Join(path, file))
				dry.PanicIfErr(err)

				err = json.Unmarshal(data, &in)
				dry.PanicIfErr(err)

				decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
					TagName:          "param",
					WeaklyTypedInput: true,
					Result:           item,
				})
				err = decoder.Decode(in)
				dry.PanicIfErr(err)

				itemValue := reflect.ValueOf(item)
				if ityp.Elem().Kind() != reflect.Ptr && itemValue.Type().Kind() == reflect.Ptr {
					itemValue = itemValue.Elem()
				}

				keyValue := reflect.New(ityp.Key()).Elem()
				CastValue(reflect.ValueOf(key), keyValue)
				pp.Println(keyValue.Interface(), itemValue.Interface())
				pp.Println(ival.Interface())
				ival.SetMapIndex(keyValue, itemValue)

			case "go":
				continue
			default:
				panic("invalid extension: " + ext)
			}
		}
	default:
		panic(ityp.String())
	}

	return nil
}

func CastValue(src, dst reflect.Value) {
	in := src.Interface()
	var out interface{}
	switch dst.Type().Kind() {
	case reflect.Bool:
		out = cast.ToBool(in)
	case reflect.Int:
		out = cast.ToInt(in)
	case reflect.Int8:
		out = cast.ToInt8(in)
	case reflect.Int16:
		out = cast.ToInt16(in)
	case reflect.Int32:
		out = cast.ToInt32(in)
	case reflect.Int64:
		out = cast.ToInt64(in)
	case reflect.Uint:
		out = cast.ToUint(in)
	case reflect.Uint8:
		out = cast.ToUint8(in)
	case reflect.Uint16:
		out = cast.ToUint16(in)
	case reflect.Uint32:
		out = cast.ToUint32(in)
	case reflect.Uint64:
		out = cast.ToUint64(in)
	case reflect.Float32:
		out = cast.ToFloat32(in)
	case reflect.Float64:
		out = cast.ToFloat64(in)
	case reflect.String:
		out = cast.ToString(in)
	default:
		panic("unsupported type " + dst.Type().String())
	}

	dst.Set(reflect.ValueOf(out))
}
