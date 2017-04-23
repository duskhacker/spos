package api

import (
	"net/http"

	"github.com/duskhacker/cqrsnu/cafe"
	"github.com/duskhacker/cqrsnu/read_models/chef_todos"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

// Test Data
var (
	food = []cafe.OrderedItem{}
)

// Test Data
func init() {
	food = append(food, cafe.NewOrderedItem(0, "Steak", false, 15.00))
	food = append(food, cafe.NewOrderedItem(1, "Burger", false, 8.00))
}

func GinEngine() *gin.Engine {
	r := gin.Default()

	r.POST("/opentab", OpenTab)
	r.POST("/placeorder", PlaceOrder)
	r.GET("/cheftodolist", ChefTodoList)

	return r
}

func OpenTab(c *gin.Context) {
	json := struct {
		WaitStaff string `json:"waitstaff" binding:"required"`
	}{}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openTabCmd := cafe.NewOpenTab(1, json.WaitStaff)
	cafe.Send(cafe.OpenTabTopic, openTabCmd)

	c.JSON(http.StatusAccepted, gin.H{"tabID": openTabCmd.ID.String()})
}

func PlaceOrder(c *gin.Context) {
	json := struct {
		TabID       string `json:"tabID" binding:"required"`
		MenuNumbers []int  `json:"items" binding:"required"`
	}{}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := []cafe.OrderedItem{}

	for _, menuNumber := range json.MenuNumbers {
		order = append(order, food[menuNumber])
	}

	cafe.Send(cafe.PlaceOrderTopic, cafe.NewPlaceOrder(uuid.Parse(json.TabID), order))
	c.String(http.StatusAccepted, "")
}

type Item struct {
	MenuNumber  int
	Description string
}

type Tab struct {
	TabID string
	Items []Item
}

type ChefTodoListResponse struct {
	Tabs []Tab
}

func ChefTodoList(c *gin.Context) {
	r := ChefTodoListResponse{}

	for _, todoList := range cheftodos.ChefTodoList {
		t := Tab{TabID: todoList.TabID.String()}
		for _, item := range todoList.Items {
			t.Items = append(t.Items, Item{MenuNumber: item.MenuNumber, Description: item.Description})
		}
		r.Tabs = append(r.Tabs, t)
	}

	c.JSON(http.StatusOK, r)
}
