package loadtest

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var mut sync.Mutex

// Message DTOs (Maybe we can extract this to separate file)

type Message struct {
	Datetime string      `json:"datetime"`
	Content  string      `json:"content"`
	Type     ContentType `json:"type"`
}

type MessageResponse struct {
	RequestedAt  string
	ElapsedNanos string
	Type         string
	InstanceId   string
	InstanceType string
}

type ContentType uint8

const (
	Text ContentType = iota
	Image
	Audio
	Unknown
)

func (e ContentType) String() string {
	switch e {
	case Text:
		return "text"
	case Image:
		return "image"
	case Audio:
		return "audio"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

// End of Message DTOs

var ( //10, 20, 30, 40,
	tpsAtIteration = []int{10, 30, 60, 80, 100, 120, 120, 130, 130, 140, 140}
	// tpsAtIteration  = []int{ 100, 110, 120, 130, 130, 140, 150, 160, 170 }
	testcase        []*Message
	loadtest_result []*MessageResponse

	validLineRegex = regexp.MustCompile(`\d+\/\d+\/\d+\, \d{2}\:\d{2} \-`)
	mediaTypeRegex = regexp.MustCompile(`([A-Z]{3})-.{16}(jpg|opus) \(file attached\)`)
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
		Type:     ResolveContentType(content),
	}

	return &message, nil
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

func RunTestCase(requestDurations *prometheus.HistogramVec) {
	log.Println("Running Test Case...")
	timout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: time.Duration(timout * int(time.Second))}

	for _, tps := range tpsAtIteration {
		log.Printf("Running at %d TPS", tps)
		limiter := time.Tick(time.Duration(float64(1000/tps) * float64(time.Millisecond)))

		for _, message := range testcase {
			<-limiter

			go Request(client, message, requestDurations)
		}

		RecordToCsv(fmt.Sprintf("./loadtest_results/result_%d.csv", tps))
	}
}

func RecordToCsv(name string) {
	log.Println("Writting result to file, ", len(loadtest_result), " records.")

	os.MkdirAll(filepath.Dir(name), 0770)
	file, err := os.Create(name)
	defer file.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	// Using WriteAll
	var data [][]string
	for _, record := range loadtest_result {
		row := []string{record.RequestedAt, record.ElapsedNanos, record.Type, record.InstanceId, record.InstanceType}
		data = append(data, row)
	}
	w.WriteAll(data)
}

func Request(client *http.Client, message *Message, requestDurations *prometheus.HistogramVec) {
	start := time.Now()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(message)

	Check(err)

	before := time.Now()
	log.Println("Requesting processor ", before, body)

	req, err := http.NewRequest(
		"POST",
		"http://proxy.default.svc.cluster.local:8082/message",
		body)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Message-Type", message.Type.String())

	response, err := client.Do(req)

	if err != nil {
		log.Println("Error: ", err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)

		// Reponse is in format ID_SIZE-TYPE, e.g. a_large-gpu, b_small-cpu
		instance := string(bodyBytes)
		s := strings.Split(instance, "_")
		inst_id, inst_type := s[0], s[1]

		(*requestDurations).WithLabelValues(inst_id).Observe(time.Since(start).Seconds())

		mr := MessageResponse{
			RequestedAt:  fmt.Sprintf("%d", before.Unix()),
			ElapsedNanos: fmt.Sprintf("%d", time.Since(before).Nanoseconds()),
			Type:         message.Type.String(),
			InstanceId:   inst_id,
			InstanceType: inst_type,
		}

		log.Printf("message { requested_at: %s, elapsed_nanos: %s, type: %s, inst_id: %s, inst_type: %s }", mr.RequestedAt, mr.ElapsedNanos, mr.Type, mr.InstanceId, mr.InstanceType)

		mut.Lock()
		loadtest_result = append(loadtest_result, &mr)
		mut.Unlock()

		log.Println("Routed to ", instance)
	}
}
