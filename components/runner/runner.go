package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"
	"runtime/pprof"
)

var tpsAtIteration = []int{10, 20, 50, 100}
var testcase []*Message

var validLineRegex = regexp.MustCompile(`\d+\/\d+\/\d+\, \d{2}\:\d{2} \-`)
var mediaTypeRegex = regexp.MustCompile(`([A-Z]{3})-.{16}(jpg|opus) \(file attached\)`)

type ContentType uint8
const (
  Text ContentType = iota
  Image
  Audio
  Unknown
)

type Message struct {
	Datetime string 	 `json:"datetime"`
	Content  string 	 `json:"content"`
	Type 	 ContentType `json:"type"`
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ParseLine(line string) (*Message, error) {

	match := validLineRegex.MatchString(line)

	if !match {
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

func ResolveContentType(content string) ContentType {
	contentType := Text

	match := mediaTypeRegex.MatchString(content)

	if match {
		switch mediaTypeRegex.FindStringSubmatch(content)[0][0:3] {
		case "IMG":
			return Image
		case "PTT":
			return Audio
		default:
			return Text
		}
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

	for _, tps := range tpsAtIteration {
		log.Printf("Running at %d TPS", tps)
		limiter := time.Tick(time.Duration(float64(1000 / tps) * float64(time.Millisecond)))

		for _, message := range testcase {
			<-limiter
			log.Println("Requesting ", time.Now())
			request(message)
		}
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
		log.Println("[ERROR] ", err)
	}

	log.Println(result)
}

func main() {
    // Start a background process that checks the threshold every 30 seconds and dumps a heap profile if necessary
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Println("Starting Runner!")

	LoadTestCase("testcase_1.txt")
	RunTestCase()

	select {
		case <-time.After(10 * time.Second):
			log.Println("missed signal")
		case <-ctx.Done():
			stop()
			log.Println("signal received")
	}
}
