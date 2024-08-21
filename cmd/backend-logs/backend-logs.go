package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/logging"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("starting BACKEND service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	r := mux.NewRouter().StrictSlash(true)

	//criteriaRepo := repos.NewSQLCriteriaRepository(cfg)
	//marketsRepo := repos.NewSQLMarketsRepository(cfg)
	//adsRepo := repos.NewAdsRepository(cfg)
	logsRepo := logging.NewLogsRepository(cfg)

	//chartsRepo := repos.NewChartsRepository(cfg)
	//chartsRepo.GetAdsPricesByStep(5000)

	//cleanupPrices(adsRepo)

	r.HandleFunc("/session/{id}", opt(logsRepo)).Methods("OPTIONS")
	r.HandleFunc("/session/{id}", deleteSession(logsRepo)).Methods("DELETE")
	r.HandleFunc("/sessions", getSessions(logsRepo)).Methods("GET")
	r.HandleFunc("/session/{id}", getSession(logsRepo)).Methods("GET")
	r.HandleFunc("/pagelogsforcriterialog/{id}", getPageLogsForCriteriaLog(logsRepo)).Methods("GET")
	//httpPort := cfg.GetString(amconfig.BackendServiceHTTPPort)
	httpPort := cfg.GetString(amconfig.AppBackendLogsPort)
	log.Printf("HTTP listening on port %s\n", httpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), r)
	errorshandler.HandleErr(err)

}

func opt(repository *logging.LogsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type optRes struct {
			res string
		}
		writeJSONResponse(optRes{res: "done"}, w)
	}
}

func deleteSession(repository *logging.LogsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getIDAsNumber(r)
		if err != nil {
			writeError(err, w)
			return
		}

		err = repository.DeleteSession(*sessionID)
		if err != nil {
			writeError(err, w)
			return
		}
		type Res struct {
			status bool
		}
		rr := Res{status: true}
		writeJSONResponse(rr, w)
	}
}

func getSessions(repository *logging.LogsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions, err := repository.GetSessions()
		if err != nil {
			panic(err)
		}
		writeJSONResponse(sessions, w)
	}
}

func getSession(repository *logging.LogsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID, err := getIDAsNumber(r)
		if err != nil {
			writeError(err, w)
			return
		}

		sessions, err := repository.GetSession(*sessionID)
		if err != nil {
			writeError(err, w)
			return
		}
		writeJSONResponse(sessions, w)
	}
}

func getPageLogsForCriteriaLog(repository *logging.LogsRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		criteriaLogID, err := getIDAsNumber(r)
		if err != nil {
			writeError(err, w)
			return
		}

		pageLogs, err := repository.GetPageLogsForCriteriaLogs(*criteriaLogID)
		if err != nil {
			writeError(err, w)
			return
		}
		writeJSONResponse(pageLogs, w)
	}
}

func writeJSONResponse(response any, w http.ResponseWriter) {

	w.Header().Set("Access-Control-Allow-Headers:", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

	res, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	//w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}

func writeError(err error, w http.ResponseWriter) {
	errStr := err.Error()
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(errStr))
}

func getIDAsNumber(r *http.Request) (*uint, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}
	id, err := strconv.Atoi(idStr)
	uintID := uint(id)
	if err != nil {
		return nil, err
	}
	return &uintID, nil
}
