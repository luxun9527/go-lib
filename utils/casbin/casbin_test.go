package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func TestBase(t *testing.T) {
	e, err := casbin.NewEnforcer("./model.base.conf", "./policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "dajun", "data1", "read")
	check(e, "lizi", "data2", "write")
	check(e, "dajun", "data1", "write")
	check(e, "dajun", "data2", "read")
}

func TestRbac(t *testing.T) {
	e, err := casbin.NewEnforcer("./model.rbac.conf", "./policy.rbac.csv")

	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "dajun", "data", "read")
	check(e, "dajun", "data", "write")
	check(e, "lizi", "data", "read")
	check(e, "lizi", "data", "write")

}
func TestMysqlRbac(t *testing.T) {
	const modelText = `
    [role_definition]
	g = _, _
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	
	[policy_effect]
	e = some(where (p.eft == allow))
`

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color

		},
	)

	dsn := "root:root@tcp(192.168.254.99:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlBb, err := db.DB()
	sqlBb.SetMaxIdleConns(100)
	sqlBb.SetMaxOpenConns(50)
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalln(err)
	}
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		log.Fatalln(err)
	}
	enforcer, _ := casbin.NewSyncedEnforcer(m, adapter)
	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalln(err)
	}
	if _, err := enforcer.AddPolicy("admin", "/getUserList", "GET"); err != nil {
		log.Fatalln(err)
	}
	if _, err := enforcer.AddPolicy("superAdmin", "/getAdminList", "GET"); err != nil {
		log.Fatalln(err)
	}
	if _, err := enforcer.AddRoleForUser("zhangsan", "admin"); err != nil {
		log.Fatalln(err)
	}
	if _, err := enforcer.AddRoleForUser("zhangsan", "superAdmin"); err != nil {
		log.Fatalln(err)
	}
	result, err := enforcer.Enforce("zhangsan", "/getUserList", "GET")
	if err != nil {
		log.Fatalln(err)
	}
	//role, err := enforcer.GetUsersForRole("zhangsan")

	log.Println(result)

}
