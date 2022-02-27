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

func ViewsApiHandler(ctx *gin.Context) {
	startDateString := ctx.Query("startDate")
	endDateString := ctx.Query("endDate")

	var startDate string
	var endDate string

	dateFormat := "2006-01-02"

	if startDateString == "" {
		startDate = time.Now().Format(dateFormat)
	} else {
		startDateTime, _ := time.Parse(dateFormat, startDateString)
		startDate = startDateTime.Format(dateFormat)
	}
	if endDateString == "" {
		endDate = time.Now().Add(time.Hour * 24).Format(dateFormat)
	} else {
		endDateTime, _ := time.Parse(dateFormat, endDateString)
		endDate = endDateTime.Format(dateFormat)
	}

	viewsCount, err := db_service.FetchViewsCount(startDate, endDate)
	uniqueUsers, err := db_service.FetchDistinctUsersCount(startDate, endDate)
	if err != nil {
		log.Println(err)
	}

	ctx.JSON(200, map[string]interface{}{
		"views":       viewsCount,
		"startDate":   startDate,
		"endDate":     endDate,
		"uniqueUsers": uniqueUsers,
	})
}

func ViewsHandler(ctx *gin.Context) {
	startDateString := ctx.Query("startDate")
	endDateString := ctx.Query("endDate")

	var startDate string
	var endDate string

	dateFormat := "2006-01-02"

	if startDateString == "" {
		startDate = time.Now().Format(dateFormat)
	} else {
		startDateTime, _ := time.Parse(dateFormat, startDateString)
		startDate = startDateTime.Format(dateFormat)
	}
	if endDateString == "" {
		endDate = time.Now().Add(time.Hour * 24).Format(dateFormat)
	} else {
		endDateTime, _ := time.Parse(dateFormat, endDateString)
		endDate = endDateTime.Format(dateFormat)
	}

	viewsCount, err := db_service.FetchViewsCount(startDate, endDate)
	uniqueUsers, err := db_service.FetchDistinctUsersCount(startDate, endDate)
	if err != nil {
		log.Println(err)
	}

	ctx.HTML(200, "views_counter.tmpl.html", map[string]interface{}{
		"views":       viewsCount,
		"startDate":   startDate,
		"endDate":     endDate,
		"uniqueUsers": uniqueUsers,
	})
}
