package jwt // Service is an interface from which our api module can access our repository of all our models

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"chat2pay/internal/pkg/logger"
	"chat2pay/internal/repositories"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

type AuthMiddleware interface {
	GenerateToken(user *entities.User) (*string, error)
	ValidateToken() fiber.Handler
}

type authMiddlewares struct {
	userRepo repositories.UserRepository
	cfg      *yaml.Config
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(userRepo repositories.UserRepository, cfg *yaml.Config) AuthMiddleware {
	return &authMiddlewares{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (m *authMiddlewares) GenerateToken(user *entities.User) (*string, error) {
	jwtKey := m.cfg.JWT.Key
	expireMinute := m.cfg.JWT.ExpiredMinute
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(time.Duration(expireMinute) * time.Minute)
	// Create the JWT claims, which includes the email and expiry time
	claims := &Claims{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return nil, err
	}

	return &tokenString, nil
}

func (m *authMiddlewares) ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			jwtKey        = m.cfg.JWT.Key
			authorization = "asdasd" //strings.Split(c.GetReqHeaders()["Authorization"], "Bearer ")
			log           = logger.NewLog("jwt_middleware_generate_token", m.cfg.Logger.Enable)
		)

		if len(authorization) != 2 {
			log.Error("authorization token not valid")
			c.Status(400)
			return c.JSON(presenter.ErrorResponse(errors.New("token not valid!")))
		}

		//token := authorization[1]

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims("token", claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			log.Error(fmt.Sprintf(`authorization failed got %s`, err))
			if err == jwt.ErrSignatureInvalid {
				c.Status(401)
				return c.JSON(presenter.ErrorResponse(err))
			}
			c.Status(400)
			return c.JSON(presenter.ErrorResponse(err))
		}
		if !tkn.Valid {
			log.Error("authorization failed token invalid")
			c.Status(401)
			return c.JSON(presenter.ErrorResponse(err))
		}

		c.Locals(consts.UserId, claims.ID)
		return c.Next()
	}
}
