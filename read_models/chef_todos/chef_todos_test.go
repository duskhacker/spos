package cheftodos

import (
	"fmt"

	"github.com/pborman/uuid"

	"github.com/duskhacker/cqrsnu/internal/github.com/bitly/go-nsq"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"

	"github.com/duskhacker/cqrsnu/cafe"
)

var pf = fmt.Printf

var testConsumers []*nsq.Consumer

var _ = Describe("Chef TODOs", func() {

	var (
		openTabCmd cafe.OpenTab
		tabID      uuid.UUID
		drinks     []cafe.OrderedItem
		food       []cafe.OrderedItem
	)

	BeforeEach(func() {
		openTabCmd = cafe.NewOpenTab(1, "Kinessa")
		tabID = openTabCmd.ID

		drinks = []cafe.OrderedItem{}
		drinks = append(drinks, cafe.NewOrderedItem(1, "Patron", true, 5.00))
		drinks = append(drinks, cafe.NewOrderedItem(2, "Scotch", true, 3.50))

		food = []cafe.OrderedItem{}
		food = append(food, cafe.NewOrderedItem(1, "Steak", false, 15.00))
		food = append(food, cafe.NewOrderedItem(2, "Burger", false, 8.00))
	})

	AfterEach(func() {
		stopAllTestConsumers()
	})

	Describe("todoList Group", func() {
		BeforeEach(func() {
			cafe.Send(cafe.OpenTabTopic, openTabCmd)
			cafe.Send(cafe.PlaceOrderTopic, cafe.NewPlaceOrder(tabID, append(food, drinks...)))
			f := func() int { return len(ChefTodoList) }
			Eventually(f).ShouldNot(BeZero())
		})

		It("creates a new group", func() {
			Expect(getTodoListGroup(tabID).Items).To(HaveLen(2))
		})

		Describe("with created group", func() {
			It("Removes an item", func() {
				cafe.Send(cafe.MarkFoodPreparedTopic, cafe.NewMarkFoodPrepared(tabID, food[0:1]))
				f := func() int {
					group := getTodoListGroup(tabID)
					if group == nil {
						return 100
					}
					return len(group.Items)
				}

				Eventually(f).Should(Equal(1))
				Expect(getTodoListGroup(tabID).Items).To(HaveLen(1))
			})
		})
	})
})

func newTestConsumer(topic, channel string, f func(*nsq.Message) error) {
	testConsumers = append(testConsumers, cafe.NewConsumer(topic, channel, f))
}

func stopAllTestConsumers() {
	for _, consumer := range testConsumers {
		consumer.Stop()
	}
}

func listenForUnexpectedException() {
	f := func(m *nsq.Message) error {
		pf("EXCEPTION: %#v\n", new(cafe.Exception).FromJSON(m.Body))
		return nil
	}
	newTestConsumer(cafe.ExceptionTopic, cafe.ExceptionTopic+"UnexpectedExceptionConsumer", f)
}
