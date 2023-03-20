package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github/im-lauson/Short-Address/conf"
	"github/im-lauson/Short-Address/global"
	"github/im-lauson/Short-Address/middleware"
	"gopkg.in/validator.v2"

	"log"
	"net/http"
)

// ===============================================================================================
// = 定义App结构体，封装了mux路由的数据结构
type App struct {
	Router     *mux.Router
	Middleware *middleware.Middleware
	Config     *conf.Env
}

// ! 定义shortenReq请求结构体
type shortenReq struct {
	URl                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

// ! 定义shortLinkResp响应结构体
type shortLinkResp struct {
	ShortLink string `json:"short_link"`
}

// ===============================================================================================
// = 初始化App函数

func (a *App) Initialize(e *conf.Env) {
	// ! 设置日志格式 时间、文件名
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// ! 初始化Router实例
	a.Router = mux.NewRouter()
	a.Middleware = &middleware.Middleware{}
	a.Config = e
	a.InitializeRoutes()
}

//	===============================================================================================
// = 绑定路由和handler之间的关系,当匹配到相应的路由的时候，
// ！执行相应的函数createShortLink、getShortLinkInfo、redirect
// ！Alice 处理panic，从panic中恢复数据

func (a *App) InitializeRoutes() {
	m := alice.New(a.Middleware.LoggingHandler, a.Middleware.RecoverHandler)
	// ! 创建短地址

	a.Router.Handle("/api/shorten",
		m.ThenFunc(a.createShortLink)).Methods("POST")
	// ! 获得短地址的详细信息

	a.Router.Handle("/api/info",
		m.ThenFunc(a.getShortLinkInfo)).Methods("GET")
	// ! 重定向接口，指定地址格式，限定位数

	a.Router.Handle("/{shortLink:[a-zA-Z0-9]{1,11}}",
		m.ThenFunc(a.redirect)).Methods("GET")

	// ! 创建短地址
	//a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("Post")
	// ! 获得短地址的详细信息
	//a.Router.HandleFunc("/api/info", a.getShortLinkInfo).Methods("Get")
	// ! 重定向接口，指定地址格式，限定位数
	//a.Router.HandleFunc("/{shortLink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("Get")
}

func (a *App) createShortLink(writer http.ResponseWriter, request *http.Request) {

	var req shortenReq
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		respondWithError(writer, global.StatusError{http.StatusBadRequest,
			fmt.Errorf("validate parameters faileed %v", request.Body),
		})
		return
	}
	if err := validator.Validate(req); err != nil {
		respondWithError(writer, global.StatusError{http.StatusBadRequest,
			fmt.Errorf("parse parameters faileed %v", req),
		})

		return
	}
	defer request.Body.Close()
	s, err := a.Config.S.Shorten(req.URl, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(writer, err)
	} else {
		respondWithJSON(writer, http.StatusCreated, shortLinkResp{
			ShortLink: s},
		)
	}
}

func (a *App) getShortLinkInfo(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	s := vals.Get("shortLink")
	d, err := a.Config.S.ShortLinkInfo(s)
	if err != nil {
		respondWithError(writer, err)
	} else {
		respondWithJSON(writer, http.StatusOK, d)
	}
}

func (a *App) redirect(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	u, err := a.Config.S.UnShorten(vars["shortLink"])
	if err != nil {
		respondWithError(writer, err)
	} else {
		http.Redirect(writer, request, u, http.StatusTemporaryRedirect)
	}
}

// ===============================================================================================
// = 定义一下错误函数
func respondWithError(writer http.ResponseWriter, err error) {
	switch e := err.(type) {
	case global.Error:
		log.Println("Http %d - %s", e.Status(), e)
		respondWithJSON(writer, e.Status(), e.Error())
	default:
		respondWithJSON(writer, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(resp)
}

// ===============================================================================================
// = 启动监听和服务

func (a *App) Run(add string) {
	// ! 声明一个绑定关系，监听到地址，交由Router实例去处理
	log.Fatal(http.ListenAndServe(add, a.Router))
}
