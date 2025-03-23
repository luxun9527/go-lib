package svc

import (
	"go-lib/frame/go-zero/apidemo/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold:             time.Second, // Slow SQL threshold
	//		LogLevel:                  logger.Info, // Log level
	//		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	//		Colorful:                  true,        // Disable color
	//	},
	//)
	//
	//dsn := "root:root@tcp(192.168.2.159:3307)/test?charset=utf8mb4&parseTime=True&loc=Local"
	//var err error
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	SkipDefaultTransaction: true,
	//	Logger:                 newLogger,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//sqlBb, err := db.DB()
	//sqlBb.SetMaxOpenConns(100)
	//sqlBb.SetMaxIdleConns(10)
	//sqlBb.SetConnMaxIdleTime(time.Hour * 5)
	//if err := db.Use(otelgorm.NewPlugin(
	//	otelgorm.WithDBName("test"),
	//	otelgorm.WithTracerProvider(otel.GetTracerProvider()),
	//)); err != nil {
	//	log.Fatal(err)
	//}
	return &ServiceContext{
		Config: c,
	}
}
