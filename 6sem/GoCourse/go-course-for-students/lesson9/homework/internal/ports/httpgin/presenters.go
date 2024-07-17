package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	AuthorID  int64     `json:"author_id"`
	Published bool      `json:"published"`
	Time      time.Time `json:"time"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type deleteAdRequest struct {
	UserID int64 `json:"user_id"`
}

func AdSuccessResponse(ad *ads.Ad) gin.H {
	return gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		},
		"error": nil,
	}
}

func AdSuccessResponses(listAd *[]ads.Ad) gin.H {
	answer := make([]adResponse, 0)
	for _, val := range *listAd {
		answer = append(answer, adResponse{
			ID:        val.ID,
			Title:     val.Title,
			Text:      val.Text,
			AuthorID:  val.AuthorID,
			Published: val.Published,
			Time:      val.Time})
	}
	return gin.H{
		"data":  answer,
		"error": nil,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userResponse struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func UserSuccessResponse(user *user.User) gin.H {
	return gin.H{
		"data": userResponse{
			Id:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		"error": nil,
	}
}
