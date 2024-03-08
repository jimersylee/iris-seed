package sqlengine

import "log"

func assert(true bool) {
	if !true {
		log.Println("assert error")
	}

}
