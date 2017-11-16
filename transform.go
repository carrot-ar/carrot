package carrot

import (
	log "github.com/sirupsen/logrus"
)

//internal controller used to update transforms of primary and secondary devices
type CarrotTransformController struct {
	sessions	SessionStore
}

func (c *CarrotTransformController) Transform(req *Request, broadcast *Broadcast) {
	if c.sessions == nil {
		c.sessions = NewDefaultSessionManager()
	}
	primaryToken, err := c.sessions.GetPrimaryDeviceToken()
	if err != nil {
		log.Errorf("There was an error retrieving the primary device token in transform.go")
	}
	session, err := c.sessions.Get(req.SessionToken)
	if err != nil {
		log.Errorf("There was an error retrieving the session in transform.go")
	}
	if req.SessionToken != primaryToken { //store T_L for the secondary device and request T_P from the primary device
		session.T_L = req.Offset
		//broadcast response to primary device that has primaryDevice token, this endpoint, empty params
		res, err := getT_PFromPrimaryDeviceRes(string(req.SessionToken))
		if err != nil {
			log.Errorf("There was an error creating a response to retrieve T_P in transform.go")
		}
		broadcast.Broadcast(res, string(primaryToken))
	} else { //store T_P from primary device
		c.sessions.Range(func(t, session interface{}) bool {
			s := session.(*Session)
			if s.T_P == nil && s.T_L != nil && s.Token != primaryToken {
				s.T_P = req.Offset
			}
			if s.T_P != nil && s.T_L != nil {
				log.Printf("session w/ token %v has filled transforms and is ready to broadcast to others!\n", s.Token)
			}
			return true
		})
	}
}