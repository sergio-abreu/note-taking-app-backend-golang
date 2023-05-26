package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w WebServer) GetInProgressNotes(c *gin.Context) {
	response, err := w.query.GetInProgressNotes(c.Request.Context(), c.Param("userID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
