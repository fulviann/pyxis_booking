package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	apierror "github.com/fulviann/pyxis_booking/back-end/utils/api-error"
	"github.com/fulviann/pyxis_booking/back-end/utils/config"
	"github.com/fulviann/pyxis_booking/back-end/utils/constants"
	contextUtil "github.com/fulviann/pyxis_booking/back-end/utils/context"
	"github.com/fulviann/pyxis_booking/back-end/utils/respond"
)

type Handler interface {
	Login(ctx *gin.Context)
	VerifyToken(ctx *gin.Context)
	Logout(ctx *gin.Context)
	Register(ctx *gin.Context)
	RegisterAdmin(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	GetPersonal(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	AddAvatar(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ResetPasswordSubmit(ctx *gin.Context)
	GoogleAuth(ctx *gin.Context)
	DeleteAvatar(ctx *gin.Context)
}

type handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service, validate *validator.Validate) Handler {
	return &handler{
		service:  service,
		validate: validate,
	}
}
func (h *handler) GetPersonal(ctx *gin.Context) {
	res, err := h.service.GetPersonal(ctx)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, res)
}

func (h *handler) Login(ctx *gin.Context) {
	var input LoginReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	if input.Role == "" {
		input.Role = constants.CUSTOMER
	}

	res, err := h.service.Login(ctx, input, ctx.Writer)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	frontend := ctx.GetHeader("X-Frontend")
	env := config.DEVELOPMENT_ENVIRONMENT
	cookieName := ""
	if frontend == "admin" {
		cookieName = "token_admin"
	} else {
		cookieName = "token_user"
	}

	cookieDomain := ""
	cookieDomainAdmin := ""
	secure := false
	sameSite := http.SameSiteLaxMode

	if env == config.PRODUCTION_ENVIRONMENT {
		cookieDomain = ""
		cookieDomainAdmin = ""
		secure = true
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    res.Token,
		Expires:  res.Expires,
		Path:     "/",
		HttpOnly: true,
		Domain:   cookieDomain,
		Secure:   secure,
		SameSite: sameSite,
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    res.Token,
		Expires:  res.Expires,
		Path:     "/",
		HttpOnly: true,
		Domain:   cookieDomainAdmin,
		Secure:   secure,
		SameSite: sameSite,
	})

	respond.Success(ctx, http.StatusOK, res)
}

func (h *handler) VerifyToken(ctx *gin.Context) {
	respond.Success(ctx, http.StatusOK, VerifyTokenRes{TokenVerified: true})
}

func (h *handler) Logout(ctx *gin.Context) {
	frontend := ctx.GetHeader("X-Frontend")
	env := config.DEVELOPMENT_ENVIRONMENT
	cookieName := ""
	if frontend == "admin" {
		cookieName = "token_admin"
	} else {
		cookieName = "token_user"
	}

	cookieDomain := ""
	cookieDomainAdmin := ""
	secure := false
	sameSite := http.SameSiteLaxMode

	if env == config.PRODUCTION_ENVIRONMENT {
		cookieDomain = ""
		cookieDomainAdmin = ""
		secure = true
		sameSite = http.SameSiteNoneMode
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Domain:   cookieDomain,
		Secure:   secure,
		SameSite: sameSite,
	})

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Domain:   cookieDomainAdmin,
		Secure:   secure,
		SameSite: sameSite,
	})

	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	input := LogoutReq{
		Token:   token.Token,
		Expires: token.Claims.ExpiresAt.Time,
	}

	res, err := h.service.Logout(ctx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, res)
}

func (h *handler) Register(ctx *gin.Context) {

	var input RegisterReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	res, err := h.service.Register(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "customer_email_key") {
			respond.Error(ctx, apierror.DuplicateEmail(input.Email))
			return
		}

		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusCreated, res)
}

func (h *handler) RegisterAdmin(ctx *gin.Context) {

	var input RegisterReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	err := h.validate.Struct(input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	res, err := h.service.RegisterAdmin(ctx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusCreated, res)
}

func (h *handler) ChangePassword(ctx *gin.Context) {
	var input ChangePasswordReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	if err := h.service.ChangePassword(ctx, input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *handler) UpdateProfile(ctx *gin.Context) {
	var input UpdateProfileReq

	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	res, err := h.service.UpdateProfile(ctx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, res)
}

func (h *handler) AddAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	req := AvatarReq{
		AvatarUrl: file,
	}

	newUrl, err := h.service.AddAvatar(ctx.Request.Context(), req)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, gin.H{
		"message":   "Avatar added successfully",
		"avatarUrl": newUrl,
	})
}

func (h *handler) ResetPassword(ctx *gin.Context) {
	var input ResetPasswordReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	if input.Role == "" {
		input.Role = constants.CUSTOMER
	}

	res, err := h.service.ResetPassword(ctx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, res)
}

func (h *handler) ResetPasswordSubmit(ctx *gin.Context) {
	var input ResetPasswordSubmitReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		respond.Error(ctx, apierror.Warn(http.StatusBadRequest, err))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}
	if input.Role == "" {
		input.Role = constants.CUSTOMER
	}

	if err := h.service.ResetPasswordSubmit(ctx, input); err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *handler) GoogleAuth(ctx *gin.Context) {
	email := ctx.GetString("google_email")
	name := ctx.GetString("google_name")
	picture := ctx.GetString("google_picture")
	googleID := ctx.GetString("google_id")

	input := GoogleAuth{
		Name:     name,
		Email:    email,
		Picture:  picture,
		GoogleID: googleID,
	}

	res, err := h.service.GoogleAuth(ctx, input)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}
	fmt.Println("GoogleAuth Response:", res)

	respond.Success(ctx, http.StatusOK, res)
}

func (h *handler) DeleteAvatar(ctx *gin.Context) {
	err := h.service.DeleteAvatar(ctx)
	if err != nil {
		respond.Error(ctx, apierror.FromErr(err))
		return
	}

	respond.Success(ctx, http.StatusOK, gin.H{"message": "Password reset successfully"})
}
