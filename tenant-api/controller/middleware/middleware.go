package middleware

import (
	"net/http"
	"strings"

	"github.com/adwinugroho/test-chat-multi-schema/config"
	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(userSvc domain.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, model.NewError(model.ErrorUnauthorized, "Missing authorization header"))
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			tokenStr = strings.TrimSpace(tokenStr)

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, model.NewError(model.ErrorUnauthorizedTokenExpired, "Invalid token"))
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, model.NewError(model.ErrorUnauthorizedTokenExpired, "Invalid token claims"))
			}

			userID, ok := claims["id"].(string)
			if !ok || userID == "" {
				return c.JSON(http.StatusUnauthorized, model.NewError(model.ErrorUnauthorized, "Invalid user ID"))
			}

			user, err := userSvc.GetUserByID(c.Request().Context(), userID)
			if err != nil {
				logger.LogError("AuthenticationMiddleware error is not user not found: " + err.Error())
				return c.JSON(http.StatusUnauthorized, model.NewError(model.ErrorUnauthorized, "Unexpected error please retry again"))
			}

			if user == nil {
				return c.JSON(http.StatusUnauthorized, model.NewError(model.ErrorUnauthorized, "User not found"))
			}

			c.Set("user", user)
			c.Set("user_role", user.Role)

			return next(c)
		}
	}
}
