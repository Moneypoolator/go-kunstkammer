package api

import (
	"encoding/json"
	"fmt"
	"kunstkammer/internal/models"
	"net/http"
	"net/url"
)

// GetUserCards возвращает список карточек для пользователя с указанным идентификатором
func (kc *KaitenClient) GetUserCards(userID int) ([]models.Card, error) {
	resp, err := kc.doRequest("GET", fmt.Sprintf("/users/%d/cards", userID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cards []models.Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

// GetUserCards возвращает список карточек, в которых участвует пользователь с указанным идентификатором
// limit - количество карточек на странице
// offset - смещение относительно начала списка
func (kc *KaitenClient) getUserCards(userID int, limit int, offset int) ([]models.Card, error) {
	params := url.Values{}
	params.Add("member_ids", fmt.Sprintf("%d", userID))
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))
	params.Add("condition", "1")

	resp, err := kc.doRequestWithBody("GET", "/cards?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var cards []models.Card
	if err := kc.decodeResponse(resp, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (kc *KaitenClient) GetUserCardsByMemberIDs(userID int, limit int, offset int) ([]models.Card, error) {
	return kc.getUserCards(userID, limit, offset)
}

func (kc *KaitenClient) GetAllUserCards(userID int) ([]models.Card, error) {
	var allCards []models.Card
	limit := 100
	offset := 0

	for {
		cards, err := kc.getUserCards(userID, limit, offset)
		if err != nil {
			return nil, err
		}

		if len(cards) == 0 {
			break
		}

		allCards = append(allCards, cards...)
		offset += limit
	}

	return allCards, nil
}

// GetCards возвращает список всех карт
func (kc *KaitenClient) GetCards() ([]models.Card, error) {
	resp, err := kc.doRequest("GET", "/cards", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cards []models.Card
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (kc *KaitenClient) getCard(cardID int) (*models.Card, error) {
	resp, err := kc.doRequestWithBody("GET", fmt.Sprintf("/cards/%d", cardID), nil)
	if err != nil {
		return nil, err
	}

	var card models.Card
	if err := kc.decodeResponse(resp, &card); err != nil {
		return nil, err
	}

	return &card, nil
}
func (kc *KaitenClient) GetCard(cardID int) (*models.Card, error) {
	return kc.getCard(cardID)
}

func (kc *KaitenClient) createCard(card *models.Card) (*models.Card, error) {
	resp, err := kc.doRequestWithBody("POST", "/cards", card)
	if err != nil {
		return nil, err
	}

	var createdCard models.Card
	if err := kc.decodeResponse(resp, &createdCard); err != nil {
		return nil, err
	}

	return &createdCard, nil
}

func (kc *KaitenClient) updateCard(cardID int, update models.CardUpdate) error {
	resp, err := kc.doRequestWithBody("PATCH", fmt.Sprintf("/cards/%d", cardID), update)
	if err != nil {
		return err
	}

	return kc.decodeResponse(resp, nil)
}

func (kc *KaitenClient) CreateCard(card *models.Card) (*models.Card, error) {
	return kc.createCard(card)
}

func (kc *KaitenClient) UpdateCard(cardID int, update models.CardUpdate) error {
	return kc.updateCard(cardID, update)
}

func (kc *KaitenClient) updateCardProperties(cardID int, properties map[string]interface{}) error {
	updateData := models.CardUpdate{
		Properties: properties,
	}

	return kc.updateCard(cardID, updateData)
}

func (kc *KaitenClient) UpdateCardProperties(cardID int, properties map[string]interface{}) error {
	return kc.updateCardProperties(cardID, properties)
}

// DeleteCard удаляет карту по ID
func (kc *KaitenClient) DeleteCard(cardID int) error {
	resp, err := kc.doRequest("DELETE", fmt.Sprintf("/cards/%d", cardID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// TagRequest описывает данные для добавления тега
type TagRequest struct {
	Name string `json:"name"` // Имя тега
}

func (kc *KaitenClient) addTagToCard(cardID int, tagName string) error {
	tagRequest := TagRequest{Name: tagName}

	resp, err := kc.doRequestWithBody("POST", fmt.Sprintf("/cards/%d/tags", cardID), tagRequest)
	if err != nil {
		return err
	}

	return kc.decodeResponse(resp, nil)
}

func (kc *KaitenClient) AddTagToCard(cardID int, tagName string) error {
	return kc.addTagToCard(cardID, tagName)
}

// Add children Request
type AddChildrenRequest struct {
	CardID int `json:"card_id"`
}

func (kc *KaitenClient) addChildrenToCard(cardID int, childrenCardID int) error {
	requestData := AddChildrenRequest{CardID: childrenCardID}

	resp, err := kc.doRequestWithBody("POST", fmt.Sprintf("/cards/%d/children", cardID), requestData)
	if err != nil {
		return err
	}

	return kc.decodeResponse(resp, nil)
}

func (kc *KaitenClient) AddChindrenToCard(cardID int, childrenCardID int) error {
	return kc.addChildrenToCard(cardID, childrenCardID)
}

func PrintCardsList(cards []models.Card, userID int) {
	fmt.Printf("Cards count=%d :\n", len(cards))
	fmt.Printf("Cards for user %d:\n", userID)
	for _, card := range cards {
		fmt.Printf("Card ID: %d, Title: %s, Description: %s, ColumnID: %d, BoardID: %d\n",
			card.ID, card.Title, card.Description, card.ColumnID, card.BoardID)
	}
}
