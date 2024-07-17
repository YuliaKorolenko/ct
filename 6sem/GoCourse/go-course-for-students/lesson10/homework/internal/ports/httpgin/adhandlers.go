package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	"homework10/internal/entities/ads"
	"net/http"
	"strconv"
)

func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		answer, err := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if errors.As(err, &app.ValidationApError{Err: app.ErrorValidation}) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&answer))
	}
}

func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ID := c.Params.ByName("ad_id")
		adID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		answer, err := a.ChangeAdStatus(int64(adID), reqBody.UserID, reqBody.Published)

		if errors.As(err, &app.PermissionDeniedError{Err: app.ErrorNoRights}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if errors.As(err, &app.ValidationApError{Err: app.ErrorValidation}) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&answer))
	}
}

func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ID := c.Params.ByName("ad_id")
		adID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		answer, err := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if errors.As(err, &app.PermissionDeniedError{Err: app.ErrorNoRights}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if errors.As(err, &app.ValidationApError{Err: app.ErrorValidation}) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&answer))
	}
}

func getAdsList(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		filter := ads.NewFilter()

		isPublished, isExist := c.GetQuery("published")
		if isExist {
			filter.UpdatePublished(isPublished)
		}

		isAuthor, isExist := c.GetQuery("author")
		if isExist {
			author, err := strconv.Atoi(isAuthor)
			if err == nil {
				filter.UpdateAuthor(int64(author))
			}
		}

		title, _ := c.GetQuery("title")
		filter.UpdateTitle(title)

		dayMy, _ := c.GetQuery("time")
		filter.UpdateTime(dayMy)

		answer, err := a.GetAdsList(filter)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponses(&answer))
	}
}

func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := c.Params.ByName("ad_id")
		adID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		answer, err := a.GetAd(int64(adID))

		if errors.As(err, &app.NoSuchAdError{Err: app.ErrorNoAd}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&answer))
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		ID := c.Params.ByName("ad_id")
		adID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		err = a.DeleteAd(int64(adID), reqBody.UserID)

		if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if errors.As(err, &app.NoSuchAdError{Err: app.ErrorNoAd}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.Status(http.StatusOK)
	}
}
