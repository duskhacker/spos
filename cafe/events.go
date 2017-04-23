package cafe

import (
	"encoding/json"
	"log"

	"github.com/pborman/uuid"
)

const (
	TabOpenedTopic     = "TabOpened"
	FoodOrderedTopic   = "FoodOrdered"
	DrinksOrderedTopic = "DrinksOrdered"
	DrinksServedTopic  = "DrinksServed"
	FoodPreparedTopic  = "FoodPrepared"
	FoodServedTopic    = "FoodServed"
	TabClosedTopic     = "TabClosed"
	ExceptionTopic     = "Exception"
)

type TabOpened struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
}

func (t TabOpened) FromJSON(data []byte) TabOpened {
	var err error
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return t
}

func NewTabOpened(guid uuid.UUID, tableNumber int, waitStaff string) TabOpened {
	return TabOpened{
		ID:          guid,
		TableNumber: tableNumber,
		WaitStaff:   waitStaff,
	}
}

// --

type DrinksOrdered struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (do DrinksOrdered) FromJSON(data []byte) DrinksOrdered {
	var err error
	err = json.Unmarshal(data, &do)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return do
}

func NewDrinksOrdered(id uuid.UUID, items []OrderedItem) DrinksOrdered {
	return DrinksOrdered{
		ID:    id,
		Items: items,
	}
}

// --

type FoodOrdered struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fo FoodOrdered) FromJSON(data []byte) FoodOrdered {
	var err error
	err = json.Unmarshal(data, &fo)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fo
}

func NewFoodOrdered(id uuid.UUID, items []OrderedItem) FoodOrdered {
	return FoodOrdered{
		ID:    id,
		Items: items,
	}
}

// --

type DrinksServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (ds DrinksServed) FromJSON(data []byte) DrinksServed {
	var err error
	err = json.Unmarshal(data, &ds)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return ds
}

func NewDrinksServed(id uuid.UUID, items []OrderedItem) DrinksServed {
	return DrinksServed{
		ID:    id,
		Items: items,
	}
}

// --

type FoodPrepared struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fp FoodPrepared) FromJSON(data []byte) FoodPrepared {
	var err error
	err = json.Unmarshal(data, &fp)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fp
}

func NewFoodPrepared(id uuid.UUID, items []OrderedItem) FoodPrepared {
	return FoodPrepared{
		ID:    id,
		Items: items,
	}
}

// --

type FoodServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (fs FoodServed) FromJSON(data []byte) FoodServed {
	var err error
	err = json.Unmarshal(data, &fs)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return fs
}

func NewFoodServed(id uuid.UUID, items []OrderedItem) FoodServed {
	return FoodServed{
		ID:    id,
		Items: items,
	}
}

// --

type TabClosed struct {
	ID         uuid.UUID
	AmountPaid float64
	OrderValue float64
	TipValue   float64
}

func (tc TabClosed) FromJSON(data []byte) TabClosed {
	var err error
	err = json.Unmarshal(data, &tc)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return tc
}

func NewTabClosed(id uuid.UUID, amountPaid, orderValue, tipValue float64) TabClosed {
	return TabClosed{
		ID:         id,
		AmountPaid: amountPaid,
		OrderValue: orderValue,
		TipValue:   tipValue,
	}
}

// --

type Exception struct {
	Type    string
	Message string
}

func (e Exception) FromJSON(data []byte) Exception {
	var err error
	err = json.Unmarshal(data, &e)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return e
}

func NewException(t string, msg string) Exception {
	return Exception{Type: t, Message: msg}
}

func (e Exception) Error() string {
	return e.Type + ":" + e.Message
}

var TabNotOpenException = NewException("TabNotOpen", "Cannot Place order without open Tab")
var DrinksNotOutstanding = NewException("DrinksNotOutstanding", "Cannot serve unordered drinks")
var FoodsNotOutstanding = NewException("FoodsNotOutstanding", "Cannot prepare unordered food")
