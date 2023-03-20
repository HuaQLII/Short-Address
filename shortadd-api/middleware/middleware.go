package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
}

// LoggingHandler ===============================================================================================
// = 记录请求需要的时间
// ! http.handler可以理解成一个适配器
func (m Middleware) LoggingHandler(next http.Handler) http.Handler {
	// ! 实现匿名函数
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

// RecoverHandler ===============================================================================================
// = 把函数从 panic 恢复出来
func (m Middleware) RecoverHandler(next http.Handler) http.Handler {
	// ! 实现匿名函数
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover from panicL: %v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
