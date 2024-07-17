package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework10/internal/app"
	"net/http"
	"strconv"
)

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		answer, err := a.CreateUser(reqBody.Nickname, reqBody.Email)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&answer))
	}
}

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		ID := c.Params.ByName("user_id")
		userID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		answer, err := a.UpdateUser(int64(userID), reqBody.Nickname, reqBody.Email)

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

		c.JSON(http.StatusOK, UserSuccessResponse(&answer))
	}
}

func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := c.Params.ByName("user_id")
		userID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		answer, err := a.GetUser(int64(userID))

		if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&answer))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := c.Params.ByName("user_id")
		userID, err := strconv.Atoi(ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		err = a.DeleteUser(int64(userID))

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

		c.Status(http.StatusOK)
	}
}

func panicServer(_ app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		panic("PANIC! ATTENTION!")
	}
}
