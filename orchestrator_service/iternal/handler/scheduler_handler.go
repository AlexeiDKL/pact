package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"dkl.ru/pact/orchestrator_service/iternal/config"
)

type SchedulerHandler struct{}

type ResultType struct{ ServiceName, Host, Port, Status string }
type JobsType struct{ Host, Port, ServiceName string }

func NewSchedulerHandler() *SchedulerHandler {
	return &SchedulerHandler{}
}

type ServiceChecker struct {
	Host       string
	Port       string
	Client     *http.Client
	MaxRetries int
	Backoff    time.Duration
}

type HealthChecker interface {
	GetStatus() (string, error)
}

func NewServiceChecker(host, port string, maxRetries int, backoff time.Duration) HealthChecker {
	return &ServiceChecker{
		Host:       host,
		Port:       port,
		Client:     &http.Client{Timeout: 2 * time.Second},
		MaxRetries: maxRetries,
		Backoff:    backoff,
	}
}

func (sc *ServiceChecker) GetStatus() (string, error) {
	var err error
	var resp *http.Response
	url := fmt.Sprintf("http://%s:%s/health", sc.Host, sc.Port)

	for attempt := 0; attempt <= sc.MaxRetries; attempt++ {
		resp, err = sc.Client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				return "", readErr
			}
			return string(body), nil
		}

		// подготовка к следующей попытке
		wait := sc.Backoff * time.Duration(attempt+1)
		time.Sleep(wait)
	}

	return "", fmt.Errorf("after %d attempts, last error: %w", sc.MaxRetries+1, err)
}

func worker(jobs <-chan JobsType, results chan<- ResultType, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		checker := NewServiceChecker(job.Host, job.Port, 3, 500*time.Millisecond)
		status, err := checker.GetStatus()
		if err != nil {
			results <- ResultType{ServiceName: job.ServiceName, Host: job.Host, Port: job.Port, Status: fmt.Sprintf("%s:%s — error: %v", job.Host, job.Port, err)}
		} else {
			results <- ResultType{ServiceName: job.ServiceName, Host: job.Host, Port: job.Port, Status: fmt.Sprintf("%s:%s — status: %s", job.Host, job.Port, status)}
		}
	}
}

func (h *SchedulerHandler) Status(w http.ResponseWriter, r *http.Request) {
	services := []JobsType{{
		config.Config.Server.BdService.Host, strconv.Itoa(config.Config.Server.BdService.Port), "BdService"},
		{config.Config.Server.Garant.Host, strconv.Itoa(config.Config.Server.Garant.Port), "Garant"},
		{config.Config.Server.DocumentService.Host, strconv.Itoa(config.Config.Server.DocumentService.Port), "DocumentService"},
		{config.Config.Server.OrchestratorService.Host, strconv.Itoa(config.Config.Server.OrchestratorService.Port), "OrchestratorService"},
	}
	jobs := make(chan JobsType, len(services))
	results := make(chan ResultType, len(services))

	// создаём пул из 3 воркеров
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// раздаём задачи
	for _, svc := range services {
		jobs <- svc
	}
	close(jobs)

	// ждём завершения воркеров и закрываем results
	wg.Wait()
	close(results)

	var output []ResultType
	for res := range results {
		output = append(output, res)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(output)
}
