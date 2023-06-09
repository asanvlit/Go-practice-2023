package handler

import (
	"Golang-practice-2023/internal/domain/apperrors"
	"Golang-practice-2023/internal/domain/logger"
	"Golang-practice-2023/internal/domain/user"
	myHttp "Golang-practice-2023/pkg/utils/http"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	router.HandleFunc("/user", h.GetWithOffsetAndLimit).Methods(http.MethodGet).Queries("offset", "{offset}", "limit", "{limit}")
	router.HandleFunc("/user", h.GetRegisteredLaterThen).Methods(http.MethodGet).Queries("date", "{date}")
	router.HandleFunc("/user/{id}", h.Delete).Methods(http.MethodDelete)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := myHttp.ValidateRequestFormat(r, contentType)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidRequestFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	var u user.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		err := HandleError(w, apperrors.ErrInvalidRequestBody)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := h.service.Create(r.Context(), &u); err != nil {
		err := HandleError(w, err)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := myHttp.WriteResponse(&u, w, contentType, http.StatusCreated); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct: %v ", u))
	}
}

func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidIdFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	u, err := h.service.GetById(r.Context(), parsedUuid)
	if err != nil {
		err := HandleError(w, err)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := myHttp.WriteResponse(&u, w, contentType, http.StatusOK); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct: %v ", u))
	}
}

func (h *UserHandler) GetWithOffsetAndLimit(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidOffsetFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidLimitFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	var users *[]user.User
	if users, err = h.service.GetWithOffsetAndLimit(r.Context(), offset, limit); err != nil {
		err := HandleError(w, err)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := myHttp.WriteResponse(&users, w, contentType, http.StatusOK); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct"))
	}
}

func (h *UserHandler) GetRegisteredLaterThen(w http.ResponseWriter, r *http.Request) {
	registerDateString := r.URL.Query().Get("date")
	limitStr := r.URL.Query().Get("limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidLimitFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidOffsetFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	var users *[]user.User
	if users, err = h.service.GetRegisteredLaterThenWithLimit(r.Context(), registerDateString, limit); err != nil {
		err := HandleError(w, err)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := myHttp.WriteResponse(&users, w, contentType, http.StatusOK); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct"))
	}
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	err := myHttp.ValidateRequestFormat(r, contentType)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidRequestFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	vars := mux.Vars(r)

	id := vars["id"]
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidIdFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		err := HandleError(w, apperrors.ErrInvalidRequestBody)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	u.ID = parsedUuid
	if err := h.service.Update(r.Context(), &u); err != nil {
		err := HandleError(w, err)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := myHttp.WriteResponse(&u, w, contentType, http.StatusAccepted); err != nil {
		h.logger.Warning(fmt.Sprintf("Error: unable to marshal order struct: %v ", u))
	}
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		err := HandleError(w, apperrors.ErrInvalidIdFormat)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	if err := h.service.Delete(r.Context(), parsedUuid); err != nil {
		err := HandleError(w, err)
		if err != nil {
			h.logger.Warning(fmt.Sprintf("Failed to write response: %s", err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
