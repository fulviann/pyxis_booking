package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/fulviann/pyxis_booking/back-end/database"
	apierror "github.com/fulviann/pyxis_booking/back-end/utils/api-error"
	"github.com/fulviann/pyxis_booking/back-end/utils/config"
	"github.com/fulviann/pyxis_booking/back-end/utils/constants"
	contextUtil "github.com/fulviann/pyxis_booking/back-end/utils/context"
	"github.com/fulviann/pyxis_booking/back-end/utils/dbselector"
	fileutils "github.com/fulviann/pyxis_booking/back-end/utils/file"
)

type Service interface {
	Login(ctx context.Context, input LoginReq, w http.ResponseWriter) (res *LoginRes, err error)
	Logout(ctx context.Context, input LogoutReq) (res *LogoutRes, err error)
	ValidateToken(ctx context.Context, token string) (err error)
	Register(ctx context.Context, input RegisterReq) (res *Customer, err error)
	RegisterAdmin(ctx context.Context, input RegisterReq) (res *Admin, err error)
	ChangePassword(ctx context.Context, input ChangePasswordReq) error
	GetPersonal(ctx context.Context) (*PersonalRes, error)
	UpdateProfile(ctx context.Context, input UpdateProfileReq) (res *PersonalRes, err error)
	AddAvatar(ctx context.Context, req AvatarReq) (string, error)
	ResetPassword(ctx context.Context, req ResetPasswordReq) (res *ResetPasswordRes, err error)
	ResetPasswordSubmit(ctx context.Context, req ResetPasswordSubmitReq) error
	GoogleAuth(ctx context.Context, input GoogleAuth) (*LoginRes, error)
	DeleteAvatar(ctx context.Context) error
}

type service struct {
	authConfig config.Auth
	dbSelector *dbselector.DBService
	CustomerDB *database.CustomerDB
	AdminDB    *database.AdminDB
}

func NewService(config *config.Config, dbSelector *dbselector.DBService, CustomerDB *database.CustomerDB, AdminDB *database.AdminDB) Service {
	return &service{
		authConfig: config.Auth,
		dbSelector: dbSelector,
		CustomerDB: CustomerDB,
		AdminDB:    AdminDB,
	}
}

func (s *service) GetPersonal(ctx context.Context) (*PersonalRes, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}

	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return nil, err
	}

	res := &PersonalRes{}

	switch token.Claims.Role {
	case constants.ADMIN:
		var admin Admin
		err = db.WithContext(ctx).
			Where("id = ?", token.Claims.UserID).
			First(&admin).Error
		if err != nil {
			return nil, err
		}

		res.ID = admin.ID
		res.Name = admin.Name
		res.Email = admin.Email
		res.AvatarUrl = admin.AvatarUrl
		res.CreatedAt = admin.CreatedAt
		res.UpdatedAt = admin.UpdatedAt

		return res, nil

	case constants.CUSTOMER:
		var customer Customer
		err = db.WithContext(ctx).
			Where("id = ?", token.Claims.UserID).
			First(&customer).Error
		if err != nil {
			return nil, err
		}

		res.ID = customer.ID
		res.Name = customer.Name
		res.Email = customer.Email
		res.AvatarUrl = customer.AvatarUrl
		res.GoogleID = customer.GoogleID
		res.HasPassword = customer.HasPassword
		res.CreatedAt = customer.CreatedAt
		res.UpdatedAt = customer.UpdatedAt

		return res, nil

	default:
		return nil, errors.New("invalid role")
	}
}

func (s *service) Login(ctx context.Context, input LoginReq, w http.ResponseWriter) (*LoginRes, error) {
	var err error
	var userID uuid.UUID

	switch input.Role {
	case constants.ADMIN:
		var admin Admin
		db := s.AdminDB.DB
		err = db.WithContext(ctx).Where("email = ?", input.Email).First(&admin).Error
		if err == nil {
			if !comparePassword(admin.Password, input.Password) {
				return nil, apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCredentials)
			}
			userID = admin.ID
		}
	case constants.CUSTOMER:
		var customer Customer
		db := s.CustomerDB.DB
		err = db.WithContext(ctx).Where("email = ?", input.Email).First(&customer).Error
		if err == nil {
			if !comparePassword(customer.Password, input.Password) {
				return nil, apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCredentials)
			}
			userID = customer.ID
		}
	default:
		return nil, apierror.NewWarn(http.StatusBadRequest, "role tidak valid")
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewWarn(http.StatusUnauthorized, ErrEmailNotFound)
		}
		return nil, apierror.FromErr(err)
	}

	expirationTime := time.Now().Add(s.authConfig.JWT.ExpireIn)
	claims := &constants.JWTClaims{
		UserID: userID,
		Email:  input.Email,
		Role:   input.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.authConfig.JWT.SecretKey))
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	return &LoginRes{
		Role:    input.Role,
		Token:   tokenString,
		Expires: expirationTime,
	}, nil
}

func (s *service) Logout(ctx context.Context, input LogoutReq) (res *LogoutRes, err error) {
	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Create(InvalidToken{
		Token:   input.Token,
		Expires: input.Expires,
	}).Error
	if err != nil {
		return nil, err
	}

	return &LogoutRes{
		LoggedOut: true,
	}, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) error {
	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		db = s.CustomerDB.DB
	}

	var invalidToken InvalidToken
	if err := db.WithContext(ctx).Where("token = ?", token).First(&invalidToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	return errors.New("token is blacklisted")
}

func (s *service) Register(ctx context.Context, input RegisterReq) (res *Customer, err error) {
	db := s.CustomerDB.DB

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	customer := Customer{
		Name:        input.Name,
		Email:       input.Email,
		Password:    hashedPassword,
		HasPassword: true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert into DB
	if err := db.WithContext(ctx).Create(&customer).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	return &customer, nil
}

func (s *service) RegisterAdmin(ctx context.Context, input RegisterReq) (res *Admin, err error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, err
	}

	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return nil, err
	}
	loggedInUsername := token.Claims.Email
	if loggedInUsername != "owner" {
		return nil, errors.New("unauthorized: only owners can register admins")
	}

	// Hash password
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	// Build admin object
	admin := Admin{
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into DB
	if err := db.WithContext(ctx).Create(&admin).Error; err != nil {
		return nil, apierror.FromErr(err)
	}

	return &admin, nil
}

func (s *service) ChangePassword(ctx context.Context, input ChangePasswordReq) error {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return err
	}

	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return err
	}

	var userID = token.Claims.UserID
	hashedPassword, err := hashPassword(input.NewPassword)
	if err != nil {
		return apierror.FromErr(err)
	}

	var role = token.Claims.Role

	switch role {
	case constants.ADMIN:
		var admin Admin
		if err = db.WithContext(ctx).Where("id = ?", userID).First(&admin).Error; err == nil {
			if !comparePassword(admin.Password, input.CurrentPassword) {
				return apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCurPassword)
			}
			admin.Password = hashedPassword
			if err = db.WithContext(ctx).Save(&admin).Error; err != nil {
				return err
			}
		}
	case constants.CUSTOMER:
		var customer Customer
		if err = db.WithContext(ctx).Where("id = ?", userID).First(&customer).Error; err == nil {
			if customer.HasPassword {
				if !comparePassword(customer.Password, input.CurrentPassword) {
					return apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCurPassword)
				}
				customer.Password = hashedPassword
				if err = db.WithContext(ctx).Save(&customer).Error; err != nil {
					return err
				}
			}

			if customer.GoogleID != "" && !customer.HasPassword {
				customer.Password = hashedPassword
				customer.HasPassword = true
				if err = db.WithContext(ctx).Save(&customer).Error; err != nil {
					return err
				}
			}
		}
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.NewWarn(http.StatusUnauthorized, ErrInvalidCurPassword)
		}
		return err
	}

	return nil
}

func (s *service) UpdateProfile(ctx context.Context, input UpdateProfileReq) (*PersonalRes, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.FromErr(err)
	}
	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return nil, err
	}
	userID := token.Claims.UserID
	role := token.Claims.Role
	var res PersonalRes

	switch role {
	case constants.CUSTOMER:
		var customer Customer
		if err := db.WithContext(ctx).
			Where("id = ?", userID).
			First(&customer).Error; err != nil {
			return nil, apierror.FromErr(err)
		}

		if customer.GoogleID != "" {
			if input.Email != customer.Email {
				return nil, apierror.ErrGoogleEmailLocked()
			}
		}

		customer.Name = input.Name
		customer.UpdatedAt = time.Now()

		if customer.GoogleID == "" {
			customer.Email = input.Email
		}

		if err := db.WithContext(ctx).Save(&customer).Error; err != nil {
			return nil, apierror.FromErr(err)
		}

		res = PersonalRes{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			AvatarUrl: customer.AvatarUrl,
			UpdatedAt: customer.UpdatedAt,
		}

	case constants.ADMIN:
		var admin Admin
		err = db.WithContext(ctx).
			Where("id = ?", userID).
			First(&admin).Error
		if err != nil {
			return nil, apierror.FromErr(err)
		}

		admin.Name = input.Name
		admin.Email = input.Email
		admin.UpdatedAt = time.Now()

		if err := db.WithContext(ctx).Save(&admin).Error; err != nil {
			return nil, apierror.FromErr(err)
		}
		res = PersonalRes{
			ID:        admin.ID,
			Name:      admin.Name,
			Email:     admin.Email,
			AvatarUrl: admin.AvatarUrl,
			UpdatedAt: admin.UpdatedAt,
		}
	}

	return &res, nil
}

func (s *service) AddAvatar(ctx context.Context, req AvatarReq) (string, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return "", err
	}
	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return "", err
	}

	userID := token.Claims.UserID
	role := token.Claims.Role

	var oldAvatar string
	switch role {
	case constants.ADMIN:
		var admin Admin
		err = db.WithContext(ctx).Where("id = ?", userID).First(&admin).Error
		if err != nil {
			return "", apierror.FromErr(err)
		}
		oldAvatar = admin.AvatarUrl

	case constants.CUSTOMER:
		var customer Customer
		err = db.WithContext(ctx).Where("id = ?", userID).First(&customer).Error
		if err != nil {
			return "", apierror.FromErr(err)
		}
		oldAvatar = customer.AvatarUrl
	}

	file := req.AvatarUrl
	ext := filepath.Ext(file.Filename)

	filename, err := fileutils.GenerateMediaName(userID.String())
	if err != nil {
		return "", apierror.FromErr(err)
	}

	filename = fmt.Sprintf("%s%s", filename, ext)
	newPath := filepath.Join("uploads", "avatars", filename)

	if err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm); err != nil {
		return "", apierror.FromErr(err)
	}

	if err := fileutils.SaveMedia(ctx, file, newPath); err != nil {
		return "", err
	}

	newUrl := "/uploads/avatars/" + filename
	switch role {
	case constants.ADMIN:
		if err := db.WithContext(ctx).Model(&Admin{}).
			Where("id = ?", userID).
			Update("avatar_url", newUrl).Error; err != nil {
			return "", err
		}

	case constants.CUSTOMER:
		if err := db.WithContext(ctx).Model(&Customer{}).
			Where("id = ?", userID).
			Update("avatar_url", newUrl).Error; err != nil {
			return "", err
		}
	}

	if oldAvatar != "" {
		oldPath := strings.TrimPrefix(oldAvatar, "/")
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			return "", apierror.FromErr(err)
		}
	}

	return newUrl, nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetPasswordReq) (res *ResetPasswordRes, err error) {
	db := s.AdminDB.DB

	switch req.Role {
	case constants.ADMIN:
		var admin Admin
		err = db.WithContext(ctx).
			Where("email = ?", req.Email).
			First(&admin).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("admin not found")
		} else if err != nil {
			return nil, err
		}
		res = &ResetPasswordRes{
			Email: admin.Email,
		}
		return res, nil

	case constants.CUSTOMER:
		var customer Customer
		err = db.WithContext(ctx).
			Where("email = ?", req.Email).
			First(&customer).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.CustomerNotFound()
		} else if err != nil {
			return nil, err
		}

		res = &ResetPasswordRes{
			Email: customer.Email,
		}
		return res, nil

	default:
		return nil, errors.New("invalid role")
	}
}

func (s *service) ResetPasswordSubmit(ctx context.Context, req ResetPasswordSubmitReq) (err error) {
	db := s.AdminDB.DB

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return apierror.FromErr(err)
	}

	switch req.Role {
	case constants.ADMIN:
		var admin Admin
		err = db.WithContext(ctx).
			Where("email = ?", req.Email).
			First(&admin).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else if err != nil {
			return err
		}
		admin.Password = hashedPassword
		if err = db.WithContext(ctx).Save(&admin).Error; err != nil {
			return err
		}

		return nil

	case constants.CUSTOMER:
		var customer Customer
		err = db.WithContext(ctx).
			Where("email = ?", req.Email).
			First(&customer).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else if err != nil {
			return err
		}

		customer.Password = hashedPassword
		if err = db.WithContext(ctx).Save(&customer).Error; err != nil {
			return err
		}

		return nil

	default:
		return nil
	}
}

func (s *service) GoogleAuth(ctx context.Context, input GoogleAuth) (*LoginRes, error) {
	db := s.CustomerDB.DB

	var existing Customer
	err := db.WithContext(ctx).
		Where("google_id = ?", input.GoogleID).
		First(&existing).Error

	var userID uuid.UUID

	if err == nil {
		userID = existing.ID
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.FromErr(err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {

		password := randomPassword()

		customer := Customer{
			Name:      input.Name,
			Email:     input.Email,
			Password:  password,
			AvatarUrl: input.Picture,
			GoogleID:  input.GoogleID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.WithContext(ctx).Create(&customer).Error; err != nil {
			return nil, apierror.FromErr(err)
		}

		userID = customer.ID
	}

	expirationTime := time.Now().Add(s.authConfig.JWT.ExpireIn)
	claims := &constants.JWTClaims{
		UserID: userID,
		Email:  input.Email,
		Role:   constants.CUSTOMER,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.authConfig.JWT.SecretKey))
	if err != nil {
		return nil, apierror.FromErr(err)
	}

	return &LoginRes{
		Role:    constants.CUSTOMER,
		Token:   tokenString,
		Expires: expirationTime,
	}, nil
}

func (s *service) DeleteAvatar(ctx context.Context) error {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return err
	}

	db, err := s.dbSelector.GetDBByRole(ctx)
	if err != nil {
		return err
	}

	userID := token.Claims.UserID
	role := token.Claims.Role

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var oldAvatar string

	switch role {
	case constants.ADMIN:
		var admin Admin
		if err := tx.WithContext(ctx).Where("id = ?", userID).First(&admin).Error; err != nil {
			tx.Rollback()
			return apierror.FromErr(err)
		}
		oldAvatar = admin.AvatarUrl

	case constants.CUSTOMER:
		var customer Customer
		if err := tx.WithContext(ctx).Where("id = ?", userID).First(&customer).Error; err != nil {
			tx.Rollback()
			return apierror.FromErr(err)
		}
		oldAvatar = customer.AvatarUrl
	}

	switch role {
	case constants.ADMIN:
		if err := tx.WithContext(ctx).Model(&Admin{}).
			Where("id = ?", userID).
			Update("avatar_url", nil).Error; err != nil {
			tx.Rollback()
			return err
		}

	case constants.CUSTOMER:
		if err := tx.WithContext(ctx).Model(&Customer{}).
			Where("id = ?", userID).
			Update("avatar_url", nil).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if oldAvatar != "" {
		oldPath := strings.TrimPrefix(oldAvatar, "/")
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			tx.Rollback()
			return apierror.FromErr(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
