package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/kr/pretty"
	"github.com/tidwall/gjson"
)

const configDir = "/.config/weather-app/"
const tempPipe = "weather-app-tempc"
const tempFeel = "weather-app-feelc"

func getTempPipePath() string {
	home := os.Getenv("HOME")
	path = home + configDir + tempPipe
	fileInfo, err := os.Stat(path)
	//if there is no file, we expect that the program will raise the error
	if err == nil { //the file exists, cause we could stat it... or we don't have perms to read it
		if (fileInfo.Mode & os.ModeNamedPipe) > 0 {
			//the pipe exists... we can read it...
			return path
		} else {
			log.Fatalf("Error: Temp Pipe File \"%s\" exists; however, it is not a named pipe... change source code to accomodate\n", path)
		}
	}
	//we are assuming the file does not exist at this point... either that or we don't have permission to use the known one
	err = syscall.Mkfifo(path, 0644) //owner can write, but everybody can read
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	} //if no errors occured, we made the temp pipe... or we already had it
	return path
}

func main() {
	home := os.Getenv("HOME")
	file, err := ioutil.ReadFile(home + configDir + "api.key")
	if err != nil {
		log.Fatalf("Please ensure that the your API Key is included in a file named \"api.key\" under your $HOME/.config/weather-app directory\nError: %s\n")
	}
	key := strings.Trim(string(file), "\n")
	res, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + key + "&q=auto:ip")
	if err != nil {
		log.Fatalf("Error: %s")
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error: %s")
	}
	//at this point we know that the string is a JSON object
	contentStr := string(content)
	pretty.Printf("Current Temp: %s\n", gjson.Get(contentStr, "current.temp_c").String())
	//write temp to pipe
}
