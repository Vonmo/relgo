/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package services

import (
	"context"
	"encoding/json"
	"github.com/elzor/relgo/config"
	"github.com/elzor/relgo/core"
	"github.com/elzor/relgo/lib"
	"github.com/elzor/relgo/log"
	"github.com/elzor/relgo/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Acounter = "acounter"

type ACounter struct {
	core.Service
	http *http.Server
}

//=================================================================================================
// service init
//=================================================================================================

func init() {
	go (&ACounter{
		Service: core.NewService(Acounter),
	}).start()
}

func (srv *ACounter) start() {
	core.WaitCore()
	srv.Ready = true
	srv.ShutdownFun = func(reason string) {
		log.Debug("acounter: soft shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.http.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
	core.Register(&srv.Service)
	log.Fatal(srv.server().ListenAndServe())
}

//=================================================================================================
// http server
//=================================================================================================

func (srv *ACounter) router() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/", srv.controllerIndex).Methods("GET")
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/increase/{name}", srv.controllerIncrease).Methods("GET")
	api.HandleFunc("/decrease/{name}", srv.controllerDecrease).Methods("GET")
	api.HandleFunc("/reset/{name}", srv.controllerReset).Methods("GET")
	api.HandleFunc("/get/{name}", srv.controllerGet).Methods("GET")
	return router
}

func (srv *ACounter) server() (server *http.Server) {
	// init cors
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	// http instance
	server = &http.Server{
		Addr:         core.System.Config.Services.Acounter.Socket.Host + ":" + strconv.Itoa(core.System.Config.Services.Acounter.Socket.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 233 * time.Second,
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(srv.router()),
	}
	srv.http = server
	return server
}

//=================================================================================================
// types
//=================================================================================================

type Response struct {
	Code      int   `json:"code"`
	Timestamp int64 `json:"ts"`
}

type ResponseOk struct {
	Response
}

type ResponseValue struct {
	Response
	Value   int    `json:"val"`
	Updated string `json:"updated"`
}

type ResponseError struct {
	Response
	Message string `json:"msg"`
}

//=================================================================================================
// controllers
//=================================================================================================

func (srv *ACounter) controllerIndex(w http.ResponseWriter, r *http.Request) {
	core.System.Metrics.Increase("http.index")
	views := int64(0)
	v := core.System.Metrics.Get("http.index")
	if v != nil {
		views = v.(int64)
	}
	content, err := lib.ParseTemplate(
		"templates/index.html",
		core.System.Boxes.Static.String("templates/index.html"),
		struct {
			Vsn   string
			Views int64
		}{
			Vsn:   config.Commit,
			Views: views,
		},
	)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Write(content)
	return
}

func (srv *ACounter) controllerIncrease(w http.ResponseWriter, r *http.Request) {
	core.System.Metrics.Increase("http.increase")

	vars := mux.Vars(r)
	name := strings.TrimSpace(vars["name"])

	c := &models.Counter{
		Name: name,
	}
	err := c.Increment()
	if err != nil {
		replyError(err, w)
		return
	}

	replyOk(w)
}

func (srv *ACounter) controllerDecrease(w http.ResponseWriter, r *http.Request) {
	core.System.Metrics.Increase("http.decrease")

	vars := mux.Vars(r)
	name := strings.TrimSpace(vars["name"])

	c := &models.Counter{
		Name: name,
	}
	err := c.Decrement()
	if err != nil {
		replyError(err, w)
		return
	}

	replyOk(w)
}

func (srv *ACounter) controllerReset(w http.ResponseWriter, r *http.Request) {
	core.System.Metrics.Increase("http.reset")

	var err error
	vars := mux.Vars(r)
	name := strings.TrimSpace(vars["name"])

	c := &models.Counter{
		Name: name,
	}
	err = c.Reset()
	if err != nil {
		replyError(err, w)
		return
	}

	replyOk(w)
}

func (srv *ACounter) controllerGet(w http.ResponseWriter, r *http.Request) {
	core.System.Metrics.Increase("http.get")

	vars := mux.Vars(r)
	name := strings.TrimSpace(vars["name"])

	c := &models.Counter{
		Name: name,
	}
	err := c.Get()
	if err != nil && err.Error() != "sql: no rows in result set" {
		replyError(err, w)
		return
	}

	replyVal(c, w)
}

//=================================================================================================
// helpers
//=================================================================================================

func replyVal(c *models.Counter, w http.ResponseWriter) {
	resp, err := encodeResult(ResponseValue{
		Response: Response{
			Code:      core.SUCCESS,
			Timestamp: time.Now().UnixNano(),
		},
		Value:   c.Value,
		Updated: c.Updated,
	})
	if err != nil {
		log.Errorf("json: %s", err.Error())
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
	w.Write(resp)
}

func replyOk(w http.ResponseWriter) {
	resp, err := encodeResult(ResponseOk{
		Response: Response{
			Code:      core.SUCCESS,
			Timestamp: time.Now().UnixNano(),
		},
	})
	if err != nil {
		log.Errorf("json: %s", err.Error())
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
	w.Write(resp)
}

func replyError(err error, w http.ResponseWriter) {
	resp, err := encodeResult(ResponseError{
		Response: Response{
			Code:      core.ERROR_INTERNAL,
			Timestamp: time.Now().UnixNano(),
		},
		Message: err.Error(),
	})
	if err != nil {
		log.Errorf("json: %s", err.Error())
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
	w.Write(resp)
}

func encodeResult(result interface{}) (bin []byte, err error) {
	bin, err = json.Marshal(result)
	return bin, err
}
