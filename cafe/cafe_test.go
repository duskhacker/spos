package cafe

import (
	"fmt"

	"github.com/pborman/uuid"

	"github.com/duskhacker/cqrsnu/internal/github.com/bitly/go-nsq"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"
)

var pf = fmt.Printf

var testConsumers []*nsq.Consumer

var _ = Describe("Main", func() {

	var (
		openTabCmd OpenTab
		tabID      uuid.UUID
		drinks     []OrderedItem
		food       []OrderedItem
	)

	BeforeEach(func() {
		Tabs = NewTabs()
		openTabCmd = NewOpenTab(1, "Kinessa")
		tabID = openTabCmd.ID

		drinks = []OrderedItem{}
		drinks = append(drinks, NewOrderedItem(1, "Patron", true, 5.00))
		drinks = append(drinks, NewOrderedItem(2, "Scotch", true, 3.50))

		food = []OrderedItem{}
		food = append(food, NewOrderedItem(1, "Steak", false, 15.00))
		food = append(food, NewOrderedItem(2, "Burger", false, 8.00))
	})

	AfterEach(func() {
		stopAllTestConsumers()
	})

	Describe("Tab", func() {
		It("opens a tab", func() {
			done := make(chan bool)

			newTestConsumer(TabOpenedTopic, TabOpenedTopic+"TestConsumer",
				func(m *nsq.Message) error {
					defer GinkgoRecover()
					Expect(new(TabOpened).FromJSON(m.Body)).To(Equal(NewTabOpened(tabID, 1, "Kinessa")))
					done <- true
					return nil
				})

			Send(OpenTabTopic, openTabCmd)

			Eventually(done).Should(Receive(BeTrue()), "No TabOpened received")
		})
	})

	Describe("Ordering", func() {
		Describe("with no tab opened", func() {
			It("receives error", func() {
				done := make(chan bool)
				command := NewPlaceOrder(nil, nil)

				newTestConsumer(ExceptionTopic, ExceptionTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						Expect(new(Exception).FromJSON(m.Body)).To(Equal(TabNotOpenException))
						done <- true
						return nil
					})

				Send(PlaceOrderTopic, command)

				Eventually(done).Should(Receive(BeTrue()), "TabNotOpenException Exception not Raised")
			})
		})

		Describe("with tab opened", func() {
			var (
				foodOrderedDone   = make(chan bool)
				drinksOrderedDone = make(chan bool)
			)

			BeforeEach(func() {

				newTestConsumer(DrinksOrderedTopic, DrinksOrderedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						order := new(DrinksOrdered).FromJSON(m.Body)
						if len(order.Items) > 0 {
							drinksOrderedDone <- true
						}
						return nil
					})

				newTestConsumer(FoodOrderedTopic, FoodOrderedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						order := new(FoodOrdered).FromJSON(m.Body)
						if len(order.Items) > 0 {
							foodOrderedDone <- true
						}
						return nil
					})

				newTestConsumer(ExceptionTopic, ExceptionTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						ex := new(Exception).FromJSON(m.Body)
						Expect(ex).To(BeNil())
						return nil
					})

				Send(OpenTabTopic, openTabCmd)
			})

			It("drinks", func() {

				Send(PlaceOrderTopic, NewPlaceOrder(tabID, drinks))

				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "DrinksOrdered not received")
			})

			It("food", func() {
				Send(PlaceOrderTopic, NewPlaceOrder(tabID, food))

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "FoodOrdered not received")
			})

			It("food and drink", func() {
				Send(PlaceOrderTopic, NewPlaceOrder(tabID, append(food, drinks...)))

				Eventually(foodOrderedDone).Should(Receive(BeTrue()), "FoodOrdered not received")
				Eventually(drinksOrderedDone).Should(Receive(BeTrue()), "DrinksOrdered not received")
			})
		})
	})

	Describe("Serving Drinks", func() {
		BeforeEach(func() {
			Send(OpenTabTopic, openTabCmd)
		})

		Describe("with 1 drink ordered", func() {
			BeforeEach(func() {
				Send(PlaceOrderTopic, NewPlaceOrder(tabID, drinks[:1]))
			})

			It("generates exception if second drink is marked served", func() {
				done := make(chan bool)

				newTestConsumer(ExceptionTopic, ExceptionTopic+"TestConsumer", func(m *nsq.Message) error {
					ex := new(Exception).FromJSON(m.Body)
					defer GinkgoRecover()
					Expect(ex).To(Equal(DrinksNotOutstanding))
					done <- true
					return nil
				})

				Send(MarkDrinksServedTopic, newMarkDrinksServed(tabID, drinks[1:2]))

				Eventually(done).Should(Receive(BeTrue()), "DrinksNotOutstanding Exception not Raised")
			})
		})

		Describe("with drinks ordered", func() {
			var (
				drinksServedDone chan bool
			)

			BeforeEach(func() {
				drinksServedDone = make(chan bool)

				Send(PlaceOrderTopic, NewPlaceOrder(tabID, drinks))
			})

			It("marks drinks served", func() {

				newTestConsumer(DrinksServedTopic, DrinksServedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(DrinksServed).FromJSON(m.Body)
						Expect(evt.Items).To(Equal(drinks))
						drinksServedDone <- true
						return nil
					})

				Send(MarkDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "DrinksServed not received")
			})

			It("does not allow drinks to be served twice", func() {
				rcvdException := make(chan bool)
				newTestConsumer(ExceptionTopic, ExceptionTopic+"TestExceptionConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						ex := new(Exception).FromJSON(m.Body)
						Expect(ex).To(Equal(DrinksNotOutstanding))
						rcvdException <- true
						return nil
					})

				newTestConsumer(DrinksServedTopic, DrinksServedTopic+"TestConsumer", func(m *nsq.Message) error {
					defer GinkgoRecover()
					evt := new(DrinksServed).FromJSON(m.Body)
					Expect(evt.Items).To(Equal(drinks))
					drinksServedDone <- true
					return nil
				})

				Send(MarkDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
				Eventually(drinksServedDone).Should(Receive(BeTrue()), "DrinksServed not received")

				Send(MarkDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
				Eventually(rcvdException).Should(Receive(BeTrue()), "DrinksNotOutstanding exception not received")
			})
		})

	})

	Describe("Food", func() {

		BeforeEach(func() {
			Send(OpenTabTopic, openTabCmd)
			Send(PlaceOrderTopic, NewPlaceOrder(tabID, food))
		})

		Describe("prepare", func() {
			It("marks food prepared", func() {
				rcvdFoodPrepared := make(chan bool)
				newTestConsumer(FoodPreparedTopic, FoodPreparedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(FoodPrepared).FromJSON(m.Body)
						Expect(evt.Items).To(Equal(food))
						rcvdFoodPrepared <- true
						return nil
					})

				Send(MarkFoodPreparedTopic, NewMarkFoodPrepared(tabID, food))
				Eventually(rcvdFoodPrepared).Should(Receive(BeTrue()), "FoodPrepared not received")
			})
		})

		Describe("serve", func() {
			BeforeEach(func() {
				rcvdFoodPrepared := make(chan bool)
				newTestConsumer(FoodPreparedTopic, FoodPreparedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						rcvdFoodPrepared <- true
						return nil
					})

				Send(MarkFoodPreparedTopic, NewMarkFoodPrepared(tabID, food))
				Eventually(rcvdFoodPrepared).Should(Receive(BeTrue()), "FoodPrepared not received")

			})

			It("marks food served", func() {
				listenForUnexpectedException()
				rcvdFoodServed := make(chan bool)
				newTestConsumer(FoodServedTopic, FoodServedTopic+"TestConsumer",
					func(m *nsq.Message) error {
						defer GinkgoRecover()
						evt := new(FoodServed).FromJSON(m.Body)
						Expect(evt.Items).To(Equal(food))
						rcvdFoodServed <- true
						return nil
					})

				Send(MarkFoodServedTopic, newMarkFoodServed(tabID, food))
				Eventually(rcvdFoodServed).Should(Receive(BeTrue()), "FoodServed not received")
			})
		})
	})

	Describe("Closing Tab", func() {
		BeforeEach(func() {
			Send(OpenTabTopic, openTabCmd)
			Send(PlaceOrderTopic, NewPlaceOrder(tabID, append(food, drinks...)))
			Send(MarkDrinksServedTopic, newMarkDrinksServed(tabID, drinks))
			Send(MarkFoodPreparedTopic, NewMarkFoodPrepared(tabID, food))
			Send(MarkFoodServedTopic, newMarkFoodServed(tabID, food))

			rcvdDrinksServed := make(chan bool)
			newTestConsumer(DrinksServedTopic, DrinksServedTopic+"TestConsumer",
				func(msg *nsq.Message) error {
					rcvdDrinksServed <- true
					return nil
				})

			Eventually(rcvdDrinksServed).Should(Receive(BeTrue()))
		})

		Describe("with tip", func() {
			It("closes tab", func() {
				tabClosedReceived := make(chan bool)
				newTestConsumer(TabClosedTopic, TabClosedTopic+"TestConsumer",
					func(msg *nsq.Message) error {
						evt := new(TabClosed).FromJSON(msg.Body)
						defer GinkgoRecover()
						Expect(evt.AmountPaid).To(Equal(31.50 + 0.50))
						Expect(evt.OrderValue).To(Equal(31.50))
						Expect(evt.TipValue).To(Equal(0.5))
						tabClosedReceived <- true
						return nil
					})

				Send(CloseTabTopic, newCloseTab(tabID, 31.50+0.50))

				Eventually(tabClosedReceived).Should(Receive(BeTrue()), "TabClosed not received")
			})
		})
	})
})

func newTestConsumer(topic, channel string, f func(*nsq.Message) error) {
	testConsumers = append(testConsumers, NewConsumer(topic, channel, f))
}

func stopAllTestConsumers() {
	for _, consumer := range testConsumers {
		consumer.Stop()
	}
}

func listenForUnexpectedException() {
	f := func(m *nsq.Message) error {
		pf("EXCEPTION: %#v\n", new(Exception).FromJSON(m.Body))
		return nil
	}
	newTestConsumer(ExceptionTopic, ExceptionTopic+"UnexpectedExceptionConsumer", f)
}
