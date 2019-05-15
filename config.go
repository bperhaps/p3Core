package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func Testa() {
	//	viper.SetDefault("test", "ms")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	//test := viper.Get("compileOptions")
	/*	for _, d := range test.(type) {
		fmt.Println(d)
	}*/
	viper.Set("test", "sm")
	fmt.Println(viper.Get("test"))
}
