package fastconfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var _viper *viper.Viper

func getExeDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	fileDir := filepath.Dir(executablePath)
	return fileDir
}

func initConfig() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	exeDir := getExeDir()
	log.Println(fmt.Sprintf("exe path is %s", exeDir))
	fi, err := os.Stat(exeDir + "/config/application.yaml")
	workPath := ""
	if err == nil && !fi.IsDir() {
		log.Println("use exe path for config")
		workPath = exeDir
	} else {
		log.Println("use pwd for config")
		workPath, err = os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println(fmt.Sprintf("work path is %s", workPath))
	_viper = viper.New()
	err = _viper.BindEnv(`APP_ENV`)
	if err != nil {
		log.Println(err)
	}
	// 通过环境变量获取当前运行的环境 如 test uat prd, 后续读取相应的配置文件
	appEnv := _viper.GetString(`APP_ENV`)
	log.Println(fmt.Sprintf("app env is %s", appEnv))
	_viper.SetConfigName("application")
	_viper.SetConfigType("yaml")
	_viper.AddConfigPath(workPath + "/config/")
	var readErr error
	// 向工作目录往上读取 20 层，防止单元测试无法读取文件
	for i := 0; i < 20; i++ {
		if readErr = _viper.ReadInConfig(); readErr == nil {
			break
		} else {
			workPath = filepath.Dir(workPath)
			_viper.AddConfigPath(workPath + "/config/")
		}
	}
	if appEnv != "" {
		_viper.SetConfigName("application_" + appEnv)
		readErr = _viper.MergeInConfig()
	}
	if readErr != nil {
		log.Println(readErr)
		log.Println("read config file err, use default config")
	}
	_viper.AutomaticEnv()
}

func init() {
	initConfig()
}

// centerLineToCamelCase 中划线转驼峰
func centerLineToCamelCase(s string) string {
	s = strings.Replace(s, "-", " ", -1)
	s = cases.Title(language.Und).String(s)
	s = strings.Replace(s, " ", "", -1)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// camelCaseToCenterLine 驼峰转中划线
func camelCaseToCenterLine(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '-')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}

// convertKey 将中划线转驼峰或将驼峰转中划线
func convertKey(s string) string {
	if strings.Contains(s, "-") {
		return centerLineToCamelCase(s)
	}
	return camelCaseToCenterLine(s)
}

func GetValue(key string, defaultVal any) any {
	if _viper.InConfig(key) {
		return _viper.Get(key)
	}
	cKey := convertKey(key)
	if _viper.InConfig(cKey) {
		return _viper.Get(cKey)
	}
	return defaultVal
}

func GetString(key string, defaultVal string) string {
	return cast.ToString(GetValue(key, defaultVal))
}

func GetInt(key string, defaultVal int) int {
	return cast.ToInt(GetValue(key, defaultVal))
}

func GetBool(key string, defaultVal bool) bool {
	return cast.ToBool(GetValue(key, defaultVal))
}
