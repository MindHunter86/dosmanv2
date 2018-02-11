package main

import "sync"
import "time"
import "encoding/json"
import "net/http"
import "github.com/gorilla/mux"

import "dosmanv2/system/config"

import "github.com/justinas/alice"

import "github.com/rs/zerolog"
import "github.com/rs/zerolog/hlog"


type vkApi struct {
	log *zerolog.Logger
	vkLogger *VKLogger

	httpSrv *http.Server
}

type vkApiResponse struct {
	Results interface{} `json:"results,omitempty"`
	Errors *vkApiResponseErr `json:"errors,omitempty"`
	Status string `json:"status"`
}

type vkApiResponseErr struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}


// vkApi internal API:
func (m *vkApi) construct(log *zerolog.Logger, cfg *config.SysConfig, vk *VKLogger) *vkApi {
	m.log = log
	m.vkLogger = vk

	var c = alice.New().Append(
		hlog.NewHandler(*m.log),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestHandler("request"),
		hlog.RefererHandler("referer"),
		hlog.UserAgentHandler("ua"))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).Msg("")
	}))

	var r = mux.NewRouter()
	r.Host(cfg.Vklogger.Http_api.Host).Subrouter()
	r.Schemes(cfg.Vklogger.Http_api.Schema)
	r.Headers("Content-Type", "application/json")

	r.HandleFunc("/", m.httpRootHandler).Methods("GET")
	r.HandleFunc("/v1/wall/{id:-?[0-9]+_?[0-9]+}", m.httpWallHandler).Methods("GET", "OPTIONS")

	m.httpSrv = &http.Server{
		Handler: c.Then(r),
		Addr: cfg.Vklogger.Http_api.Listen,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second}
	return m
}

func (m *vkApi) bootstrap(wg *sync.WaitGroup) {
	wg.Add(1)
	if e := m.httpSrv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		m.log.Error().Err(e).Msg("VkApi HTTP Serve Error!")}
	wg.Done()
}

// vkApi HTTP handlers:
func (m *vkApi) httpRootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK", "version": "1.0"}`))
}

func (m *vkApi) httpWallHandler(w http.ResponseWriter, r *http.Request) {
	var mv =  mux.Vars(r)

	wallId,e := m.vkLogger.vkGetWallPostWiget(mv["id"]); if e != nil {
		m.log.Warn().Err(e).Msg("Could not get requested hash!")
		m.respondJSON(w, http.StatusInternalServerError, "ERROR_WALL_HASH_IS_UNDEFINED", nil, &vkApiResponseErr{
			Name: "ERROR_WALL_HASH_IS_UNDEFINED",
			Desc: "Could not get requested hash: " + e.Error()})
		return
	}

	m.respondJSON(w, http.StatusOK, "OK", struct {
		Hash string `json:"hash"`
	}{wallId}, nil)
}

// vkapi HTTP responder:
func (m *vkApi) respondJSON(w http.ResponseWriter, status int, statusCode string, payload interface{}, errors *vkApiResponseErr) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(&vkApiResponse{
		Results: payload,
		Errors: errors,
		Status: statusCode})
}
