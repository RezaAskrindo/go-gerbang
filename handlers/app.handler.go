package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"go-gerbang/types"

	// "github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func ParseBody(c *fiber.Ctx, body interface{}) error {
	if err := c.BodyParser(body); err != nil {
		return BadRequestErrorResponse(c, fmt.Errorf("Failed To Parse Body"))
	}

	return nil
}

func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

func RandomString(length int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	b := make([]byte, length)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[idx.Int64()]
	}
	return string(b)
}

// CONVERT
func StringToInt(stringData string) int {
	n, err := strconv.Atoi(stringData)
	if err != nil {
		return 0
	}
	return n
}

func StringToFloat64(stringData string) float64 {
	n, err := strconv.ParseFloat(stringData, 64)
	if err != nil {
		return 0
	}
	return n
}

// FASTHTTP CLIENT
var Client = fasthttp.Client{}

// UUID google
var UUID = uuid.New()

var (
	MapMicroServiceMutex sync.RWMutex
	MapMicroService      *types.ConfigServices
)

func LoadConfig(filename string) (*types.ConfigServices, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config types.ConfigServices
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return &config, nil
}

// NOT WORKING
// func WatchConfigFile(filename string, done chan bool) {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatalf("Error creating file watcher: %v", err)
// 	}
// 	defer watcher.Close()

// 	err = watcher.Add(filename)
// 	if err != nil {
// 		log.Fatalf("Error adding file to watcher: %v", err)
// 	}

// 	for {
// 		select {
// 		case <-done:
// 			return
// 		case event := <-watcher.Events:
// 			if event.Op&fsnotify.Write == fsnotify.Write {
// 				log.Println("Config file modified, reloading...")
// 				newConfig, err := LoadConfig(filename)
// 				if err == nil {
// 					MapMicroService = newConfig
// 					// SaveToRedis("proxy-route", newConfig)
// 					log.Println("Config reloaded successfully")
// 				} else {
// 					log.Printf("Error reloading config: %v", err)
// 				}
// 			}
// 		case err := <-watcher.Errors:
// 			log.Printf("Watcher error: %v", err)
// 		}
// 	}
// }

// Response
type SuccessStruct struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   *int64      `json:"total"`
}

type ErrorStruct struct {
	Message interface{} `json:"message"`
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
}

func SuccessResponse(c *fiber.Ctx, status bool, message string, data interface{}, total *int64) error {
	return c.Status(fiber.StatusOK).JSON(&SuccessStruct{
		Status:  status,
		Message: message,
		Data:    data,
		Total:   total,
	})
}

func BadRequestErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusBadRequest,
	})
}

func ConflictErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusConflict).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusConflict,
	})
}

func InternalServerErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusInternalServerError,
	})
}

func UnauthorizedErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusUnauthorized,
	})
}

func ForbiddenErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusForbidden,
	})
}

func UnprocessableEntityErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusUnprocessableEntity,
	})
}

func NotFoundErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusNotFound).JSON(&ErrorStruct{
		Message: err.Error(),
		Status:  false,
		Code:    fiber.StatusNotFound,
	})
}

var validate = validator.New()

func ValidateStruct(data interface{}) map[string]map[string]interface{} {
	errors := make(map[string]map[string]interface{})

	tagDescriptionValidation := map[string]string{
		"required": "Tidak Boleh Kosong",
		"email":    "Harus Email",
	}

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[LowerFirstCase(err.StructField())] = map[string]interface{}{
				"invalid": true,
				"desc":    tagDescriptionValidation[err.Tag()],
				"descRaw": err.Tag(),
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func LowerFirstCase(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToLower(string(str[0])) + str[1:]
}
