package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w WebServer) MarkNoteAsComplete(c *gin.Context) {
	response, err := w.app.MarkNoteAsComplete(c.Request.Context(), c.Param("userID"), c.Param("noteID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
