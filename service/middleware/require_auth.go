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

		roleRepository := repository.ProvideRoleRepository(db, config)
		if user.RoleId == "" {
			userRole, err := roleRepository.FindByMapping("user")
			if userRole.ID == "" {
				userRoleEnt := &entity.Role{
					Title:   "User",
					Mapping: "user",
				}
				userRole, err = roleRepository.CreateRole(userRoleEnt, nil)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Cannot create user role",
					})
				}
			}
			user.RoleId = userRole.ID
			user, err = accountRepository.UpdateUser(user)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot update user",
				})
			}
		}

		role, err := roleRepository.FindById(user.RoleId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot get user role",
			})
		}

		permissionRepository := repository.ProvidePermissionRepository(db, config)
		permissions, err := permissionRepository.GetByRoleId(role.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot get permissions",
			})
		}

		permissionNames := make([]string, len(permissions))
		for _, permission := range permissions {
			permissionNames = append(permissionNames, permission.Mapping)
		}

		c.Locals("user", user)
		c.Locals("role", role)
		c.Locals("permissions", permissionNames)
		c.Locals("token", jwtEnt)
		return c.Next()
	}
}
