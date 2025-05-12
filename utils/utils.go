package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

func Error(c fiber.Ctx, err error) error {
	return c.Status(500).JSON(map[string]string{
		"message": fmt.Sprintf("Error: %v", err.Error()),
	})
}

func GetLocals(c fiber.Ctx, name string, result interface{}) {
	json.Unmarshal([]byte(fmt.Sprintf("%v", c.Locals(name))), &result)
}

func SetLocals(c fiber.Ctx, name string, data interface{}) {
	bytes, _ := json.Marshal(data)
	c.Locals(name, string(bytes))
}
