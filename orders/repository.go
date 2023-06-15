package orders

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Order, error)
	FindByID(ID int) (Order, error)
	Save(order Order) (Order, error)
	Update(order Order) (Order, error)
	Delete(order Order) (Order, error)
	Person() (Person, error)
}

type repository struct {
	db            *gorm.DB
	personService string
}

func NewRepository(db *gorm.DB, personService string) *repository {
	return &repository{db, personService}
}

func (r *repository) FindAll() ([]Order, error) {
	var orders []Order

	err := r.db.Preload("Items").Find(&orders).Error
	if err != nil {
		return orders, err
	}

	return orders, nil
}

func (r *repository) FindByID(ID int) (Order, error) {
	var order Order

	err := r.db.Preload("Items").Where("id = ?", ID).Find(&order).Error
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *repository) Save(order Order) (Order, error) {
	err := r.db.Create(&order).Error
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *repository) Update(order Order) (Order, error) {
	err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&order).Error
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *repository) Delete(order Order) (Order, error) {
	err := r.db.Delete(&order).Error
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *repository) Person() (Person, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, r.personService+"/data.php?qty=1&apikey=7f8fc96e-de1f-4aab-9c62-3dd1de365e66", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var person Person
	err = json.Unmarshal(resBody, &person)
	if err != nil {
		log.Fatal(err.Error())
	}

	return person, nil
}
