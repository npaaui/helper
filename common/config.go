package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig(configFolder, configType string) {
	configFolder = strings.TrimRight(configFolder, string(os.PathSeparator)) + string(os.PathSeparator)
	configName := "config." + configType
	configPath := ""

	envNames := [2]string{"local", "release"}
	for _, item := range envNames {
		if _, err := os.OpenFile(configFolder+item+string(os.PathSeparator)+configName, os.O_RDONLY, 0); err == nil {
			configPath = configFolder + item + string(os.PathSeparator)
			break
		}
	}

	if configPath == "" {
		panic("config file empty")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
