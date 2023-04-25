package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergio-abreu/note-taking-app-backend-golang/application/notes"
)

func (w WebServer) ScheduleReminder(c *gin.Context) {
	var r notes.ScheduleReminderRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := w.command.ScheduleReminder(c.Request.Context(), c.Param("userID"), c.Param("noteID"), r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}
