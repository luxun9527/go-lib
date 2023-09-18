package handler

import (
	"github.com/labstack/echo/v4"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/service"
	"net/http"
	"strconv"
)

func GetUser(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError,"parameter error")
	}
	user,err := service.UserRepository.Get(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
