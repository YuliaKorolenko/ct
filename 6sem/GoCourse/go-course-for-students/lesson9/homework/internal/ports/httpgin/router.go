package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(title) объявления
	r.GET("/ads", getAdsList(a))                   // Метод для получения объявлений
	r.GET("/ads/:ad_id", getAd(a))                 // Метод для получения объявления по id
	r.DELETE("/ads/:ad_id", deleteAd(a))           // Метод для удаление ad по id
	r.POST("/users", createUser(a))                // Метод для создания (user)
	r.PUT("/users/:user_id", updateUser(a))        // Метод для обновления текста или заголовка объявления
	r.GET("/users/:user_id", getUser(a))           // Метод для получения user
	r.DELETE("/users/:user_id", deleteUser(a))     // Метод для удаления user
}
