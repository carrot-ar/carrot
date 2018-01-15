// when building response messages:
// 	var e_l, e_p //obtain event in primary device's coordinates
// 	if sender is a secondary device  {  //establish primary perspective

// 		//calculte object in primary device's coordinates
// 			find primary device
// 			e_l = getE_L(sender, primaryDevice, event)

// 	} else { //sender is primary device & already in perspective
// 		e_l = event
// 	}

// 	//transform and/or broadcast primary device's event 
// 	for the list of recipient devices
// 		if recipient != primaryDevice && sender != recipient {
// 			e_p = getE_P(recipient, e_l)
// 			//calculate e_p from primary to secondary device
// 				//e_p = e_l - o_p
// 					//where o_p = t_p - t_l
// 		} else {//device == primary device || self to self broadcast
// 			e_p = e_l
// 		}
// 		broadcast e_p as the offset in the response for this device
		

