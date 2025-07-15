package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"strings"
)

func ValidateRequest(c *fiber.Ctx, entity any) error {
	body := c.Body()

	if err := json.Unmarshal(body, entity); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(entity); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func GenerateRandomFilename(original string) (string, error) {
	ext := filepath.Ext(original)
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	randomStr := hex.EncodeToString(randomBytes)
	return randomStr + strings.ToLower(ext), nil
}
