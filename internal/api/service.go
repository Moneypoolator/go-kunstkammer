package api

import (
	"fmt"
	"kunstkammer/internal/models"
)

type CRUD interface {
	Create(interface{}) (interface{}, error)
	Update(int, interface{}) error
	Delete(int) error
}

type CardService struct {
	client *KaitenClient
}

func (cs *CardService) Create(item interface{}) (interface{}, error) {
	card, ok := item.(*models.Card)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	return cs.client.CreateCard(card)
}

func (cs *CardService) Update(id int, item interface{}) error {
	update, ok := item.(models.CardUpdate)
	if !ok {
		return fmt.Errorf("invalid type")
	}
	return cs.client.UpdateCard(id, update)
}

func (cs *CardService) Delete(id int) error {
	return cs.client.DeleteCard(id)
}
