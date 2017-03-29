package api

import (
	"net/http"

	"github.com/cjduffett/stork/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type APIController struct {
	session *mgo.Session
	dbname  string
}

func NewAPIController(session *mgo.Session, config *config.StorkConfig) *APIController {
	return &APIController{
		session: session,
		dbname:  config.DatabaseName,
	}
}

func (a *APIController) APIRoot(c *gin.Context) {
	c.String(http.StatusOK, "Stork API root")
}
