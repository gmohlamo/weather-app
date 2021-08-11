package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/tidwall/gjson"
)

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
	contentStr := string(content)
	pretty.Printf("Current Temp: %s\n", gjson.Get(contentStr, "current.temp_c").String())
	pretty.Printf("Current Feel: %s\n", gjson.Get(contentStr, "current.feelslike_c").String())
}
