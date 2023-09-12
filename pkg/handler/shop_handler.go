package handler

import (
	"net/http"
	"strconv"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) create(c *gin.Context) {
	var data models.Good

	if err := c.BindJSON(&data); err != nil {
		responseWithError(c, http.StatusBadRequest, ErrParseJSON.Error())
		return
	}
	if data.Name == "" {
		responseWithError(c, http.StatusBadRequest, ErrInvalidRequest.Error())
		return
	}
	projectId := c.Query("projectId")
	if projectId == "" {
		responseWithError(c, http.StatusBadRequest, ErrInvalidRequest.Error())
		return
	}

	data.ProjectId, _ = strconv.Atoi(projectId)
	res, err := h.services.Shop.Create(data)
	if err != nil {
		responseWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	responseSuccessful(c, res)
}

func (h *Handler) update(c *gin.Context) {
	var data models.Good

	if err := c.BindJSON(&data); err != nil {
		responseWithError(c, http.StatusBadRequest, ErrParseJSON.Error())
		return
	}
	if data.Name == "" {
		responseWithError(c, http.StatusBadRequest, ErrInvalidRequest.Error())
		return
	}

	projectId, goodId := c.Query("projectId"), c.Query("id")
	if projectId == "" || goodId == "" {
		responseWithError(c, http.StatusBadRequest, ErrInvalidRequest.Error())
		return
	}

	data.ProjectId, _ = strconv.Atoi(projectId)
	data.Id, _ = strconv.Atoi(goodId)
	res, err := h.services.Shop.Update(data)
	if err != nil {
		responseWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccessful(c, res)
}

func (h *Handler) delete(c *gin.Context) {}

func (h *Handler) list(c *gin.Context) {}

func (h *Handler) reprioritiize(c *gin.Context) {}