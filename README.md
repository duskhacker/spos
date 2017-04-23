# CQRSNU

This is my take on [Edument's CQRS and Intentful Testing Tutorial](http://cqrs.nu/).

From Edument's site:

> What is this?

> A bunch of C# code to help you get started with writing intentful tests for a domain, expressing it as commands, >
events and exceptions. These ideas are often associated with the CQRS pattern.


Edument's tutorial provides a "starter kit" written in C# that provides the plumbing necessary to complete the tutorial.
[Edument CQRS and Intentful Testing Starter Kit](https://github.com/edumentab/cqrs-starter-kit). I'll go over the big
ideas behind the tutorial here, but I highly recommend browsing through Edument's site and their Github repo for deeper
insights into these concepts.

Since I'm interested in CQRS, and my day job requires me to be a Go programmer, I have attempted to complete the
tutorial using [Go](https://golang.org) as the programming language and [NSQ](http://nsq.io/) for the message bus.

## The Big Idea

The big idea behind the tutorial, and CQRS in general is that you write your software by concentrating on the verbs of
the domain, and turn those verbs into commands. For instance, in the Cafe domain of the tutorial, you would open a tab
(`OpenTab` command), Place an order (`PlaceOrder` command), Prepare Food (`MarkFoodPrepared` command), serve food
(`MarkFoodServed` command). These commands are received by processors that track state and send events. Again, in the
Cafe domain, some of the events would be `TabOpened`, `OrderPlaced`,`FoodPrepared`,`FoodServed`. This is the "C" part of
"CQRS".
 
For the query portion (the "Q" portion of "CQRS"), different processes (not necessarily a process as in Unix process)
listen to the generated events and build up their own models that are specifically for querying. For instance, in the
Cafe domain, the chefs will need a list of items to prepare. A `ChefsTodo` process can be implemented that listens to
the `FoodOrdered` events and build up a list of items that must be prepared. This process would also listen for
`FoodPrepared` events and remove those items from the `ChefsTodos` list.
  
Lastly, not everything in a domain is worth using the CQRS methodology. For instance, the lists of Wait Staff, Food an
Drink items, etc. are just data that a simple CRUD interface is sufficient to maintain.
  
## This Repository

This repository implements the full Cafe, but due to time constraints, only the `ChefsTodos` "Read Model" is
implemented. Further, an HTTP API for opening a Tab, placing an order, and getting the chef's todos is also implemented
using [Gin](https://github.com/gin-gonic/gin). Since I'm a back-end developer and not interested in doing front-end
work, that's where my implementation stops.