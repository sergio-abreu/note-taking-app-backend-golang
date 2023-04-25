package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergio-abreu/note-taking-app-backend-golang/application/notes"
)

func (w WebServer) RescheduleReminder(c *gin.Context) {
	var r notes.RescheduleReminderRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := w.app.RescheduleReminder(c.Request.Context(), c.Param("userID"), c.Param("noteID"), c.Param("reminderID"), r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
