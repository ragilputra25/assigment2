package handler

import (
	"net/http"

	"go-hacktiv8-2/auth"
	"go-hacktiv8-2/helper"
	"go-hacktiv8-2/orders"
	"go-hacktiv8-2/users"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	service     orders.Service
	authService auth.Service
	userService users.Service
}

func NewHandler(service orders.Service, authService auth.Service, userService users.Service) *orderHandler {
	return &orderHandler{service, authService, userService}
}

// FindAll godoc
// @Summary Show all orders
// @Description Get all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Response
// @Router /orders [get]
func (h *orderHandler) FindAll(c *gin.Context) {
	findOrders, err := h.service.FindAll()
	if err != nil {
		response := helper.APIResponse("Error to get orders", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of orders", http.StatusOK, "success", orders.FormatOrders(findOrders))
	c.JSON(http.StatusOK, response)
}

// FindByID godoc
// @Summary Show orders by user id
// @Description Get orders by user id
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_id path int true "User ID"
// @Success 200 {object} helper.Response
// @Router /orders/{order_id} [get]
func (h *orderHandler) FindByID(c *gin.Context) {
	var input orders.FindOrderInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to get detail of order", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	findOrder, err := h.service.FindByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of order", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if findOrder.ID == 0 {
		response := helper.APIResponse("Order not found", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Detail of order", http.StatusOK, "success", orders.FormatOrder(findOrder))
	c.JSON(http.StatusOK, response)
}

// Save godoc
// @Summary Create new order
// @Description Create new order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body orders.SaveOrderInput true "Order"
// @Success 200 {object} helper.Response
// @Router /orders [post]
func (h *orderHandler) Save(c *gin.Context) {
	var input orders.SaveOrderInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create order", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tokenAuth, err := auth.ExtractTokenAuth(c.Request)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	foundAuth, err := h.authService.FetchAuth(tokenAuth)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userID := foundAuth.UserID

	user, err := h.userService.GetUserById(int(userID))
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userName := user.Name

	newOrder, err := h.service.Save(userName, input)
	if err != nil {
		response := helper.APIResponse("Failed to create order", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create order", http.StatusOK, "success", orders.FormatOrder(newOrder))
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary Update order
// @Description Update order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_id path int true "Order ID"
// @Param order body orders.UpdateOrderInput true "Order"
// @Success 200 {object} helper.Response
// @Router /orders/{order_id} [put]
func (h *orderHandler) Update(c *gin.Context) {
	var inputID orders.FindOrderInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update order", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData orders.UpdateOrderInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update order", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tokenAuth, err := auth.ExtractTokenAuth(c.Request)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	foundAuth, err := h.authService.FetchAuth(tokenAuth)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userID := foundAuth.UserID

	user, err := h.userService.GetUserById(int(userID))
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userName := user.Name

	updatedOrder, err := h.service.Update(inputID, inputData, userName)
	if err != nil {
		response := helper.APIResponse("Failed to update order", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update order", http.StatusOK, "success", orders.FormatOrder(updatedOrder))
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Delete order
// @Description Delete order
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_id path int true "Order ID"
// @Success 200 {object} helper.Response
// @Router /orders/{order_id} [delete]
func (h *orderHandler) Delete(c *gin.Context) {
	var input orders.FindOrderInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to delete order", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	tokenAuth, err := auth.ExtractTokenAuth(c.Request)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	foundAuth, err := h.authService.FetchAuth(tokenAuth)
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userID := foundAuth.UserID

	user, err := h.userService.GetUserById(int(userID))
	if err != nil {
		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userName := user.Name

	_, err = h.service.Delete(input, userName)
	if err != nil {
		response := helper.APIResponse("Failed to delete order", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to delete order", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

// FindOrderPerson godoc
// @Summary Show order person by order id
// @Description Get order person by order id
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_id path int true "Order ID"
// @Success 200 {object} helper.Response
// @Router /orders/person/{order_id} [get]
func (h *orderHandler) FindOrderPerson(c *gin.Context) {
	var input orders.FindOrderInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to get detail of order", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	findOrder, err := h.service.FindOrderPerson(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of order", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Detail of order", http.StatusOK, "success", orders.FormatOrderPerson(findOrder))
	c.JSON(http.StatusOK, response)
}
