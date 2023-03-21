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

func (m Middleware) Cors(f http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		w.Header().Add("content-type", "application/json;charset=UTF-8")
		w.Header().Add("Access-Control-Max-Age", "259200")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		f.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
