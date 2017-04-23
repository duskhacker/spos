package cafe

import (
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"
)

var _ = Describe("TabAggegate", func() {
	BeforeEach(func() {
		Tabs = NewTabs()
	})
	Describe("DeleteOrderedItem", func() {
		It("deletes an item", func() {
			drink0 := NewOrderedItem(5, "drink1", true, 0.0)
			drink1 := NewOrderedItem(1, "drink2", true, 0.0)
			drink2 := NewOrderedItem(7, "drink3", true, 0.0)

			drinks := []OrderedItem{}
			drinks = append(drinks, drink0)
			drinks = append(drinks, drink1)
			drinks = append(drinks, drink2)

			tab := NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks[1:2])
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink0, drink2}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks[0:1])
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink1, drink2}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks[2:])
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{drink0, drink1}))

			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			tab.DeleteOutstandingDrinks(drinks)
			Expect(tab.OutstandingDrinks).To(ConsistOf([]OrderedItem{}))
		})
	})

	Describe("AreDrinksOutstanding", func() {
		var (
			drinks []OrderedItem
			tab    *Tab
		)

		BeforeEach(func() {
			drink0 := NewOrderedItem(5, "drink1", true, 0.0)
			drink1 := NewOrderedItem(1, "drink2", true, 0.0)
			drink2 := NewOrderedItem(7, "drink3", true, 0.0)

			drinks = []OrderedItem{}
			drinks = append(drinks, drink0)
			drinks = append(drinks, drink1)
			drinks = append(drinks, drink2)

		})

		It("returns true if all are outstanding", func() {
			tab = NewTab(nil, 0, "", drinks, nil, false, 0)
			Expect(tab.AreDrinksOutstanding(drinks)).To(BeTrue())
		})

		It("returns false if any are not outstanding", func() {
			tab = NewTab(nil, 0, "", drinks[:2], nil, false, 0)
			Expect(tab.AreDrinksOutstanding(drinks)).To(BeFalse())

		})

	})

})
