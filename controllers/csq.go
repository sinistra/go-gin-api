package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"rits/carriage/api/models"
	"rits/carriage/api/repository"
	"rits/carriage/api/utils"
	"time"
)

type CsqController struct{}

var Csqs []models.Csq

func (c CsqController) GetCsqs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error

		var csqs []models.Csq
		csqRepo := repository.CsqRepository{}
		csqs, err := csqRepo.GetActiveCsqs()

		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.RespondWithJSON(w, csqs)
	}
}

func (c CsqController) GetCsq() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csq models.Csq
		var error models.Error

		params := mux.Vars(r)

		csqRepo := repository.CsqRepository{}

		id, _ := params["id"]

		csq, err := csqRepo.GetCsq(id)

		if err != nil {
			if err.Error() == "not found" {
				error.Message = fmt.Sprintf("%s not found.", id)
				utils.RespondWithError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server error"
				utils.RespondWithError(w, http.StatusInternalServerError, error)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		utils.RespondWithJSON(w, csq)
	}
}

func (c CsqController) AddCsq() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csq models.Csq
		var error models.Error

		json.NewDecoder(r.Body).Decode(&csq)

		if csq.CustomerName == "" {
			error.Message = "Enter missing fields."
			utils.RespondWithError(w, http.StatusBadRequest, error) //400
			return
		}
		csq.LastUpdate = time.Now()
		csq.RequestDate = time.Now()

		csqRepo := repository.CsqRepository{}
		csqID, err := csqRepo.AddCsq(csq)

		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.RespondWithJSON(w, csqID)
	}
}

func (c CsqController) UpdateCsq() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var csq models.Csq
		var error models.Error

		json.NewDecoder(r.Body).Decode(&csq)

		if csq.ID == "" || csq.CustomerName == "" {
			error.Message = "missing fields."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		csq.LastUpdate = time.Now()

		csqRepo := repository.CsqRepository{}
		err := csqRepo.UpdateCsq(csq)

		if err != nil {
			if err.Error() == "not found" {
				error.Message = fmt.Sprintf("%s not found.", csq.ID.Hex())
				utils.RespondWithError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server error"
				utils.RespondWithError(w, http.StatusInternalServerError, error) //500
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.RespondWithJSON(w, fmt.Sprintf("%s updated", csq.ID.Hex()))
	}
}

func (c CsqController) RemoveCsq() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)
		csqRepo := repository.CsqRepository{}
		id := params["id"]

		err := csqRepo.RemoveCsq(id)

		if err != nil {
			if err.Error() == "not found" {
				error.Message = fmt.Sprintf("%s not found.", id)
				utils.RespondWithError(w, http.StatusNotFound, error) //404
				return
			} else {
				error.Message = "Server error."
				utils.RespondWithError(w, http.StatusInternalServerError, error) //500
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.RespondWithJSON(w, fmt.Sprintf("%s removed.", id))
	}
}
