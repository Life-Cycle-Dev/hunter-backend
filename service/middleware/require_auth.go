package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
	"hunter-backend/repository"
	"hunter-backend/util"
	"time"
)

func RequireAuth(db *gorm.DB, config config.AppConfig, tokenType entity.JsonTokenType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		tokenParsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, fiber.ErrUnauthorized
			}
			pubKey, err := util.LoadRSAPublicKey()
			if err != nil {
				return nil, err
			}
			return pubKey, nil
		})

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid exp claim",
			})
		}

		if time.Now().Unix() > int64(exp) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}

		if err != nil || !tokenParsed.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		jsonWebTokenRepository := repository.ProvideJsonWebTokenRepository(db, config)
		jwtId := claims["sub"].(string)
		jwtEnt, err := jsonWebTokenRepository.GetTokenById(jwtId)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token Not Found",
			})
		}

		if jwtEnt.Type != tokenType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Wrong token type",
			})
		}

		if jwtEnt.Revoked {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token is revoked",
			})
		}

		accountRepository := repository.ProvideUserRepository(db, config)
		user, err := accountRepository.FindById(jwtEnt.UserId)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Account Not Found",
			})
		}

		c.Locals("user", user)
		c.Locals("token", jwtEnt)
		return c.Next()
	}
}
