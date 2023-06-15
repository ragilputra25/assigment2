package orders

import "time"

// Order represents the model for an order
type Order struct {
	ID           int
	CustomerName string
	CreatedAt    time.Time
	Items        []Item `gorm:"Constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Item represents the model for an item in the order
type Item struct {
	ID          int
	Code        string
	Description string
	Qty         int
	OrderID     int
}

// OrderPerson represents the input for finding an order
type OrderPerson struct {
	Order  Order
	Person Person
}

// Person represents the model for a person
type Person struct {
	Status struct {
		Code        int
		Description string
	}
	Result []struct {
		Firstname string
		Lastname  string
		Username  string
		Phone     string
		Email     string
		UUID      string
	}
}
