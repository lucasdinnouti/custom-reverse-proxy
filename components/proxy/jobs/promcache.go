package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var PROMETHEUS_URL = "http://prometheus.monitoring.svc.cluster.local:9090/api/v1/query"
var CPU_QUERY = "(sum(rate(container_cpu_usage_seconds_total{container=~'processor-.'}[1m])) by (container) / (sum(container_spec_cpu_quota{container=~'processor-.'}) by (container) / sum(container_spec_cpu_period{container=~'processor-.'}) by (container))) * 100"
var MEMORY_QUERY = "(sum(container_memory_usage_bytes{container=~'processor-.'}) by (container) / sum (container_spec_memory_limit_bytes{container=~'processor-.'}) by (container)) * 100"

type PrometheusResponseResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type PrometheusResponseData struct {
	ResultType string                     `json:"resultType"`
	Result     []PrometheusResponseResult `json:"result"`
}

type PrometheusResponse struct {
	Status string                 `json:"status"`
	Data   PrometheusResponseData `json:"data"`
}

type PromCache struct {
	CpuUsage        map[string]float32
	MemoryUsage     map[string]float32
	IntervalSeconds int
}

func NewPromCache(hosts []string, intervalSeconds int) *PromCache {
	cpuUsage := map[string]float32{}
	memoryUsage := map[string]float32{}

	for _, h := range hosts {
		cpuUsage[h] = 0
		memoryUsage[h] = 0
	}

	return &PromCache{
		CpuUsage:        cpuUsage,
		MemoryUsage:     memoryUsage,
		IntervalSeconds: intervalSeconds}
}

func (p *PromCache) requestPrometheus(query string) (PrometheusResponse, error) {
	dt := fmt.Sprintf("%d", time.Now().Unix())

	resp, err := http.PostForm(
		PROMETHEUS_URL,
		url.Values{"query": {query}, "time": {dt}})

	if err != nil {
		return PrometheusResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return PrometheusResponse{}, err
	}

	var result PrometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return PrometheusResponse{}, err
	}

	return result, nil
}

func (p *PromCache) getCpuUsage() error {
	resp, err := p.requestPrometheus(CPU_QUERY)

	if err != nil {
		return err
	}

	for _, r := range resp.Data.Result {
		host := strings.Replace(r.Metric["container"], "processor-", "", 1)

		usageStr, ok := r.Value[1].(string)
		if !ok {
			return errors.New("cannot cast prometheus cpu return value")
		}

		usage, err := strconv.ParseFloat(usageStr, 32)
		if err != nil {
			return err
		}

		p.CpuUsage[host] = float32(usage)
	}

	return nil
}

func (p *PromCache) getMemoryUsage() error {
	resp, err := p.requestPrometheus(MEMORY_QUERY)

	if err != nil {
		return err
	}

	for _, r := range resp.Data.Result {
		host := strings.Replace(r.Metric["container"], "processor-", "", 1)
		usageStr, ok := r.Value[1].(string)

		if !ok {
			return errors.New("cannot cast prometheus memory return value")
		}

		usage, err := strconv.ParseFloat(usageStr, 32)
		if err != nil {
			return err
		}

		p.MemoryUsage[host] = float32(usage)
	}

	return nil
}

func (p *PromCache) Run() {
	for {
		log.Println("Starting loading prometheus data...")
		err := p.getCpuUsage()
		if err != nil {
			log.Fatal("Cannot get cpu usage", err)
		}

		err = p.getMemoryUsage()
		if err != nil {
			log.Fatal("Cannot get memory usage", err)
		}

		log.Println("PromCache sleep...")
		time.Sleep(time.Duration(p.IntervalSeconds) * time.Second)
	}
}
