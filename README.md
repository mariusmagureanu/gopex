# Kinly Controller

This document is draft proposal for a component which is meant to drive Kinly cloud onwards.

### What is it?

Kinly Controller is a component which leverages access to Pexip's client api along with providing custom functionalities built on top of the said api.

### How does it work?

Kinly Controller is at its core a simple http server. However, there are two ways of talking with the controller:
    
   * rest api
   
     Clients will use the controller's rest interface in order to fulfill their business tasks.
   
   * events / messages
   
     Clients can subscribe to a message broker by specifying what messages they are interested in. This use-case goes hand-in-hand with Pexip's server sent events.

The entire business logic is handled on the controller's side. The clients - on the other hand - are simple consumers for all the functionalities provided by the controller.

The controller does not keep any global state for itself, instead all state is persisted in a storage layer - Postgres - for the time being. The idea behind is to enable scalability.

The picture below outlines how the system is supposed to work.

![alt text](https://github.com/mariusmagureanu/gopex/blob/master/kc.png)

The controller can run in multiple instances as kubernetes pods. Clients will talk to any of these pods through a router (load balancer) and will not care which pod they're actually talking to.

From the other end of the spectrum, Pexip will serve sse's towards the controller. The controller will pick up the server sent events and publish messages against a message broker. The approach here is "fire&forget". It is up to the client(s) to subscribe on the same message broker and listen for the events they're interested in.
