package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergio-abreu/note-taking-app-backend-golang/application/notes"
)

func (w WebServer) CreateNote(c *gin.Context) {
	var r notes.CreateNoteRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := w.app.CreateNote(c.Request.Context(), c.Param("userID"), r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}
