<p align="center">
<img src="https://github.com/carrot-ar/carrot-ios/wiki/resources/Carrot@2x.png" alt="Carrot" width="300">
</p>
<p align="center">
<a href="https://travis-ci.org/carrot-ar/carrot"><img src="https://travis-ci.org/carrot-ar/carrot.svg?branch=master" alt="build status"></a>
<a href=""><img src="https://codecov.io/gh/carrot-ar/carrot/branch/master/graph/badge.svg" alt="code coverage"></a>
</p>



## Building an application with Carrot
```
package main

import (
  "github.com/carrot-ar/carrot"
)

type PingController struct{}

func (c *PingController) Ping(req *carrot.Request, res *carrot.Responder) {
	res.broadcast <- []byte("Pong!")
}

func main() {

  // Register endpoints here in the order of endpoint, controller, method,
  // and whether the route will accept streaming
  carrot.Add("ping", PingController{}, "Ping", false)

  // Run the server
  carrot.Run()
}
```




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

## Controller Types
Controllers take on two forms. EventControllers and StreamControllers. 

### EventController
EventControllers handle one time events such as placing an object. When a request is sent to an EventController, a new instance of the controller is instantiated, the request is routed to the correct method, and once the request is finished the controller is dereferenced.

### StreamController
StreamControllers handle persistent events such as drawing a line or moving an object in real time. When a request is sent to a StreamController, a new instance of the controller is instantiated similar to the EventController. However, after the request is acted on by the correct method, the controller is maintained in a map from `SessionToken` to `StreamController`. This allows you to call the same route multiple times in a row to the same instance of a controller. Some benefits of this are a performance boost (no need to instantiate a controller, just a simple lookup is required) as well as the ability to set fields in the struct and have them persist between requests. 

Currently, a StreamController will remain open indefinitely. 

