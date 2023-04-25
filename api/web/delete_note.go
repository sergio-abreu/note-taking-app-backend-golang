package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w WebServer) DeleteNote(c *gin.Context) {
	err := w.command.DeleteNote(c.Request.Context(), c.Param("userID"), c.Param("noteID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
