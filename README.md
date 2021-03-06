<p align="center">
<img src="https://github.com/carrot-ar/carrot-ios/wiki/resources/Carrot@2x.png" alt="Carrot" width="300">
</p>

<p align="center">
<a href="https://travis-ci.org/carrot-ar/carrot"><img src="https://travis-ci.org/carrot-ar/carrot.svg?branch=master" alt="build status"></a>
<a href=""><img src="https://codecov.io/gh/carrot-ar/carrot/branch/master/graph/badge.svg" alt="code coverage"></a>
</p>

Carrot is an easy-to-use, real-time framework for building applications with multi-device AR capabilities. It works using WebSockets, Golang, client libraries written for iOS, and a unique location tracking system based on iBeacons that we aptly named The Picnic Protocol. Using Carrot, multi-device AR apps can be created with high accuracy location tracking to provide rich and lifelike experiences. To see for yourself, check out Scribbles, a multiplayer drawing application made with Carrot. You can see a demo video [here](https://www.youtube.com/watch?v=6EVtb0pJPgk) and the code [here](https://github.com/carrot-ar/scribbles).

To see documentation for the iOS Client library visit the [README for carrot-ios](https://github.com/carrot-ar/carrot-ios/blob/master/README.md)

|    | 🗂 Table of Contents |
|:--:|----------------------
| ✨ | [Features](#features)
| 📋 | [To-Do](#to-do)
| ⚙️ | [Design](#design)
| 🛠 | [Building an Application with Carrot](#building-an-application-with-carrot)
| 🥗 | [The Picnic Protocol](#the-picnic-protocol)
| ✉️ | [Message Format](#message-format)
| 🎙 | [Sending Messages to Carrot](#sending-messages-to-carrot)
| 📨 | [Receiving Messages from Carrot](#receiving-messages-from-carrot)
| 📺 | [Broadcasting Responses](#broadcasting-responses)
| 🌎 | [Sessions](#sessions)

## Features
- Rapid development of multi-device AR applications with little Go or server knowledge
- WebSocket connection and state management
- High accuracy AR location tracking with the [Picnic Protocol](https://github.com/carrot-ar/carrot-ios#-the-picnic-protocol)
- Sessions and session management
- Middleware 
- Extensible controllers
- Custom endpoints 
- Performance optimizations
- High throughput (30k messages/second tested on 4 CPU test machine)
- ~100 microsecond average to service one request 

## To-Do
- Support for external session management using Redis, memcached, etc
- Bug fixes for picnic protocol implementation 
- Object Relational Mapping library to have a true Model-Controller design
- [Universal Scene Description](https://github.com/PixarAnimationStudios/USD) support 

## Design

Bellow is a high level design of the Carrot framework. More detailed diagrams will be provided when time permits.


<img src="https://i.imgur.com/sHkF7Hl.png" alt="carrot flow">


## Building an application with Carrot

Building applications with Carrot is incredibly simple. Check out this echo application that echos a payload from one device into the AR space of all connected devices: 

``` go
package main

import (
	"fmt"
	"github.com/carrot-ar/carrot"
)

// Controller declaration
type EchoController struct{}

//Controller method implementation
func (c *EchoController) Echo(req *carrot.Request, br *carrot.Broadcast) {
	message, err := carrot.CreateDefaultResponse(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	br.Broadcast(message)
}

func main() {

	// Register endpoint by providing endpoint, controller, and method, and if endpoint requires streaming
	carrot.Add("echo", EchoController{}, "Echo", true)

	// Run the server to handle traffic
	carrot.Run()
}
```
The example above ommits extra functionality to showcase the basic components required to connect carrot to your application. The required components are the following:

* Import inclusion
* Controller declaration
* Controller method(s) implementation
* Main method that
	* Registers connections between methods and controllers for Carrot to route and maintain
	* Runs the Carrot server
	
Controller methods receieve requests and broadcast responses to clients. Requests can be passed as-is, demonstrated with the CreateDefaultResponse method above, or information can be appended before responses are broadcasted.

To make the framework interact with platform-specific code, developers will need to implement the Carrot client framework. Currently, only iOS support exists. To see how to do so, visit the carrot-ios repository [https://github.com/carrot-ar/carrot-ios](https://github.com/carrot-ar/carrot-ios)

## The Picnic Protocol

The Picnic Protocol (patent pending) is a set of rules and standards that provide a way for devices to communicate local AR events as well as understand foreign ones. More specifically, however, it relies on both decentralized and centralized network topologies in order to solve the problem of understanding events that happen in foreign coordinate spaces.
The protocol's "handshake" begins by designating the first device to join the session as the primary device. The primary device has two responsibilities:

- It must provide other devices a way to know that they are immediately next to it in physical space, which we'll refer to as the "immediate ping". On iOS, this is achieved by broadcasting iBeacon signals from the primary device.

- It must let the server know what it's current position in physical space is whenever the server asks for it, which the server does by sending a message with a reserved endpoint. At the moment of the immediate ping, the server asks the primary device for its position in physical space. We'll refer to this as TP, or the primary device's transform.

The rest of the devices in a session are referred to as secondary devices. Secondary devices must be able to listen for the immediate ping from the primary device and let the server know that they received this immediate ping by sending it their own position in physical space at that moment in time. We refer to this as TL, or the local transform. 
The state of the environment between a secondary device and the primary device at the moment of the immediate ping is illustrated below.

![figure 1](https://i.imgur.com/yTr9OEg.png)
*The first step of the invention’s handshake, shown from the perspective of the secondary device. TL is the vector reflecting where the secondary device travelled to receive the immediate ping from the primary device. TP reflects where the primary device travelled to send the immediate ping to the secondary device.*

After receiving the immediate ping, a secondary device is considered to be authenticated and ready to interact with other devices in the session. The invention uses the TL and TP relationship between every secondary device and the primary device in order to calculate the primary device’s origin in the secondary device’s coordinate space. This equation, explained in the image below, acts as the bridge between a secondary device and any other authenticated device in the network, whether that be the primary device or another secondary device.

This is the core of protocol. Clients are responsible for being able to send and receive the immediate ping and the server is responsible for maintaining the TL and TP relationship for every authenticated device in the network, as well as converting the locations in messages themselves before broadcasting them to clients.

![figure 2](https://i.imgur.com/IEsfau0.png)
*The calculation of OP, which is the vector resulting from the difference of TL and TP. Visually speaking, OP can be calculated by “walking along” TL and then walking in the opposite direction of TP. Being able to derive OP via this relationship allows the server to convert a coordinate that originated in the coordinate system of a secondary device to one that is now relative to the origin of the primary device. This equation can be applied a second time over in order to do secondary device to secondary device conversions.*

The final step in picnic’s coordinate conversion work is to take the coordinates of a local event, referred to as EL, and convert it to the primary device’s coordinate space. This results in a new vector, EP, which the server populates the outgoing message with before broadcasting it to the primary device. The calculation of EP is illustrated below.

The protocol is a platform-agnostic way of performing coordinate conversion. It was designed specifically for multi-device augmented reality on mobile devices though, and is well tailored for that use case. The only thing it requires devices to have, however,  is the ability to connect to a network. Although Bluetooth and iBeacon technologies were chosen as the way to do inter-device communication in the iOS framework, one can imagine this happening over something like a P2P network instead, for example.

![figure 3](https://i.imgur.com/uRsqDEH.png)
*The calculation of EP, which is EL converted to be relative to the primary device’s origin. The server mutates the message sent by the secondary device who rendered the event at EL, effectively replacing EL with EP. This allows the primary device to take the incoming message and render it as is, without having to even consider where the coordinate originated from. It allows the primary device to treat all incoming messages as if they originated locally.*


## Message Format

Carrot has two message types: request and responses. These are represented by the []byte type.

Requests are created and sent by the client framework to the server framework. Conversely, responses are created and sent by a developer defined controller back to the client framework. The end of a request's path marks the beginning of the corresponding response's path. 

The structure of messages are identical, so the two types of messages represent the opposite directions (and ultimate paths) data travels. Requests and responses are in the form of the following JSON:

	{
		"session_token": "E621E1F8-C36C-495A-93FC-0C247A3E6E5F",
		"endpoint": "echo",
		"payload": {
			"offset": {
				"x": 3.2,
				"y": 1.3,
				"z": 4.0
			},
			"params": {
				"foo": "bar"
			}
		}
	}

## Sending Messages to Carrot

Carrot expects to recieve requests from clients.

In order for a client to send messages to other devices, a route must be defined so the server knows the destination controller and method to handle the incoming requests. Consider this example application snippet:

``` go
// Controller declaration
type ExampleController struct{}

//Controller method implementation
func (c *ExampleController) PrintFooParameterToConsole(req *carrot.Request, br *carrot.Broadcast) {
	fmt.Println(req.Params["foo"])
}

func main() {
  	// Register endpoint by providing endpoint, controller, and method, and if endpoint requires streaming
	carrot.Add("print_foo_param", ExampleController{}, "PrintFooParameterToConsole", true)
}
```

In order to send a message to this endpoint, all we need to do is specify the endpoint in our request. A request to be sent to the endpoint defined above could look something like this :

```
{
	"session_token": "E621E1F8-C36C-495A-93FC-0C247A3E6E5F",
	"endpoint": "print_foo_param",
	"payload": {
		"offset": {
			"x": 3,
			"y": 1,
			"z": 4
		},
		"params": {
			"foo": "bar"
		}
	}
}
```

If this request is sent to the server over the WebSocket connection established, it will create a new instance of the `ExampleController` and route the message to the `PrintFooParameterToConsole` function. Inside the function, we will log the value of the `foo` parameter to the server log. 

Once the request reached its intended controller method, it has reached the end of its life cycle.

## Receiving Messages from Carrot

Clients can expect to receive responses from Carrot. Assuming the request is not explicitly modified in the developer-defined controller, then the response should be exactly the same as the request. 

In the case that the developer wants to modify requests before the devices are broadcasted and recieve these responses, custom responses can be built.

### Responses

There are two different types of responses: default and custom.

#### Default Responses

If the developer wants to forward requests as they are, then they can use the create a default response. The following snippet shows how to create a default response:

``` go

func (c *ExampleController) CreateDefaultResponse(req *carrot.Request, br *carrot.Broadcast) {
	defaultResponse, err := carrot.CreateDefaultResponse(req)
	if err != nil {
		fmt.Println(err)
		return
	}
}

```

As you can see above, the contents of the request received by the controller are dumped into the generating the contents of the new response. Behind the scenes, the JSON representing the message stucture is created and returned as a []byte ready to be broadcasted. Once created, a default reponse cannot be modified. This option is therefore used for sake of brevity in relevant use cases.

#### Custom Responses

If the developer wants to add extra information to the message, then a custom response must be created.  Like the default response, the contents are the request are placed in the response. However, the custom response requires the developer to explictly call more functions to create it. 

The first two functions define the contents for all of the fields that are copied from the request. 

``` go 

func (c *EchoController) EchoExtendable(req *carrot.Request, br *carrot.Broadcast) {
	token := string(req.SessionToken)
	payload, err := carrot.NewPayload(token, req.Offset, req.Params)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := carrot.NewResponse(token, "Echo", payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	...
}

```

The next and most interesting functions to note are the `AddParam` and `AddParams` functions. They allow the developer to append key-value pairs to a custom response.

``` go

func (c *EchoController) EchoExtendable(req *carrot.Request, br *carrot.Broadcast) {
	...
	oneFishTwoFish := ResponseParams{"red": "fish", "blue": "fish"}
	res.AddParam("someKey", "someValue")
	res.AddParams(oneFishTwoFish)
	...
}
	
```

Finally, a custom response must be packaged into a JSON object. 

``` go

func (c *EchoController) EchoExtendable(req *carrot.Request, br *carrot.Broadcast) {
	...
	message, err := res.Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	...
}
	
```

The resulting message is ready to be broadcasted to clients and can no longer be modified.

## Broadcasting Responses

The broadcast module, available in all controller implementations, has a few options for narrowing down which clients to send a message to. Since all clients have a session associated with them, there is a 1-to-1 relationship between sessions and clients. Thus, every client has a `SessionToken` which is accessible within the client and within the session store internal to Carrot. 

#### Broadcasting to all clients

``` go
carrot.Broadcast(/* carrot response  */)
```

#### Broadcasting to a subset of clients

``` go
carrot.Broadcast(/* carrot response */, sessionToken1, sessiontoken2)
```
or
``` go
recipients := []string{sessionToken1, sessionToken2, sessionToken3)
carrot.Broadcast(/* carrot response */, recipients)
```

One way to make the best use of this feature is to keep sessions associated with users in a datastore connected to carrot. Then, a simple query can return a set of sessions that should be sent a response. Here is some pseudocode demonstrating this:
``` go
func (c *ExampleController) SendHelloToAll(r *carrot.Request, b *carrot.Broadcast) {
	/* build up a response here */
	/* database call to get a list of session tokens based on a query */
	b.Broadcast(/* response */, /* array with session tokens */)
}
```

Once the response reaches its intended recipient(s), it has reached the end of its life cycle.

## Sessions

Maintaining the state of clients are done using sessions inside of Carrot. Due to the shared, concurrent access of sessions throughout the lifecycle of a request, they are stored isnide of Golang's `sync.Map`. An interface is provided such that future extensibility could be easily integrated into Carrot to allow session storage within in-memory data stores like Redis. 

Within Carrot, the session store is maintained using a singleton pattern and within any point of carrot the current state of sessions is accessible by calling the `NewDefaultSessionManager()`, which returns a pointer to the `SessionStore` interface. See the GoDoc for the `SessionStore` interface and `DefaultSessionStore` struct for more details.

### Resuming a Session

To be implemented
