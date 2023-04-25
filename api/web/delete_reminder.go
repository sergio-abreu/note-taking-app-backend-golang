package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w WebServer) DeleteReminder(c *gin.Context) {
	err := w.command.DeleteReminder(c.Request.Context(), c.Param("userID"), c.Param("noteID"), c.Param("reminderID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
