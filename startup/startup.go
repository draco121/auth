package startup

import (
	"encoding/json"
	"io"
	"os"
)

type Configuration struct {
	Port       string `json:"port"`
	Mongodburi string `json:"mongodburi"`
	SecretKey  string `json:"secretkey"`
}

var Config *Configuration = nil

func Initialize() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	} else {
		bytevalue, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		} else {
			json.Unmarshal(bytevalue, &Config)
		}
	}
}
