package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
)

type S struct {
	ReminderId string `json:"reminder_id"`
	NoteId     string `json:"note_id"`
	UserId     string `json:"user_id"`
}

func main() {
	r := gin.Default()
	r.POST("/v1/webhooks/reminders/", func(c *gin.Context) {
		var s S
		err := c.BindJSON(&s)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println(s)

		c.JSON(http.StatusOK, s)
		return
	})

	if err := r.Run(":80"); err != nil {
		log.Fatalln(err.Error())
	}
}
