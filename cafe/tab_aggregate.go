package cafe

import (
	"fmt"
	"log"

	"github.com/pborman/uuid"
)

var (
	Tabs map[string]*Tab
)

func NewTabs() map[string]*Tab {
	return make(map[string]*Tab)
}

type Tab struct {
	ID                uuid.UUID
	TableNumber       int
	WaitStaff         string
	OutstandingDrinks []OrderedItem
	OutstandingFoods  []OrderedItem
	Open              bool
	ServedItemsValue  float64
}

func NewTab(id uuid.UUID, table int, staff string, drinks []OrderedItem, food []OrderedItem, open bool, siv float64) *Tab {
	mutex.Lock()
	defer mutex.Unlock()
	tab := &Tab{
		ID:                id,
		TableNumber:       table,
		WaitStaff:         staff,
		OutstandingDrinks: drinks,
		OutstandingFoods:  food,
		Open:              open,
		ServedItemsValue:  siv,
	}
	Tabs[id.String()] = tab
	return tab
}

func GetTab(id uuid.UUID) *Tab {
	tab, ok := Tabs[id.String()]
	if !ok {
		return nil
	}
	return tab
}

func (t Tab) AreDrinksOutstanding(drinks []OrderedItem) bool {
	for _, drink := range drinks {
		if indexOfOrderedItem(t.OutstandingDrinks, drink) < 0 {
			return false
		}
	}
	return true
}

func (t Tab) AreFoodsOutstanding(foods []OrderedItem) bool {
	for _, food := range foods {
		if indexOfOrderedItem(t.OutstandingFoods, food) < 0 {
			return false
		}
	}
	return true
}

func (t *Tab) AddServedItemsValue(items []OrderedItem) {
	for _, item := range items {
		t.ServedItemsValue += item.Price
	}
}

func (t *Tab) DeleteOutstandingDrinks(items []OrderedItem) error {
	for _, item := range items {
		drinks, err := deleteOrderedItem(t.OutstandingDrinks, item)
		if err != nil {
			return err
		}
		t.OutstandingDrinks = drinks
	}
	return nil
}

func (t *Tab) DeleteOutstandingFoods(items []OrderedItem) error {
	for _, item := range items {
		foods, err := deleteOrderedItem(t.OutstandingFoods, item)
		if err != nil {
			return err
		}
		t.OutstandingFoods = foods
	}
	return nil
}

func deleteOrderedItem(items []OrderedItem, item OrderedItem) ([]OrderedItem, error) {
	idx := indexOfOrderedItem(items, item)
	if idx < 0 {
		return nil, fmt.Errorf("no item %#v in tab", item)
	}
	a := make([]OrderedItem, len(items))
	n := copy(a, items)
	if n <= 0 {
		log.Fatalf("error copying data for deleteOutstandingDrinks")
	}
	return append(a[:idx], a[idx+1:]...), nil
}

func indexOfOrderedItem(items []OrderedItem, item OrderedItem) int {
	for i, e := range items {
		if e == item {
			return i
		}
	}
	return -1
}

// -

type OrderedItem struct {
	MenuNumber  int
	Description string
	IsDrink     bool
	Price       float64
}

func NewOrderedItem(menuNumber int, description string, isDrink bool, price float64) OrderedItem {
	return OrderedItem{
		MenuNumber:  menuNumber,
		Description: description,
		IsDrink:     isDrink,
		Price:       price,
	}
}
