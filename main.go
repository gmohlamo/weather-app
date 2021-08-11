package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kr/pretty"
)

func getTemp(object map[string]interface{}) {
}

func main() {
	home := os.Getenv("HOME")
	file, err := ioutil.ReadFile(home + "/.config/weather-app/api.key")
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
	var object map[string]interface{}
	json.Unmarshal(content, &object)
	pretty.Println(object)
}
