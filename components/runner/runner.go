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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tpsAtIteration = []int{10, 20, 50, 100}
	testcase []*Message
	
	validLineRegex = regexp.MustCompile(`\d+\/\d+\/\d+\, \d{2}\:\d{2} \-`)
	mediaTypeRegex = regexp.MustCompile(`([A-Z]{3})-.{16}(jpg|opus) \(file attached\)`)

	requestDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "full_request_latency",
		Help: "Latency of requests to processor",
		Buckets: []float64{0.001, 0.005, 0.010, 0.100, 0.500, 1, 2, 3, 4, 5, 6, 8, 9, 10}})
)

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

	s := strings.Split(line, " - ")
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
			request(message)
		}
	}
}

func request(message *Message) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(message)
	
	Check(err)

	log.Println("Requesting processor ", time.Now(), body)

	before := time.Now() 
	timer := prometheus.NewTimer(requestDurations)

	result, err := http.Post(
		"http://proxy.default.svc.cluster.local:8082/message",
		"application/json",
		body)

	timer.ObserveDuration()
	log.Println("Time elapsed", time.Since(before))	

	if err != nil {
		log.Println("Error: ", err)
	}

	log.Println("Response: ", result)
}

func main() {

	prometheus.MustRegister(requestDurations)

	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(":8081", nil)

	log.Println("Starting Runner!")

	LoadTestCase("testcase_1.txt")
	RunTestCase()
}
