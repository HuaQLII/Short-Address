package conf

import (
	"github/im-lauson/Short-Address/middleware"
	"github/im-lauson/Short-Address/service/dto"
	"log"
	"os"
	"strconv"
)

type Env struct {
	S dto.Storage
}

func GetEnv() *Env {
	addr := os.Getenv("App_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	passwd := os.Getenv("App_REDIS_PASSWD")
	if passwd == "" {
		passwd = ""
	}
	dbS := os.Getenv("App_REDIS_DB")
	if dbS == "" {
		dbS = "0"
	}
	db, err := strconv.Atoi(dbS)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connect to redis (addr:%s passwd :%s db:%d)", addr, passwd, db)

	r := middleware.NewRedisCli("127.0.0.1:63791", "", 1)

	return &Env{S: r}
}
