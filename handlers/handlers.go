package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerajay/tracker-service/db_service"
	"github.com/thinkerajay/tracker-service/queue_service"
	"github.com/thinkerajay/tracker-service/types"
	"log"
	"net/http"
	"time"
)

func HomePageHandler(ctx *gin.Context) {
	userId, err := ctx.Cookie("_tsuid")
	if err != nil {
		ctx.SetCookie("_tsuid", "random_id", 9000, "/", "localhost", true, true)
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
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		PageUrl:   ctx.Request.URL.Path,
		UserId:    userId,
	}

	err = queue_service.Enqueue(event)
	if err != nil {
		log.Println(err)
	}
	// redundant
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.HTML(200, "awesome_page.tmpl.html", nil)

}

func ViewsHandler(ctx *gin.Context) {
	startDateString := ctx.Query("startDate")

	var startDate string
	dateFormat := "2006-01-02"

	startDateTime, err := time.Parse(dateFormat, startDateString)
	startDate = startDateTime.Format(dateFormat)
	if err != nil || startDateString == "" {
		log.Println(err)
		startDate = time.Now().Format(dateFormat)
	}

	var endDate string

	endDateString := ctx.Query("endDate")
	endDateTime, err := time.Parse(dateFormat, endDateString)
	endDate = endDateTime.Format(dateFormat)
	if err != nil || endDateString == "" {
		endDate = time.Now().Add(time.Hour * 24).Format(dateFormat)
	}
	log.Println(startDate)
	log.Println(endDate)
	viewsCount, err := db_service.FetchViewsCount(startDate, endDate)
	if err != nil {
		log.Println(err)
	}

	ctx.HTML(200,"views_counter.tmpl.html", map[string]interface{}{
		"views":     viewsCount,
		"startDate": startDate,
		"endDate":   endDate,
	})
}
