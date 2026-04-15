package dbselector

import (
	"context"
	"net/http"

	"github.com/fulviann/pyxis_booking/back-end/database"
	apierror "github.com/fulviann/pyxis_booking/back-end/utils/api-error"
	"github.com/fulviann/pyxis_booking/back-end/utils/constants"
	contextUtil "github.com/fulviann/pyxis_booking/back-end/utils/context"
	"gorm.io/gorm"
)

type DBService struct {
	AdminDB    *database.AdminDB
	CustomerDB *database.CustomerDB
}

func NewDBService(adminDB *database.AdminDB, customerDB *database.CustomerDB) *DBService {
	return &DBService{
		AdminDB:    adminDB,
		CustomerDB: customerDB,
	}
}

// helper pilih DB sesuai role
func (s *DBService) GetDBByRole(ctx context.Context) (*gorm.DB, error) {
	token, err := contextUtil.GetTokenClaims(ctx)
	if err != nil {
		return nil, apierror.NewWarn(http.StatusUnauthorized)
	}

	switch token.Claims.Role {
	case constants.ADMIN:
		return s.AdminDB.DB, nil
	case constants.CUSTOMER:
		return s.CustomerDB.DB, nil
	default:
		return nil, apierror.NewWarn(http.StatusUnauthorized)
	}
}
