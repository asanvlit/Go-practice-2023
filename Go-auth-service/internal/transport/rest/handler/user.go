package handler

import (
	"Golang-practice-2023/internal/domain/logger"
	"Golang-practice-2023/internal/domain/user"
	myHttp "Golang-practice-2023/pkg/utils/http"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

const contentType = "application/json"

type UserHandler struct {
	service user.Service
	logger  logger.Logger
}

func New(service user.Service, logger logger.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

func (h *UserHandler) InitRoutes(router *mux.Router) {
	router.HandleFunc("/user", h.Create).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}", h.Update).Methods(http.MethodPut)
	router.HandleFunc("/user/{id}", h.GetById).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}", h.Delete).Methods(http.MethodDelete)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := myHttp.ValidateRequestFormat(r, contentType)
	if err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: fmt.Sprintf("Error: Content-Type header is not %s", contentType)},
			w, contentType, http.StatusBadRequest)
		return
	}

	var u user.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Info("Error: Failed to decode request body. " + err.Error())
		_ = myHttp.BuildResponse(Error{Code: 400, Message: "Error: failed to decode request body"},
			w, contentType, http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &u); err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: "Error: failed to create user"}, w, contentType, http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Warning("Failed to publish info about created user")
		return
	}

	if err := myHttp.BuildResponse(&u, w, contentType, http.StatusCreated); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct: %v ", u))
	}

	return
}

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: fmt.Sprintf("The id %s does not have the Uuid format", id)},
			w, contentType, http.StatusBadRequest)
		return
	}

	u, err := h.service.GetById(r.Context(), parsedUuid)
	if err != nil {
		HandleError(w, err)
		return
	}

	if err := myHttp.BuildResponse(&u, w, contentType, http.StatusOK); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct: %v ", u))
	}
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	err := myHttp.ValidateRequestFormat(r, contentType)
	if err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: fmt.Sprintf("Error: Content-Type header is not %s", contentType)},
			w, contentType, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	id := vars["id"]
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: fmt.Sprintf("The id %s does not have the Uuid format", id)},
			w, contentType, http.StatusBadRequest)
		return
	}

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: "Error: failed to decode request body"},
			w, contentType, http.StatusBadRequest)
		return
	}

	u.ID = parsedUuid
	if err := h.service.Update(r.Context(), &u); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = myHttp.BuildResponse(Error{Code: 404, Message: fmt.Sprintf("Error: user with id %s does not exist", parsedUuid)},
				w, contentType, http.StatusNotFound)
			return
		}
		_ = myHttp.BuildResponse(Error{Code: 400, Message: "Error: failed to update user"}, w, contentType, http.StatusBadRequest)
		return
	}

	if err := myHttp.BuildResponse(&u, w, contentType, http.StatusAccepted); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct: %v ", u))
	}
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		_ = myHttp.BuildResponse(Error{Code: 400, Message: fmt.Sprintf("The id %s does not have the Uuid format", parsedUuid)},
			w, contentType, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), parsedUuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = myHttp.BuildResponse(Error{Code: 404, Message: fmt.Sprintf("Error: user with id %s does not exist", parsedUuid)},
				w, contentType, http.StatusNotFound)
			return
		}
		_ = myHttp.BuildResponse(Error{Code: 400, Message: "Error: failed to delete user"}, w, contentType, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
