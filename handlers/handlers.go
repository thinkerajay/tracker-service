package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerajay/tracker-service/queue_service"
	"github.com/thinkerajay/tracker-service/types"
	"log"
	"net/http"
	"time"
)

func HomePageHandler(ctx *gin.Context) {
	userId, err := ctx.Cookie("_tsuid")
	if err != nil {
		ctx.SetCookie("_tsuid", "random_id", 9000, "/", "heroku.com", true, true)
	}
	log.Println(userId)
	ctx.HTML(200, "index.tmpl.html", nil)

}

func AwesomePageHandler(ctx *gin.Context) {
	userId, err := ctx.Cookie("_tsuid")
	if err != nil {
		userId = "guest-user"
	}

	event := &types.Event{
		Type:      "view",
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
		PageUrl:   ctx.Request.URL.Path,
		User: struct {
			Id string
		}{Id: userId},
	}

	err = queue_service.Enqueue(event)
	if err != nil {
		log.Println(err)
	}
	// redundant
	ctx.Writer.WriteHeader(http.StatusOK)

}

func ViewsHandler(ctx *gin.Context) {

}
