//go:build wireinject
// +build wireinject

package wireinject

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/google/subcommands"
	"github.com/google/wire"

	"github.com/fulviann/pyxis_booking/back-end/database"
	"github.com/fulviann/pyxis_booking/back-end/domains/user"
	"github.com/fulviann/pyxis_booking/back-end/middlewares"
	"github.com/fulviann/pyxis_booking/back-end/routes"
	"github.com/fulviann/pyxis_booking/back-end/utils/config"
	"github.com/fulviann/pyxis_booking/back-end/utils/dbselector"
)

var dbSet = wire.NewSet(
	database.NewDBCustomer,
	database.NewDBAdmin,
)

var dbSelectorSet = wire.NewSet(
	dbselector.NewDBService,
)

var userSet = wire.NewSet(
	user.NewService,
	user.NewHandler,
)

func initializeDependency(config *config.Config) (*routes.Dependency, error) {

	wire.Build(
		dbSet,
		dbSelectorSet,
		validator.New,
		middlewares.NewMiddlewares,
		routes.NewDependency,
		userSet,
	)

	return nil, nil
}
