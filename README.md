<img src="https://github.com/carrot-ar/carrot-ios/wiki/resources/Carrot@2x.png" alt="Carrot" width="300">

[![Build Status](https://travis-ci.org/carrot-ar/carrot.svg?branch=master)](https://travis-ci.org/carrot-ar/carrot) 
[![codecov](https://codecov.io/gh/carrot-ar/carrot/branch/master/graph/badge.svg)](https://codecov.io/gh/carrot-ar/carrot)

Carrot is an easy-to-use, real-time framework for building multiplayer applications in Augmented Reality. Currently, not many AR frameworks exist with multiplayer in mind. There are a few reasons for this, with the most important being the difficulty of resolving location to an acceptable degree of accuracy with traditional GPS based coordinates. This is where Carrot flourishes. 

By implementing the [Picnic Protocol](https://github.com/carrot-ar/carrot-ios/wiki/The-Picnic-Protocol%E2%84%A2) into the server and client's respective frameworks, we have decreased the error size for location resolution from 10-65 meters with GPS down to less than one foot. This enables developers (i.e. you) to focus on creating applications with rich content and need not worry about the finer details such as cross-device accuracy and networking. 

## Building an application with Carrot

Building applications on Carrot is incredibly simple. Check out this simple echo application that echos text input from one device into the AR space of all connected devices: 

```
package main

import (
  "github.com/carrot-ar/carrot"
)

type EchoController struct{}

func (c *PingController) Echo(req *carrot.Request, broadcast *carrot.Broadcast) {
	responseText := req.Params["foo"]
	response := carrot.Response(responseText)
	broadcast.Send(response)
}

func main() {

  // Register endpoints here in the order of endpoint, controller, method
  carrot.Add("echo", EchoController{}, "Echo")

  // Run the server
  carrot.Run()
}
```

## Message Format
Carrot has two message types: request and responses. 

### Requests

### Responses

## Sessions



