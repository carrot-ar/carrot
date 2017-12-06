<p align="center">
<img src="https://github.com/carrot-ar/carrot-ios/wiki/resources/Carrot@2x.png" alt="Carrot" width="300">
</p>

<p align="center">
<a href="https://travis-ci.org/carrot-ar/carrot"><img src="https://travis-ci.org/carrot-ar/carrot.svg?branch=master" alt="build status"></a>
<a href=""><img src="https://codecov.io/gh/carrot-ar/carrot/branch/master/graph/badge.svg" alt="code coverage"></a>
</p>

Carrot is an easy-to-use, real-time framework for building multiplayer applications in Augmented Reality. Currently, not many AR frameworks exist with multiplayer in mind. There are a few reasons for this, with the most important being the difficulty of resolving location to an acceptable degree of accuracy with traditional GPS based coordinates. This is where Carrot flourishes.

By implementing the Picnic Protocol into the server and client's respective frameworks, we have decreased the error size for location resolution from 10-65 meters with GPS down to less than one foot. This enables developers (i.e. you) to focus on creating applications with rich content and need not worry about the finer details such as cross-device accuracy and networking.

## Response Groups
The broadcast module, available in all controller implementations, has a few options for narrowing down which clients to send a message to. Since all clients have a session associated with them, there is a 1-to-1 relationship between sessions and clients. Thus, every client has a `SessionToken` which is accessible within the client and within the session store internal to Carrot. 

#### Broadcasting to all clients
```
carrot.Broadcast(/* carrot response message */)
```

#### Broadcasting to a subset of clients
```
carrot.Broadcast(/* carrot response message */, sessionToken1, sessiontoken2)
```
or
```
recipients := []string{sessionToken1, sessionToken2, sessionToken3)
carrot.Broadcast(/* carrot response message */, recipients)
```

One way to make the best use of this feature is to keep sessions associated with users in a datastore connected to carrot. Then, a simple query can return a set of sessions that should be sent a response. Here is some pseudocode demonstrating this:
```
func (c *ExampleController) SendHelloToAll(r *carrot.Request, b *carrot.Broadcast) {
	/* build up a response message here */
	/* database call to get a list of session tokens based on a query */
	b.Broadcast(/* response message */, /* array with session tokens */)
}
```

## Building an application with Carrot

Building applications on Carrot is incredibly simple. Check out this simple echo application that echos text input from one device into the AR space of all connected devices: 

```
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

	// Register endpoint connection here by providing endpoint, controller, and method, respectively
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

To make the framework interact with platform-specific code, developers will need to implement the Carrot client framework. Currently, only iOS support exists. To see how to do so, visit the carrot-ios repository [https://github.com/carrot-ar/carrot-ios](https://github.com/carrot-ar/carrot-ios)

## Message Format
Carrot has two message types: request and responses.

Requests are sent by the client framework to the server framework. Conversely, responses are sent by a developer defined controller back to the client framework. The structure of messages are identical, so the two types of messages represent the opposite directions (and ultimate paths) data travels.

Requests and responses are in the form of the following JSON:

	{
		"session_token": "E621E1F8-C36C-495A-93FC-0C247A3E6E5F",
		"endpoint": "test_endpoint",
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

Depending on the the intent of the message, the payload may be completely or partially empty. These scenarios are explained below.

** BELOW THIS IS OUT OF DATE! **

## Sending Messages

In order to send messages, a route must be defined so the server knows the destination controller and method to handle the incoming message. Consider this example application:
```
// Controller definition
type ExampleController struct {}

// Method implementation for a controller
func (c *ExampleController) PrintFooParameterToConsole(req *buddy.Request) {
  fmt.Println(req["foo"])
}

func main() {
  // Route registration in the form of endpoint, controller, method, and whether the endpoint requires streaming
  buddy.Add("print_foo_param", ExampleController{}, "PrintFooParameter", false)
}
```

In order to send a message to this endpoint, all we need to do is specify the endpoint in our message. Messages for carrot take on the following form:
```
{
	"session_token": "KjIQhKUPNrvHkUHv1VySBg==",
	"endpoint": "print_foo_param",
	"origin": {
		"longitude": 45.501689,
		"latitude": -73.567256
	},
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

If this message is sent to the server over the WebSocket connection established, it will create a new instance of the `ExampleController` and route the message to the `PrintFooParam` function. Inside the function, we will log the value of the `foo` parameter to the server log. 

## Receiving Messages

At the moment, all messages are echo'd back to every client. This will change once the `responder` module is complete. 

### Responses

## Sessions

## New Session with client secrets disabled

To connect to the server, connect to the WebSocket url the server is running on. For this example, `localhost:8080` will be used and a Ruby WebSocket client will be used for demo purposes.

```
ws = WebSocket::Client::Simple.connect 'ws://localhost:8080/ws'
```

Once the client connects, a welcoming message consisting of the `SessionToken` will be sent. It will look like this:
```
KjIQhKUPNrvHkUHv1VySBg==
```
Save this token. It will be required to be attached to every message sent to the server.

From this point on you can begin sending/receiving messages to the server. 

## Resuming a Session
When a WebSocket connection is closed, the session state is maintained for a period of time determined by the application configuration. 

tbd



