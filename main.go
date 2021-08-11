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
	fmt.Println("Getting temp file path")
	home := os.Getenv("HOME")
	path := home + configDir + tempPipe
	fileInfo, err := os.Stat(path)
	//if there is no file, we expect that the program will raise the error
	if err == nil { //the file exists, cause we could stat it... or we don't have perms to read it
		if (fileInfo.Mode() & os.ModeNamedPipe) > 0 {
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

func getFeelPipePath() string {
	fmt.Println("Getting feel file path")
	home := os.Getenv("HOME")
	path := home + configDir + tempFeel
	fileInfo, err := os.Stat(path)
	//if there is no file, we expect that the program will raise the error
	if err == nil { //the file exists, cause we could stat it... or we don't have perms to read it
		if (fileInfo.Mode() & os.ModeNamedPipe) > 0 {
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

//open these fifo files
func getFifo(fileName string) *os.File {
	fmt.Println("Opening FIFO")
	f, err := os.OpenFile(fileName, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Error opening FIFO: %s\n", err)
	}
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
		fmt.Println("Making iteration")
		res, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + key + "&q=auto:ip")
		if err != nil {
			log.Fatalf("Error: %s")
		}
		content, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatalf("Error: %s")
		}
		//at this point we know that the string is a JSON object
		contentStr := string(content)
		fmt.Fprintf(tempFile, "%s", gjson.Get(contentStr, "current.temp_c").String())
		fmt.Printf("Temp: %s\n", gjson.Get(contentStr, "current.temp_c").String())
		fmt.Fprintf(feelFile, "%s", gjson.Get(contentStr, "current.feelslike_c").String())
		fmt.Printf("Feel: %s\n", gjson.Get(contentStr, "current.feelslike_c").String())
		time.Sleep(30 * time.Second) //sleep for an hour
		//tempFile.Close()             //close files after writing to them
		//feelFile.Close()
	}
}

func main() {
	fmt.Printf("Starting process")
	writePipes()
}
