package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var testcase []*Message

type Message struct {
	Datetime string `json:"datetime"`
	Content  string `json:"content"`
}

func parseLine(line string) (*Message, error) {

	match, err := regexp.MatchString(`\d+\/\d+\/\d+\, \d{2}\:\d{2} \- .*`, line)

	if err != nil || !match {
		return nil, errors.New("Invalid Line")
	}

	sep := " - "
	s := strings.Split(line, sep)
	datetime, content := s[0], s[1]

	message := Message{
		Datetime: datetime,
		Content:  content,
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
		if err == nil {
			testcase = append(testcase, message)
			log.Println(message)
		}
	}

	defer file.Close()

}

func run_testcase() {
	log.Println("Running Test Case...")

	for _, message := range testcase {
		request(message)
		time.Sleep(1 * time.Second)
	}
}

func request(message *Message) {
	log.Println("Making Request...")

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(message)

	log.Println("Requesting processor...")
	log.Println(body)

	result, err := http.Post(
		"http://proxy.default.svc.cluster.local:8082/message",
		"application/json",
		body)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(result)
}

func main() {
	log.Println("Starting Runner!")

	load_testcase("testcase_1.txt")
	run_testcase()
}
