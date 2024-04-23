package main

import (
	"bufio"
	"bytes"
	"errors"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var mutex = &sync.Mutex{}
var testcase []*Message

type Message struct {
	datetime string
    content string
}

func parseLine(line string) (*Message, error) {

	match, err := regexp.MatchString(`\d+\/\d+\/\d+\, \d{2}\:\d{2} \- .*`, line)

    if (err != nil || !match) {
        return nil, errors.New("Invalid Line")
    }

    sep := " - "
	s := strings.Split(line, sep)
    datetime, content := s[0], s[1]

    message := Message{
        datetime: datetime,
        content: content,
    }

	return &message, nil
}

func load_testcase(filename string) {

	file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }

	fileScanner := bufio.NewScanner(file)
	line := ""
 
    fileScanner.Split(bufio.ScanLines)
  
    for fileScanner.Scan() {
		line = fileScanner.Text()
		
		message, err := parseLine(line)
		if (err == nil) {
			testcase = append(testcase, message)
			log.Println(message)
		}
    }

    defer file.Close()

}

func run_testcase() {
	for _, message := range testcase {
		request(message)
	}
}

func request(message *Message) {
	body, err := json.Marshal(message)
	log.Println("Requesting processor...")
	log.Println(body)

	result, err := http.Post(
		"http://proxy.default.svc.cluster.local:8082/test",
		"application/json",
		bytes.NewBuffer(body))

	if err != nil {
		log.Fatalln(err)
	}

    log.Println(result)
}

func main() {
	load_testcase("../../datasets/chat_dos/chat_dos.txt")
	run_testcase()
}
