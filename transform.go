package carrot

import (
	log "github.com/sirupsen/logrus"
)

//based on Picnic Protocol requirements
type transforms struct {
	T_L *offset //local transform
	T_P *offset //primary device transform
}

type dictionary map[SessionToken]transforms

//internal controller used to update primary and secondary devices
type CarrotTransformController struct {
	dict		dictionary
	sessions	SessionStore
}

func (c *CarrotTransformController) Transform(req *Request, broadcast *Broadcast) {
	if c.dict == nil { //this req's device is the first and therefore primary device
		c.dict = make(map[SessionToken]transforms)
		c.sessions = NewDefaultSessionManager()
	}
	primaryToken, err := c.sessions.GetPrimaryDeviceToken()
	if err != nil {
		log.Errorf("There was an error retrieving the primary device token in Transform")
	}
	//request and update primary device transform
	if req.SessionToken == primaryToken {
		c.dict[primaryToken] = transforms{
			T_P:	req.Offset,
			T_L:	nil,
		}
	} else {
	//request and update secondary device transforms
		c.dict[req.SessionToken] = transforms {
			T_P:	c.dict[primaryToken].T_P,
			T_L:	req.Offset,
		}
	}
}