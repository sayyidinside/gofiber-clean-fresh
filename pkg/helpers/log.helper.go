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
	"github.com/natefinch/lumberjack"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	apiLogger    *zap.Logger
	systemLogger *zap.Logger
)

func InitLogger() {
	// Create logs directory if it does not exist
	if err := os.MkdirAll("storage/logs/api", os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll("storage/logs/system", os.ModePerm); err != nil {
		panic(err)
	}

	// Encoder configuration
	cfg := config.AppConfig
	debugMode := cfg.Debug

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Dynamic log filename based on current month
	currentTime := time.Now()
	logAPIFilename := "storage/logs/api/api_" + currentTime.Format("2006-01") + ".log" // Format as YYYY-MM

	// API Logger setup with lumberjack for log rotation
	apiFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logAPIFilename, // Log file path
		MaxAge:     365,            // Keep log files for 365 days (1 year)
		MaxBackups: 12,             // Keep 12 backups (1 per month for a year)
		Compress:   true,           // Compress old log files
	})

	// Dynamic log filename based on current month
	logSysFilename := "storage/logs/system/system_" + currentTime.Format("2006-01") + ".log" // Format as YYYY-MM

	// System Logger setup with lumberjack for log rotation
	systemFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logSysFilename, // Log file path
		MaxAge:     365,            // Keep log files for 365 days (1 year)
		MaxBackups: 12,             // Keep 12 backups (1 per month for a year)
		Compress:   true,           // Compress old log files
	})

	// Console output configuration
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	var apiCores []zapcore.Core
	var systemCores []zapcore.Core

	if debugMode {
		apiCores = append(apiCores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
		systemCores = append(systemCores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
	}
	apiCores = append(apiCores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), apiFileWriter, zapcore.InfoLevel))
	systemCores = append(systemCores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), systemFileWriter, zapcore.InfoLevel))

	apiLogger = zap.New(zapcore.NewTee(apiCores...))
	systemLogger = zap.New(zapcore.NewTee(systemCores...))
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
