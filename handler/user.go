package handler

import (
	"log"
	"net/http"

	"go-hacktiv8-2/auth"
	"go-hacktiv8-2/helper"
	"go-hacktiv8-2/users"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService users.Service
	authService auth.Service
}

func NewUserHandler(userService users.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

// Register godoc
// @Summary Register new user
// @Description Register new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body users.RegisterUserInput true "User"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /register [post]
func (h *userHandler) RegisterUser(c *gin.Context) {
	var input users.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Unprocessable Entity", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := users.FormatRegister(newUser)
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// Login godoc
// @Summary Login user
// @Description Login user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body users.LoginUserInput true "User"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /login [post]
func (h *userHandler) LoginUser(c *gin.Context) {
	var input users.LoginUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Unprocessable Entity", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.LoginUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	findUser := h.authService.FindAuthUser(loggedInUser.ID)
	if findUser {
		err = h.authService.DeleteAuthUser(loggedInUser.ID)
		if err != nil {
			errorMessage := gin.H{"errors": err.Error()}
			response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	authData, err := h.authService.CreateAuth(uint64(loggedInUser.ID))
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	authD := auth.AuthDetails{}
	authD.UserID = authData.UserID
	authD.AuthUUID = authData.AuthUUID

	token, err := h.authService.Login(authD)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := users.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Successfuly logged in", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Logout(c *gin.Context) {
	au, err := auth.ExtractTokenAuth(c.Request)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	findUser := h.authService.FindAuth(au)
	log.Println(findUser)
	if findUser {
		err = h.authService.DeleteAuth(au)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.JSON(http.StatusUnauthorized, response)
			return
		}
	} else {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := helper.APIResponse("Successfuly logged out", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
