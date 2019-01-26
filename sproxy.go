package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/chinx/sproxy/config"
	"github.com/chinx/sproxy/server"
)

func main() {
	flag.Parse()
	configFile := ""
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	if configFile != "" {
		if exists, err := server.IsFileExists(configFile); !exists || err != nil {
			configFile = ""
		}
	}

	if configFile == "" {
		binDir := path.Dir(os.Args[0])
		if binDir == "" || binDir == "." {
			log.Fatalln("not found config file")
		}
		oldConfig := configFile
		configFile = path.Join(binDir, "config.yaml")
		log.Printf("%s not found, try config file %s\n", oldConfig, configFile)
	}
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Println(err)
		return
	}

	err = server.NewServer(conf.Listener.Socks, conf.Listener.HTTP).
		ListenAndProxy(conf.Server.Addr, conf.Server.Method, conf.Server.Password)
	if err != nil {
		log.Fatal("Run server failed:", err)
	}
}
