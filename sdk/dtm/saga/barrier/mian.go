package main

import (
	"database/sql"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:123456@tcp(192.168.2.99:3306)/dtm?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

const qsBusiAPI = "/api/busi_start"
const qsBusiPort = 8082

var qsBusi = fmt.Sprintf("http://192.168.2.218:%d%s", qsBusiPort, qsBusiAPI)

func main() {
	if err := initDB(); err != nil {
		log.Println("err", err)
		return
	}
	QsStartSvr()
	_ = QsFireRequest()
	time.Sleep(3 * time.Minute)
}

// QsStartSvr quick start: start server
func QsStartSvr() {
	app := gin.New()
	log.Printf("quick start examples listening at %d", qsBusiPort)
	qsAddRoute(app)
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", qsBusiPort))
	}()
	time.Sleep(100 * time.Millisecond)
}

func qsAddRoute(app *gin.Engine) {
	app.POST(qsBusiAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn")
		//if true {
		//	c.JSON(409, "")
		//	return
		//}
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			c.JSON(409, "")
			return
		}

		if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec("update user_account set balance = balance + ? where user_id = ?", 10, 2); err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(409, "")
			return
		}

		//c.JSON(500, "")
		c.JSON(200, "") // Status 409 for Failure. Won't be retried
	})
	app.POST(qsBusiAPI+"/TransInCompensate", func(c *gin.Context) {
		log.Printf("TransInCompensate")
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			c.JSON(409, "")
			return
		}
		if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec("update user_account set balance = balance - ? where user_id = ?", 10, 2); err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(409, "")
			return
		}
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOut", func(c *gin.Context) {
		log.Printf("TransOut")
		if true {
			c.JSON(409, "")
			return
		}
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			c.JSON(409, "")
			return
		}
		if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec("update user_account set balance = balance - ? where user_id = ?", 10, 1); err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(409, "")
			return
		}
		c.JSON(200, "")
	})
	app.POST(qsBusiAPI+"/TransOutCompensate", func(c *gin.Context) {
		log.Printf("TransOutCompensate")
		barrier, err := dtmcli.BarrierFromQuery(c.Request.URL.Query())
		if err != nil {
			c.JSON(409, "")
			return
		}
		if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
			if _, err := tx.Exec("update user_account set balance = balance + ? where user_id = ?", 10, 1); err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(409, "")
			return
		}
		log.Println("TransOutCompensate finish")
		//c.JSON(409, "")
		return
	})
}

const dtmServer = "http://192.168.2.99:36789/api/dtmsvr"

// QsFireRequest quick start: fire request
func QsFireRequest() string {
	req := &gin.H{"amount": 30} // the payload of requests
	gid := shortuuid.New()
	// DtmServer is the address of dtm
	saga := dtmcli.NewSaga(dtmServer, gid).
		// add a branch transaction，action url is: qsBusi+"/TransOut"， compensate url: qsBusi+"/TransOutCompensate"
		Add(qsBusi+"/TransOut", qsBusi+"/TransOutCompensate", req).
		// add a branch transaction，action url is: qsBusi+"/TransIn"， compensate url: qsBusi+"/TransInCompensate"
		Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// submit saga global transaction，dtm will finish all action and compensation
	err := saga.Submit()

	if err != nil {
		panic(err)
	}
	log.Printf("transaction: %s submitted", saga.Gid)
	return saga.Gid
}
