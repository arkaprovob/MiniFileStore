package pkg

import (
	"encoding/json"
	"os"
)

type Config struct {
	FileStore   string `json:"file_store"`
	RecordStore string `json:"record_store"`
}

func GetConfig() (Config, error) {
	var config Config

	// Open the config file
	file, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Decode the config file into the Config struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
