package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type BackupConfig struct {
	BackupFiles   string
	Destination    string
}

func GetHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal( err )
	}
	return dirname
}

func GetConfig() BackupConfig {

	file, err := os.Open(GetHomeDir() + "/.config/backups/" + "config.json")

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	b, err := ioutil.ReadAll(file)

	configJson := string(b)

	var conf BackupConfig
	err = json.Unmarshal([]byte(configJson), &conf)
	if err != nil {
		log.Fatal("Something went wrong while trying to parse the config")
	}

	if len(conf.BackupFiles) <= 0 {
		log.Fatal("Please specify a directory to backup")
	}

	if len(conf.Destination) <= 0 {
		log.Fatal("Please provide a destination")
	}

	return conf
}

func hostName () string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return hn
}

func archiveFileName () string {
	var name string

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()

	y := strconv.Itoa(currentYear)
	m := strconv.Itoa(int(currentMonth))
	d := strconv.Itoa(currentDay)
	h := hostName()

	name += h + "-" + y + "-" + m + "-" + d

	return name
}

func zipFolder(config BackupConfig, name string) {
	out, err := exec.Command("tar", "-zcvf", config.Destination + "/" + name + ".tar.gz", config.BackupFiles).Output()
	if err != nil {
		log.Fatal(err)
	}
	println(string(out))
}

func main() {
	arch := fmt.Sprintf("%s.tgz",  archiveFileName())

	conf := GetConfig()

	println("Backing up " + conf.BackupFiles + " To " + conf.Destination + " as " + arch)

	zipFolder(conf, arch)
}
