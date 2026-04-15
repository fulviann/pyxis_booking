package routes

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/devanadindra/signlink-mobile/back-end/database"
	"github.com/devanadindra/signlink-mobile/back-end/utils/logger"
)

type Dependency struct {
	handler    *gin.Engine
	AdminDB    *database.AdminDB
	CustomerDB *database.CustomerDB
}

func (d *Dependency) Close() {
	ctx := context.Background()

	if d.CustomerDB != nil {
		if db, err := d.CustomerDB.DB.DB(); err == nil {
			_ = db.Close()
		} else {
			logger.Error(ctx, "Error closing customer DB: %v", err)
		}
	}

	if d.AdminDB != nil {
		if db, err := d.AdminDB.DB.DB(); err == nil {
			_ = db.Close()
		} else {
			logger.Error(ctx, "Error closing admin DB: %v", err)
		}
	}
}

func (d *Dependency) GetHandler() *gin.Engine {
	return d.handler
}

func (d *Dependency) GetCustomerSQLDB() *sql.DB {
	db, _ := d.CustomerDB.DB.DB()
	return db
}

func (d *Dependency) GetAdminSQLDB() *sql.DB {
	db, _ := d.AdminDB.DB.DB()
	return db
}
