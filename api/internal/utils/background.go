package utils

import (
	"log"
	"sync"
)

func Background(fn func()) {
	var wg sync.WaitGroup
	wg.Go(func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()

		fn()
	})
}
