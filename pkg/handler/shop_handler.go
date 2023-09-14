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
	res, err := h.services.Shop.Create(c, data)
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
	res, err := h.services.Shop.Update(c, data)
	if err != nil {
		responseWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccessful(c, res)
}

func (h *Handler) delete(c *gin.Context) {
	projectId, goodId := c.Query("projectId"), c.Query("id")
	if projectId == "" || goodId == "" {
		responseWithError(c, http.StatusBadRequest, ErrInvalidRequest.Error())
		return
	}

	var data models.Good

	data.ProjectId, _ = strconv.Atoi(projectId)
	data.Id, _ = strconv.Atoi(goodId)

	res, err := h.services.Shop.Delete(c, data)
	if err != nil {
		responseWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccessful(c, res)
}

func (h *Handler) list(c *gin.Context) {
	limit, offset := c.Query("limit"), c.Query("offset")

	var limitInt, offsetInt = 10, 0
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	if offset != "" {
		offsetInt, _ = strconv.Atoi(offset)
	}

	res, err := h.services.Shop.GetAll(c, limitInt, offsetInt)
	if err != nil {
		responseWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccessful(c, res)

}

func (h *Handler) reprioritiize(c *gin.Context) {
	projectId, goodId := c.Query("projectId"), c.Query("id")
	if projectId == "" || goodId == "" {
		responseWithError(c, http.StatusBadRequest, ErrInvalidRequest.Error())
		return
	}

	type Reprioritiize struct {
		NewPriority *int `json:"new_priority"`
		ProjectId   int
		GoodId      int
	}

	var data Reprioritiize
	if err := c.BindJSON(&data); err != nil {
		responseWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if data.NewPriority == nil {
		responseWithError(c, http.StatusBadRequest, ErrParseJSON.Error())
		return
	}

	data.ProjectId, _ = strconv.Atoi(projectId)
	data.GoodId, _ = strconv.Atoi(goodId)

	res, err := h.services.Shop.Reprioritiize(c, data.GoodId, data.ProjectId, *data.NewPriority)
	if err != nil {
		responseWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccessful(c, res)

}
