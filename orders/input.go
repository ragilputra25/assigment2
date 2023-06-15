package orders

type SaveOrderInput struct {
	Items []SaveItemInput `json:"items" binding:"dive"`
}

type SaveItemInput struct {
	Code        string `json:"itemCode" binding:"required"`
	Description string `json:"description" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type UpdateOrderInput struct {
	Items []UpdateItemInput `json:"items" binding:"dive"`
}

type UpdateItemInput struct {
	ID          int    `json:"lineItemId" binding:"required"`
	Code        string `json:"itemCode" binding:"required"`
	Description string `json:"description" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type FindOrderInput struct {
	ID int `uri:"id" binding:"required"`
}
