package cheftodos

import (
	"sync"

	"github.com/duskhacker/cqrsnu/cafe"
	"github.com/pborman/uuid"
)

var (
	mutex sync.RWMutex
)

var ChefTodoList []*todoListGroup

type todoListItem struct {
	MenuNumber  int
	Description string
}

type todoListGroup struct {
	TabID uuid.UUID
	Items []todoListItem
}

func getTodoListGroup(tabID uuid.UUID) *todoListGroup {
	for _, list := range ChefTodoList {
		if list.TabID.String() == tabID.String() {
			return list
		}
	}
	return nil
}

func newTodoListGroup(tabID uuid.UUID, items []cafe.OrderedItem) *todoListGroup {
	group := todoListGroup{TabID: tabID}
	for _, item := range items {
		group.Items = append(group.Items, todoListItem{MenuNumber: item.MenuNumber, Description: item.Description})
	}
	return &group
}

func Init() {
	consumer := cafe.NewConsumer(cafe.FoodPreparedTopic, cafe.FoodPreparedTopic+"ChefTodoList", FoodPreparedHandler)
	consumers = append(consumers, consumer)
	consumer = cafe.NewConsumer(cafe.FoodOrderedTopic, cafe.FoodOrderedTopic+"ChefTodoList", FoodOrderedHandler)
	consumers = append(consumers, consumer)
}
