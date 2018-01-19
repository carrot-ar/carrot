// when building response messages:

// 	event //unconverted event coordinates that may be recalculated twice depending on a particular broadcast's sender and recipient's roles
// 	var e_l //obtain event in primary device's coordinates (first potential calcuation)
// 	var e_p //obtain event in secondary device's coordinates (second potential calcuation)

// 	if sender is a secondary device  {  //establish primary perspective

// 		//calculate object in primary device's coordinates
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
// 		} else { //device == primary device || self to self broadcast
// 			e_p = e_l
// 		}
// 		// broadcast e_p as the offset in the response for this recipient device
		

