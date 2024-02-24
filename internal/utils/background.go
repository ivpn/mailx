package utils

import (
	"log"
	"sync"
)

func Background(fn func()) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
		}()

		fn()
	}()
}
