package helpers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	apiLogger    *zap.Logger
	systemLogger *zap.Logger
)

func InitLogger() {
	// Create logs directory if it does not exist
	if err := os.MkdirAll("storage/logs", os.ModePerm); err != nil {
		panic(err)
	}

	// Encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Console output configuration
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// API Logger setup
	apiLogFile, err := os.OpenFile("storage/logs/api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	apiFileWriter := zapcore.AddSync(apiLogFile)
	apiCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), apiFileWriter, zapcore.InfoLevel),
	)
	apiLogger = zap.New(apiCore)

	// System Logger setup
	systemLogFile, err := os.OpenFile("storage/logs/system.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	systemFileWriter := zapcore.AddSync(systemLogFile)
	systemCore := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), systemFileWriter, zapcore.InfoLevel),
	)
	systemLogger = zap.New(systemCore)
}

func APILogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Read request body
		var requestBody interface{}
		contentType := c.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Handle multipart form-data
			form, err := c.MultipartForm()
			if err != nil {
				requestBody = fmt.Sprintf("Error parsing form-data: %v", err)
			} else {
				formData := make(map[string]interface{})
				for key, values := range form.Value {
					if key == "password" || key == "re_password" || key == "old_password" {
						formData[key] = "[REDACTED]"
					} else {
						formData[key] = values
					}
				}
				requestBody = formData
			}
		} else {
			bodyBytes := c.Body()
			if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				// Handle form-urlencoded
				formData, err := url.ParseQuery(string(bodyBytes))
				if err != nil {
					requestBody = string(bodyBytes)
				} else {
					jsonFormData := make(map[string]interface{})
					for key, values := range formData {
						if key == "password" || key == "re_password" || key == "old_password" || key == "raw" {
							jsonFormData[key] = "[REDACTED]"
						} else {
							jsonFormData[key] = values
						}
					}
					requestBody = jsonFormData
				}
			} else {
				err := json.Unmarshal(bodyBytes, &requestBody)
				if err != nil {
					requestBody = string(bodyBytes)
				}
			}
		}

		// Let the request proceed
		err := c.Next()

		// Capture response
		statusCode := c.Response().StatusCode()
		responseBody := c.Response().Body()
		endTime := time.Now()

		// Process response body for logging
		var jsonResponseBody map[string]interface{}
		if err := json.Unmarshal(responseBody, &jsonResponseBody); err != nil {
			jsonResponseBody = map[string]interface{}{
				"raw": string(responseBody),
			}
		} else {
			redactFields(jsonResponseBody, []string{"key", "token", "password", "re_password", "old_password", "raw"})
		}

		// Redact Authorization header
		headers := c.GetReqHeaders()
		if _, exists := headers["Authorization"]; exists {
			headers["Authorization"] = []string{"[REDACTED]"}
		}

		// Get other data
		identifier := c.GetRespHeader(fiber.HeaderXRequestID)
		userAgent := string(c.Request().Header.UserAgent())
		clientIP := c.IP()
		endpoint := c.Path()
		queryParams := c.Request().URI().QueryArgs()
		httpMethod := c.Method()
		statusCodeString := strconv.Itoa(statusCode)
		originalURL := c.OriginalURL()

		// Get username from session or set as empty string if nil
		var username string
		if sessionUsername := c.Locals("username"); sessionUsername != nil {
			username = sessionUsername.(string)
		} else {
			username = ""
		}

		// Log message
		message := "API LOG"
		if msg, ok := jsonResponseBody["message"]; ok {
			message = fmt.Sprintf("API Log | %s", msg.(string))
		}

		// Log the request and response
		logFunc := apiLogger.Info
		if statusCode >= 500 {
			logFunc = apiLogger.Error
		} else if statusCode >= 400 {
			logFunc = apiLogger.Warn
		}

		fmt.Println(time.Now())
		logFunc(message,
			zap.String("identifier", identifier),
			zap.Time("timestamp", time.Now()),
			zap.String("http_method", httpMethod),
			zap.Any("request_header", headers),
			zap.Any("query_params", queryParams),
			zap.Any("request_body", requestBody),
			zap.String("response_code", statusCodeString),
			zap.Any("response_body", jsonResponseBody),
			zap.String("endpoint", endpoint),
			zap.String("original_url", originalURL),
			zap.String("user_agent", userAgent),
			zap.String("client_ip", clientIP),
			zap.String("username", username),
			zap.Time("start_time", startTime),
			zap.Time("end_time", endTime),
		)

		return err
	}
}

type LogSystemParam struct {
	Identifier string
	StatusCode int
	Location   string
	Message    string
	StartTime  time.Time
	EndTime    time.Time
	Username   string
	Err        interface{}
}

func LogSystem(logData LogSystemParam) {
	var (
		category         string
		humanTime        = logData.EndTime.Format(time.RFC1123)
		statusCodeString = strconv.Itoa(logData.StatusCode)
	)

	switch {
	case logData.StatusCode >= 500:
		category = "FATAL"
	case logData.StatusCode >= 400:
		category = "ERROR"
	default:
		category = "INFO"
	}

	systemLogger.Info("System Log",
		zap.Time("timestamp", time.Now()),
		zap.String("category", category),
		zap.String("response_code", statusCodeString),
		zap.String("location", logData.Location),
		zap.String("message", logData.Message),
		zap.Time("start_time", logData.StartTime),
		zap.Time("end_time", logData.EndTime),
		zap.String("identifier", logData.Identifier),
		zap.Any("username", logData.Username),
		zap.Any("errors", logData.Err),
		zap.String("human_time", humanTime),
	)
}

// Helper function to redact sensitive fields in a map
func redactFields(data map[string]interface{}, fields []string) {
	for _, field := range fields {
		if _, ok := data[field]; ok {
			data[field] = "[REDACTED]"
		}
	}
	for _, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			redactFields(nestedMap, fields)
		} else if nestedArray, ok := value.([]interface{}); ok {
			for _, item := range nestedArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					redactFields(itemMap, fields)
				}
			}
		}
	}
}

// GetAPILogger returns the initialized apiLogger instance
func GetAPILogger() *zap.Logger {
	if apiLogger == nil {
		InitLogger()
	}
	return apiLogger
}

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func GetFunctionAndStructName(i interface{}) (string, string, string) {
	pc, file, _, ok := runtime.Caller(2) // 2 untuk mengambil caller dari fungsi pemanggil
	if !ok {
		return "", "", ""
	}
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := fullFuncName[strings.LastIndex(fullFuncName, ".")+1:]

	contrName := reflect.TypeOf(i).Elem().Name()

	projectRoot := getProjectRoot()

	relativePath, err := filepath.Rel(projectRoot, file)
	if err != nil {
		relativePath = file
	}

	packagePath := strings.ReplaceAll(relativePath, string(filepath.Separator), "/")
	packagePath = strings.TrimSuffix(packagePath, filepath.Ext(packagePath))

	projectName := filepath.Base(projectRoot)
	if idx := strings.Index(packagePath, projectName); idx != -1 {
		packagePath = packagePath[idx+len(projectName)+1:]
	}

	return funcName, contrName, packagePath
}

func CreateLog(i interface{}) Log {
	funcName, contrName, packagePath := GetFunctionAndStructName(i)
	return Log{
		StartTime: time.Now(),
		Location:  fmt.Sprintf("%s/%s.%s", packagePath, contrName, funcName),
	}
}
