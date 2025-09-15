package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"dkl.ru/pact/orchestrator_service/iternal/config"
	"dkl.ru/pact/orchestrator_service/iternal/core"
	"dkl.ru/pact/orchestrator_service/iternal/logger"
	"github.com/go-chi/chi/v5"
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

func (h *SchedulerHandler) GetAllVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result := map[string]map[string]int{
		"versions": {
			"en": 1644614400,
			"am": 1644614400,
			"by": 1644614400,
			"kz": 1644614400,
			"kg": 1764614400,
			"ru": 1764614400,
		},
	}
	json.NewEncoder(w).Encode(result)
}

type LanguageDates struct {
	Languages map[string]string `json:"languages"`
}

func (h *SchedulerHandler) GetLanguagesUpdateStatus(w http.ResponseWriter, r *http.Request) {
	var req LanguageDates
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невозможно прочитать JSON", http.StatusBadRequest)
		return
	}

	// Преобразуем даты из запроса в Unix timestamp
	langV := map[string]int{}
	for lang, dateStr := range req.Languages {
		parsedTime, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Неверный формат даты для языка %s: %v", lang, err), http.StatusBadRequest)
			return
		}
		logger.Logger.Info(fmt.Sprintf("GetLanguagesUpdateStatus: язык=%s, дата=%s, timestamp=%d", lang, dateStr, parsedTime.Unix()))
		langV[lang] = int(parsedTime.Unix())
	}

	fmt.Println(langV)

	// Ожидаемые версии
	langVersion := core.GetVersions(allFiles)

	// Сравнение и формирование результата
	updateStatus := make(map[string]bool)
	for lang, expectedVersion := range langVersion {
		if reqVersion, ok := langV[lang]; ok {
			updateStatus[lang] = reqVersion < expectedVersion
		} else {
			updateStatus[lang] = false
		}
	}

	// Ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"languages_update_status": updateStatus,
	})
}

type LanguageRequest struct {
	Languages []string `json:"languages"`
}

type CreateTestVersionRequest struct {
	Language string `json:"lang"`
	AddText  string `json:"text"`
}

var allFiles = map[string]map[string]any{
	"RU": {
		"text":    "agree_RU.txt",
		"json":    "agree_RU.json",
		"version": 1756425600,
	},
	"EN": {
		"text":    "agree_EN.txt",
		"json":    "agree_EN.json",
		"version": 1756512000,
	},
	"AM": {
		"text":    "agree_AM.txt",
		"json":    "agree_AM.json",
		"version": 1756339200,
	},
	"BY": {
		"text":    "agree_BY.txt",
		"json":    "agree_BY.json",
		"version": 1756252800,
	},
	"KZ": {
		"text":    "agree_KZ.txt",
		"json":    "agree_KZ.json",
		"version": 1756166400,
	},
	"KG": {
		"text":    "agree_KG.txt",
		"json":    "agree_KG.json",
		"version": 1756080000,
	},
}

func (h *SchedulerHandler) CreateTestVersion(w http.ResponseWriter, r *http.Request) {
	// Декодируем тело запроса
	// json {"text":"helloWorld","lang":"RU"}
	var req CreateTestVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невозможно прочитать JSON", http.StatusBadRequest)
		return
	}

	lang := strings.ToUpper(strings.TrimSpace(req.Language))
	if len(lang) != 2 {
		http.Error(w, "Некорректный код языка", http.StatusBadRequest)
		return
	}

	jsonOldPath := allFiles[lang]["json"].(string)
	txtOldPath := allFiles[lang]["text"].(string)

	jsonNewPath := strings.ReplaceAll(jsonOldPath, ".json", "_test.json")
	txtNewPath := strings.ReplaceAll(txtOldPath, ".txt", "_test.txt")

	// читаем старый текстовый файл
	oldTxtContent, err := ioutil.ReadFile("C:/project/golang/pact/files/" + txtOldPath)
	if err != nil {
		http.Error(w, "Ошибка при чтении текстового файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// дописываем в него строку из запроса
	newTxtContent := req.AddText + "\n" + string(oldTxtContent)
	// сохраняем в новый файл с суффиксом _test
	if err := ioutil.WriteFile("C:/project/golang/pact/files/"+txtNewPath, []byte(newTxtContent), 0644); err != nil {
		http.Error(w, "Ошибка при записи текстового файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// читаем старый json файл
	oldJsonContent, err := ioutil.ReadFile("C:/project/golang/pact/files/" + jsonOldPath)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// вызываем AddOffsetToPositions с оффсетом len(req.AddText)
	offset := len(req.AddText)
	var rootObj core.Root
	if err := json.Unmarshal(oldJsonContent, &rootObj); err != nil {
		http.Error(w, "Ошибка при парсинге JSON файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	core.AddOffsetToPositions(rootObj.Item, offset)
	// сохраняем в новый файл с суффиксом _test
	newJsonContent, err := json.MarshalIndent(rootObj, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка при сериализации JSON файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := ioutil.WriteFile("C:/project/golang/pact/files/"+jsonNewPath, newJsonContent, 0644); err != nil {
		http.Error(w, "Ошибка при записи JSON файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// добавляем созданные файлы в allFiles[lang]
	// старые файлы сохраняем в allFiles[lang_old]
	allFiles[lang+"_old"] = map[string]any{
		"text":    txtOldPath,
		"json":    jsonOldPath,
		"version": allFiles[lang]["version"],
	}
	versionInt, ok := allFiles[lang]["version"].(int)
	if !ok {
		http.Error(w, "Некорректный тип версии", http.StatusInternalServerError)
		return
	}
	allFiles[lang] = map[string]any{
		"text":    txtNewPath,
		"json":    jsonNewPath,
		"version": versionInt + 1,
	}

	fmt.Println(allFiles)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *SchedulerHandler) DeleteTestVersion(w http.ResponseWriter, r *http.Request) {
	// читаем весь allFiles
	// для каждого языка, у которого есть суффикс _old
	// удаляем файлы с суффиксом _test
	// восстанавливаем файлы из allFiles[lang_old] в allFiles[lang]
	for lang, files := range allFiles {
		if strings.Contains(lang, "_old") {
			fmt.Println("Удаляем тестовую версию для языка: " + lang)
			originalLang := strings.TrimSuffix(lang, "_old")
			testTxtPath := strings.ReplaceAll(files["text"].(string), ".txt", "_test.txt")
			testJsonPath := strings.ReplaceAll(files["json"].(string), ".json", "_test.json")
			// удаляем тестовые файлы
			if err := os.Remove("C:/project/golang/pact/files/" + testTxtPath); err != nil {
				logger.Logger.Error("Ошибка при удалении тестового текстового файла: " + err.Error())
			}
			if err := os.Remove("C:/project/golang/pact/files/" + testJsonPath); err != nil {
				logger.Logger.Error("Ошибка при удалении тестового JSON файла: " + err.Error())
			}
			// восстанавливаем оригинальные файлы в allFiles
			versionInt, ok := allFiles[lang]["version"].(int)
			if !ok {
				http.Error(w, "Некорректный тип версии", http.StatusInternalServerError)
				return
			}

			allFiles[originalLang] = map[string]any{
				"text":    files["text"],
				"json":    files["json"],
				"version": versionInt + 2,
			}
			// удаляем запись с суффиксом _old
			delete(allFiles, lang)
			// удаляем файлы с суффиксом _test
			logger.Logger.Info("Удалены тестовые файлы для языка: " + originalLang)
		}
	}

	fmt.Println(allFiles)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *SchedulerHandler) GetFileList(w http.ResponseWriter, r *http.Request) {
	// Декодируем тело запроса
	var req LanguageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невозможно прочитать JSON", http.StatusBadRequest)
		return
	}

	// Формируем ответ только по нормализованным языкам
	selectedFiles := make(map[string]map[string]any)
	for _, lang := range req.Languages {
		normalized := strings.ToUpper(strings.TrimSpace(lang))
		if len(normalized) != 2 {
			continue // Пропускаем некорректные коды
		}
		if files, ok := allFiles[normalized]; ok {
			selectedFiles[normalized] = files
		}
	}

	logger.Logger.Info(fmt.Sprintf("GetFileList: requested=%v, returned=%v", req.Languages, selectedFiles))

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectedFiles)
}

func (h *SchedulerHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Получаем имя файла из path-параметра или query-параметра
	fileName := chi.URLParam(r, "file")
	if fileName == "" {
		fileName = r.URL.Query().Get("file")
	}

	if fileName == "" {
		logger.Logger.Warn("DownloadFile: не указано имя файла")
		http.Error(w, "Не указано имя файла", http.StatusBadRequest)
		return
	}

	// Защита от path traversal
	if strings.Contains(fileName, "..") || strings.ContainsAny(fileName, "/\\") {
		logger.Logger.Warn(fmt.Sprintf("DownloadFile: попытка небезопасного доступа к файлу: %s", fileName))
		http.Error(w, "Недопустимое имя файла", http.StatusBadRequest)
		return
	}

	// Путь к директории с файлами
	basePath := "C:/project/golang/pact/files" // TODO: вынести в конфиг
	fullPath := filepath.Join(basePath, fileName)

	logger.Logger.Info(fmt.Sprintf("DownloadFile: запрошен файл %s", fileName))
	logger.Logger.Debug(fmt.Sprintf("DownloadFile: полный путь %s", fullPath))

	// Проверка существования файла
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Logger.Error(fmt.Sprintf("DownloadFile: файл не найден: %s", fullPath))
			http.Error(w, "Файл не найден", http.StatusNotFound)
		} else {
			logger.Logger.Error(fmt.Sprintf("DownloadFile: ошибка при доступе к файлу: %v", err))
			http.Error(w, "Ошибка при доступе к файлу", http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		logger.Logger.Warn(fmt.Sprintf("DownloadFile: путь указывает на директорию: %s", fullPath))
		http.Error(w, "Ожидался файл, но получена директория", http.StatusBadRequest)
		return
	}

	// Открываем файл
	file, err := os.Open(fullPath)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("DownloadFile: ошибка при открытии файла: %v", err))
		http.Error(w, "Ошибка при открытии файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Устанавливаем заголовки
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Отправляем файл
	if _, err := io.Copy(w, file); err != nil {
		logger.Logger.Error(fmt.Sprintf("DownloadFile: ошибка при отправке файла: %v", err))
		http.Error(w, "Ошибка при отправке файла", http.StatusInternalServerError)
	}
}
