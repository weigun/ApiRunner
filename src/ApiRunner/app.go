// app.go

package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func bootstrap() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	// ---- Serve our controllers. ----

	// Prepare our repositories and services.
	db := datasource.GetDBInstance()
	cache := cache.NewCache()
	repo := repositories.NewStockRepository(db)
	stockService := services.NewStockService(repo, cache)
	uploadService := services.NewUploadService(repo, cache)
	cron := cronService.NewCronService()
	cron.StartService()

	// "/users" based mvc application.
	stocks := mvc.New(app.Party("/stocks"))
	// Add the basic authentication(admin:password) middleware
	// for the /users based requests.
	//	users.Router.Use(middleware.BasicAuth)
	// Bind the "userService" to the UserController's Service (interface) field.
	stocks.Register(stockService, uploadService)
	stocks.Handle(new(controllers.StocksController))

	stock := mvc.New(app.Party("/stock"))
	stock.Register(
		stockService,
		sessManager.Start,
	)
	stock.Handle(new(controllers.StockController))

	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8080"),
		iris.WithPostMaxMemory(MAXSIZE),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
