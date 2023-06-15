package orders

import "errors"

type Service interface {
	FindAll() ([]Order, error)
	FindByID(ID FindOrderInput) (Order, error)
	Save(customerName string, orderInput SaveOrderInput) (Order, error)
	Update(ID FindOrderInput, orderInput UpdateOrderInput, customerName string) (Order, error)
	Delete(ID FindOrderInput, customerName string) (Order, error)
	FindOrderPerson(ID FindOrderInput) (OrderPerson, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() ([]Order, error) {
	orders, err := s.repository.FindAll()
	if err != nil {
		return orders, err
	}

	return orders, nil
}

func (s *service) FindByID(ID FindOrderInput) (Order, error) {
	order, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (s *service) Save(customerName string, orderInput SaveOrderInput) (Order, error) {
	order := Order{}
	order.CustomerName = customerName

	var orderItems []Item
	for _, item := range orderInput.Items {
		orderItem := Item{}
		orderItem.Code = item.Code
		orderItem.Description = item.Description
		orderItem.Qty = item.Quantity
		orderItems = append(orderItems, orderItem)
	}
	order.Items = orderItems

	newOrder, err := s.repository.Save(order)
	if err != nil {
		return newOrder, err
	}

	return newOrder, nil
}

func (s *service) Update(ID FindOrderInput, orderInput UpdateOrderInput, customerName string) (Order, error) {
	order, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return order, err
	}

	if order.CustomerName != customerName {
		return order, errors.New("unauthorized")
	}

	order.CustomerName = customerName

	var orderItems []Item
	for _, item := range orderInput.Items {
		orderItem := Item{}
		for i := 0; i < len(order.Items); i++ {
			if order.Items[i].ID != item.ID {
				return order, errors.New("unauthorized")
			}
		}
		orderItem.ID = item.ID
		orderItem.Code = item.Code
		orderItem.Description = item.Description
		orderItem.Qty = item.Quantity
		orderItems = append(orderItems, orderItem)
	}

	order.Items = orderItems

	updatedOrder, err := s.repository.Update(order)
	if err != nil {
		return updatedOrder, err
	}

	return updatedOrder, nil
}

func (s *service) Delete(ID FindOrderInput, customerName string) (Order, error) {
	order, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return order, err
	}

	if order.CustomerName != customerName {
		return order, errors.New("unauthorized")
	}

	deletedOrder, err := s.repository.Delete(order)
	if err != nil {
		return deletedOrder, err
	}

	return deletedOrder, nil
}

func (s *service) FindOrderPerson(ID FindOrderInput) (OrderPerson, error) {
	order, err := s.repository.FindByID(ID.ID)
	if err != nil {
		return OrderPerson{}, err
	}

	person, err := s.repository.Person()
	if err != nil {
		return OrderPerson{}, err
	}

	orderPerson := OrderPerson{}
	orderPerson.Order = order
	orderPerson.Person = person

	return orderPerson, nil
}
