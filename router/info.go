package router

import "server/utils"

func CountConn() int {
	count := 0
	for _, j := range utils.LiveConn {
		if j {
			count++
		}
	}
	return count
}
