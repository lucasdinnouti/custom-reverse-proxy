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
	"time"
)

var testcase []*Message

var validLineRegex = `\d+\/\d+\/\d+\, \d{2}\:\d{2} \-`
var imageMessageRegex = `IMG-.{16}jpg \(file attached\)`
var audioMessageRegex = `PTT-.{16}opus \(file attached\)`

type Message struct {
	Datetime string `json:"datetime"`
	Content  string `json:"content"`
	Type  	 string `json:"type"`
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ParseLine(line string) (*Message, error) {

	match, err := regexp.MatchString(validLineRegex, line)

	if err != nil || !match {
		return nil, errors.New("invalid line")
	}

	sep := " - "
	s := strings.Split(line, sep)
	datetime, content := s[0], s[1]

	message := Message{
		Datetime: datetime,
		Content:  content,
		Type: ResolveContentType(content),
	}

	return &message, nil
}

func ResolveContentType(content string) string {
	contentType := "text"

	match, err := regexp.MatchString(imageMessageRegex, content)
	if err != nil && match {
		contentType = "image"
	}
	
	match, err = regexp.MatchString(audioMessageRegex, content)
	if err != nil && match {
		contentType = "audio"
	}

	return contentType
}

func LoadTestCase(filename string) {

	file, err := os.Open(filename)
	Check(err)

	fileScanner := bufio.NewScanner(file)
	line := ""

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line = fileScanner.Text()

		message, err := ParseLine(line)
		if err == nil {
			testcase = append(testcase, message)
			log.Println(message)
		}
	}

	defer file.Close()

}

func RunTestCase() {
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
	
	Check(err)

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

	LoadTestCase("testcase_1.txt")
	RunTestCase()
}
