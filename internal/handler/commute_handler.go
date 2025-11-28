package handler

import (
	"fmt"
	"net/http"

	"github.com/ekastn/commute-analyzer/internal/dto"
	"github.com/ekastn/commute-analyzer/internal/response"
	"github.com/ekastn/commute-analyzer/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommuteHandler struct {
	service *service.CommuteService
}

func NewCommuteHandler(s *service.CommuteService) *CommuteHandler {
	return &CommuteHandler{service: s}
}

func (h *CommuteHandler) CreateCommute(c *gin.Context) {
	var req dto.CreateCommuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to bind JSON: %s", err.Error()))
		return
	}

	commute, err := h.service.CreateCommute(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create commute: %s", err.Error()))
		return
	}

	response.Success(c, http.StatusCreated, commute)
}

func (h *CommuteHandler) ListCommutes(c *gin.Context) {
	deviceID := c.Query("device_id")
	if deviceID == "" {
		response.Error(c, http.StatusBadRequest, "Missing device_id")
		return
	}

	list, err := h.service.ListCommutes(c.Request.Context(), deviceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to list commutes: %s", err.Error()))
		return
	}

	response.Success(c, http.StatusOK, list)
}

func (h *CommuteHandler) UpdateCommute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateCommuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to bind JSON: %s", err.Error()))
		return
	}

	commute, err := h.service.UpdateCommute(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update commute: %s", err.Error()))
		return
	}

	response.Success(c, http.StatusOK, commute)
}

func (h *CommuteHandler) DeleteCommute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.DeleteCommute(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete commute: %s", err.Error()))
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}
