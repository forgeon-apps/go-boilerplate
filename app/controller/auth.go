package controller

import (
	"fmt"
	"time"

	"github.com/byeblogs/go-boilerplate/app/model"
	repo "github.com/byeblogs/go-boilerplate/app/repository"
	"github.com/byeblogs/go-boilerplate/pkg/config"
	"github.com/byeblogs/go-boilerplate/platform/database"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetNewAccessToken method for create a new access token.
// @Description Create a new access token.
// @Summary Create a new access token
// @Tags Token
// @Accept json
// @Produce json
// @Param login body model.Auth true "Request for token"
// @Success 200 {object} model.TokenResponse "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /v1/token/new [post]
func GetNewAccessToken(c *fiber.Ctx) error {
	login := &model.Auth{}
	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	user, err := userRepo.GetByUsername(login.Username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "username not found"})
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": "user has no password set"})
	}

	if !IsValidPassword([]byte(*user.PasswordHash), []byte(login.Password)) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": "password is wrong"})
	}

	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": "user not active anymore."})
	}

	token, err := GenerateNewAccessToken(user.ID, user.IsAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{
		"msg":          fmt.Sprintf("Token will be expired within %d minutes", config.AppCfg().JWTSecretExpireMinutesCount),
		"access_token": token,
	})
}

func GenerateNewAccessToken(userID uuid.UUID, isAdmin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID.String()
	claims["admin"] = isAdmin
	claims["exp"] = time.Now().
		Add(time.Minute * time.Duration(config.AppCfg().JWTSecretExpireMinutesCount)).
		Unix()

	return token.SignedString([]byte(config.AppCfg().JWTSecretKey))
}

func GeneratePasswordHash(password []byte) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func IsValidPassword(hash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}
