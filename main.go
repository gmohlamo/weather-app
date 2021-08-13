package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/tidwall/gjson"
)

const configDir = "/.config/weather-app/"
const tempPipe = "weather-app-tempc"
const tempFeel = "weather-app-feelc"

func getTempPipePath() string {
	home := os.Getenv("HOME")
	path := home + configDir + tempPipe
	fileInfo, err := os.Stat(path)
	if err == nil {
		if (fileInfo.Mode() & os.ModeNamedPipe) > 0 {
			return path
		} else {
			log.Fatalf("Error: Temp Pipe File \"%s\" exists; however, it is not a named pipe... change source code to accomodate\n", path)
		}
	}
	err = syscall.Mkfifo(path, 0644) //owner can write, but everybody can read
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	return path
}

func getFeelPipePath() string {
	home := os.Getenv("HOME")
	path := home + configDir + tempFeel
	fileInfo, err := os.Stat(path)
	if err == nil {
		if (fileInfo.Mode() & os.ModeNamedPipe) > 0 {
			return path
		} else {
			log.Fatalf("Error: Temp Pipe File \"%s\" exists; however, it is not a named pipe... change source code to accomodate\n", path)
		}
	}
	err = syscall.Mkfifo(path, 0644) //owner can write, but everybody can read
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	return path
}

func getFifo(fileName string) *os.File {
	f, _ := os.OpenFile(fileName, os.O_WRONLY|syscall.O_NONBLOCK, os.ModeNamedPipe)
	return f
}

func writePipes() {
	home := os.Getenv("HOME")
	file, err := ioutil.ReadFile(home + configDir + "api.key")
	if err != nil {
		log.Fatalf("Please ensure that the your API Key is included in a file named \"api.key\" under your $HOME/.config/weather-app directory\nError: %s\n")
	}
	key := strings.Trim(string(file), "\n")
	tempFile := getFifo(getTempPipePath())
	feelFile := getFifo(getFeelPipePath())
	for {
		res, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + key + "&q=auto:ip")
		if err != nil {
			log.Fatalf("Error: %s")
		}
		content, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatalf("Error: %s")
		}
		contentStr := string(content)
		_, err = fmt.Fprintf(tempFile, "%s", gjson.Get(contentStr, "current.temp_c").String())
		if err != nil {
			tempFile.Close()
			tempFile = getFifo(getTempPipePath())
		}
		fmt.Printf("Temp: %s\n", gjson.Get(contentStr, "current.temp_c").String())
		_, err = fmt.Fprintf(feelFile, "%s", gjson.Get(contentStr, "current.feelslike_c").String())
		if err != nil {
			feelFile.Close()
			feelFile = getFifo(getFeelPipePath())
		}
		fmt.Printf("Feel: %s\n", gjson.Get(contentStr, "current.feelslike_c").String())
		time.Sleep(30 * time.Second)
	}
}

func main() {
	fmt.Printf("Starting process")
	writePipes()
}
