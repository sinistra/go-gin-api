package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"rits/carriage/api/models"
	"rits/carriage/api/repository"
	"rits/carriage/api/utils"
)

type ServiceController struct{}

var Services []models.Service

func (s ServiceController) GetServices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error

		var services []models.Service
		serviceRepo := repository.ServiceRepository{}
		services, err := serviceRepo.GetActiveServices()

		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.RespondWithJSON(w, services)
	}
}

func (s ServiceController) GetService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var service models.Service
		var error models.Error

		params := mux.Vars(r)

		serviceRepo := repository.ServiceRepository{}

		id, _ := params["id"]

		service, err := serviceRepo.GetService(id)

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
		utils.RespondWithJSON(w, service)
	}
}

func (s ServiceController) AddService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var service models.Service
		var error models.Error

		json.NewDecoder(r.Body).Decode(&service)

		if service.Name == "" {
			error.Message = "Enter missing fields."
			utils.RespondWithError(w, http.StatusBadRequest, error) //400
			return
		}
		serviceRepo := repository.ServiceRepository{}
		serviceID, err := serviceRepo.AddService(service)

		if err != nil {
			error.Message = "Server error"
			utils.RespondWithError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.RespondWithJSON(w, serviceID)
	}
}

func (s ServiceController) UpdateService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var service models.Service
		var error models.Error

		json.NewDecoder(r.Body).Decode(&service)

		if service.ID == "" || service.Name == "" {
			error.Message = "missing fields."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		serviceRepo := repository.ServiceRepository{}
		err := serviceRepo.UpdateService(service)

		if err != nil {
			if err.Error() == "not found" {
				error.Message = fmt.Sprintf("%s not found.", service.ID.Hex())
				utils.RespondWithError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server error"
				utils.RespondWithError(w, http.StatusInternalServerError, error) //500
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.RespondWithJSON(w, fmt.Sprintf("%s updated", service.ID.Hex()))
	}
}

func (s ServiceController) RemoveService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)
		serviceRepo := repository.ServiceRepository{}
		id := params["id"]

		err := serviceRepo.RemoveService(id)

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
