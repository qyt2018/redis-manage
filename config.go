package main

import(
	"os"
	"path/filepath"
	"encoding/json"
	"fmt"
)

const accountFile string = "account.conf"

func saveConfig(mt *server){
	config := make(map[string]server)
	wd, _ := os.Getwd()
	configFile := filepath.Join(wd, accountFile)
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			fs, _ := os.OpenFile(configFile, os.O_CREATE|os.O_RDWR, 0755)
			defer fs.Close()
			config[mt.Name] = *mt
			jconfig, _ := json.Marshal(config)
			fmt.Println(string(jconfig))
			fs.Write(jconfig)
			return
		}
	}
	//读
	rfs, _ := os.OpenFile(configFile, os.O_RDONLY, 0755)
	configData := make([]byte, 1024000)
	n, _ := rfs.Read(configData)
	json.Unmarshal(configData[:n], &config)
	rfs.Close()
	//值处理
	config[mt.Host] = *mt
	jconfig, _ := json.Marshal(config)
	fmt.Println(string(jconfig))
	//写
	wfs, _ := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC, 0755)
	defer wfs.Close()
	wfs.Write(jconfig)
	return
}

func getConfig() string {
	wd, _ := os.Getwd()
	configFile := filepath.Join(wd, accountFile)
	fs, _ := os.OpenFile(configFile, os.O_RDONLY, 0755)
	defer fs.Close()
	configData := make([]byte, 1024000)
	n, _ := fs.Read(configData)
	conf := string(configData[:n])
	if conf == "" {
		return "{}"
	}
	return conf
}
