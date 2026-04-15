package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/time/rate"

	"github.com/devanadindra/signlink-mobile/back-end/domains/user"
	apierror "github.com/devanadindra/signlink-mobile/back-end/utils/api-error"
	"github.com/devanadindra/signlink-mobile/back-end/utils/config"
	"github.com/devanadindra/signlink-mobile/back-end/utils/constants"
	contextUtil "github.com/devanadindra/signlink-mobile/back-end/utils/context"
	"github.com/devanadindra/signlink-mobile/back-end/utils/logger"
	"github.com/devanadindra/signlink-mobile/back-end/utils/respond"
)

type Middlewares interface {
	AddRequestId(ctx *gin.Context)
	Logging(ctx *gin.Context)
	BasicAuth(ctx *gin.Context)
	JWT(roles ...constants.ROLE) func(ctx *gin.Context)
	Recover(ctx *gin.Context)
	RateLimiter(ctx *gin.Context)
	OptionalJWT(roles ...constants.ROLE) gin.HandlerFunc
	GoogleAuth() gin.HandlerFunc
}

type middlewares struct {
	conf        *config.Config
	rateLimiter *rate.Limiter
	userService user.Service
}

// Constructor untuk middlewares
func NewMiddlewares(conf *config.Config, userService user.Service) Middlewares {
	return &middlewares{
		conf:        conf,
		rateLimiter: rate.NewLimiter(rate.Limit(conf.RateLimiter.Rps), conf.RateLimiter.Bursts),
		userService: userService,
	}
}

func (m *middlewares) AddRequestId(ctx *gin.Context) {
	requestId := uuid.New()
	ctx = contextUtil.GinWithCtx(ctx, contextUtil.SetRequestId(ctx, requestId))
	ctx.Header("Request-Id", requestId.String())
	ctx.Next()
}

func (m *middlewares) Logging(ctx *gin.Context) {
	start := time.Now()
	reqPayload := getRequestPayload(ctx)

	ctx.Next()

	logPayload := logger.LogPayload{
		Method:         ctx.Request.Method,
		Path:           ctx.Request.URL.Path,
		StatusCode:     ctx.Writer.Status(),
		Took:           time.Since(start),
		RequestPayload: reqPayload,
	}

	var err error
	errAny, ok := ctx.Get("error")
	if !ok {
		err = nil
	} else {
		err, ok = errAny.(error)
		if !ok {
			err = nil
		}
	}

	logger.Log(ctx, logPayload, err)
}

func (m *middlewares) BasicAuth(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte(m.conf.Auth.Basic.Username))
		expectedPasswordHash := sha256.Sum256([]byte(m.conf.Auth.Basic.Password))
		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)
		if usernameMatch && passwordMatch {
			ctx.Next()
			return
		}
	}
	respond.Error(ctx, apierror.Unauthorized())
}

func (m *middlewares) JWT(roles ...constants.ROLE) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var tokenStr string

		frontend := ctx.GetHeader("X-Frontend")
		var cookieName string
		if frontend == "admin" {
			cookieName = "token_admin"
		} else {
			cookieName = "token_user"
		}

		cookie, err := ctx.Request.Cookie(cookieName)
		if err != nil {
			authorization := ctx.GetHeader(constants.AUTHORIZATION)
			if authorization == "" {
				authorization = ctx.GetHeader(constants.AUTH)
			}
			authorizationSplit := strings.Split(authorization, " ")
			if len(authorizationSplit) < 2 {
				respond.Error(ctx, apierror.Unauthorized())
				return
			}
			tokenStr = authorizationSplit[1]
		} else {
			tokenStr = cookie.Value
		}

		// Parse token & claims
		claims := constants.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.conf.Auth.JWT.SecretKey), nil
		})
		if err != nil || !token.Valid {
			respond.Error(ctx, apierror.Unauthorized())
			return
		}

		// Validasi ke DB blacklist
		err = m.userService.ValidateToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			respond.Error(ctx, apierror.Unauthorized())
			return
		}

		// Cek role
		if !slices.Contains(roles, claims.Role) {
			respond.Error(ctx, apierror.Unauthorized())
			return
		}

		// Simpan claims ke context
		ctx = contextUtil.GinWithCtx(ctx, contextUtil.SetTokenClaims(ctx, constants.Token{
			Token:  tokenStr,
			Claims: claims,
		}))

		ctx.Next()
	}
}

func (m *middlewares) OptionalJWT(roles ...constants.ROLE) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenStr string

		frontend := ctx.GetHeader("X-Frontend")
		var cookieName string
		if frontend == "admin" {
			cookieName = "token_admin"
		} else {
			cookieName = "token_user"
		}

		// Cek cookie
		cookie, err := ctx.Request.Cookie(cookieName)
		if err != nil {
			authorization := ctx.GetHeader(constants.AUTHORIZATION)
			if authorization == "" {
				authorization = ctx.GetHeader(constants.AUTH)
			}
			authorizationSplit := strings.Split(authorization, " ")
			if len(authorizationSplit) < 2 {
				ctx.Next()
				return
			}
			tokenStr = authorizationSplit[1]
		} else {
			tokenStr = cookie.Value
		}

		// Parse token & claims
		claims := constants.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.conf.Auth.JWT.SecretKey), nil
		})
		if err != nil || !token.Valid {
			ctx.Next()
			return
		}

		// validasi token ke DB blacklist
		err = m.userService.ValidateToken(ctx.Request.Context(), tokenStr)
		if err != nil {
			ctx.Next()
			return
		}

		if len(roles) > 0 && !slices.Contains(roles, claims.Role) {
			ctx.Next()
			return
		}

		ctx = contextUtil.GinWithCtx(ctx, contextUtil.SetTokenClaims(ctx, constants.Token{
			Token:  tokenStr,
			Claims: claims,
		}))

		ctx.Next()
	}
}

func (m *middlewares) GoogleAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idToken := ctx.GetHeader(constants.AUTHORIZATION)
		if idToken == "" {
			var body struct {
				IDToken string `json:"id_token"`
			}
			if err := ctx.ShouldBindJSON(&body); err != nil || body.IDToken == "" {
				respond.Error(ctx, apierror.Unauthorized())
				ctx.Abort()
				return
			}
			idToken = body.IDToken
		}

		payload, err := idtoken.Validate(ctx, idToken, m.conf.GoogleAuth.ClientID)
		if err != nil {
			respond.Error(ctx, apierror.Unauthorized())
			ctx.Abort()
			return
		}

		email, _ := payload.Claims["email"].(string)
		name, _ := payload.Claims["name"].(string)
		picture, _ := payload.Claims["picture"].(string)
		googleID := payload.Subject

		ctx.Set("google_email", email)
		ctx.Set("google_name", name)
		ctx.Set("google_picture", picture)
		ctx.Set("google_id", googleID)

		ctx.Next()
	}
}

func (m *middlewares) Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			respond.Error(ctx, apierror.NewError(http.StatusInternalServerError, fmt.Sprintf("Panic : %v", r)))
		}
	}()
	ctx.Next()
}

func (m *middlewares) RateLimiter(ctx *gin.Context) {
	if !m.rateLimiter.Allow() {
		respond.Error(ctx, apierror.NewWarn(http.StatusTooManyRequests, "Too many request"))
		return
	}
	ctx.Next()
}
