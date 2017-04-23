package cafe

import (
	"encoding/json"
	"log"

	"github.com/pborman/uuid"
)

const (
	OpenTabTopic          = "OpenTab"
	PlaceOrderTopic       = "PlaceOrder"
	MarkDrinksServedTopic = "MarkDrinksServed"
	MarkFoodPreparedTopic = "MarkFoodPrepared"
	MarkFoodServedTopic   = "MarkFoodServed"
	CloseTabTopic         = "CloseTab"
)

type OpenTab struct {
	ID          uuid.UUID
	TableNumber int
	WaitStaff   string
}

func NewOpenTab(tableNumber int, waiter string) OpenTab {
	return OpenTab{
		ID:          uuid.NewRandom(),
		TableNumber: tableNumber,
		WaitStaff:   waiter,
	}
}

func (o OpenTab) FromJSON(data []byte) OpenTab {
	var err error
	err = json.Unmarshal(data, &o)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return o
}

// --

type PlaceOrder struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (po PlaceOrder) FromJSON(data []byte) PlaceOrder {
	var err error
	err = json.Unmarshal(data, &po)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return po
}

func NewPlaceOrder(id uuid.UUID, items []OrderedItem) PlaceOrder {
	return PlaceOrder{
		ID:    id,
		Items: items,
	}
}

// --

type markDrinksServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mds markDrinksServed) fromJSON(data []byte) markDrinksServed {
	var err error
	err = json.Unmarshal(data, &mds)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mds
}

func newMarkDrinksServed(id uuid.UUID, items []OrderedItem) markDrinksServed {
	return markDrinksServed{
		ID:    id,
		Items: items,
	}
}

// --

type MarkFoodPrepared struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mfp MarkFoodPrepared) FromJSON(data []byte) MarkFoodPrepared {
	var err error
	err = json.Unmarshal(data, &mfp)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mfp
}

func NewMarkFoodPrepared(id uuid.UUID, items []OrderedItem) MarkFoodPrepared {
	return MarkFoodPrepared{
		ID:    id,
		Items: items,
	}
}

// --

type markFoodServed struct {
	ID    uuid.UUID
	Items []OrderedItem
}

func (mfs markFoodServed) fromJSON(data []byte) markFoodServed {
	var err error
	err = json.Unmarshal(data, &mfs)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return mfs
}

func newMarkFoodServed(id uuid.UUID, items []OrderedItem) markFoodServed {
	return markFoodServed{
		ID:    id,
		Items: items,
	}
}

// --

type closeTab struct {
	ID         uuid.UUID
	AmountPaid float64
}

func (ct closeTab) fromJSON(data []byte) closeTab {
	var err error
	err = json.Unmarshal(data, &ct)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s\n'", err)
	}
	return ct
}

func newCloseTab(id uuid.UUID, amountPaid float64) closeTab {
	return closeTab{
		ID:         id,
		AmountPaid: amountPaid,
	}
}
