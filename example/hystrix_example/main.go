package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
	"time"
)

func main() {
	hystrix.ConfigureCommand("koala_rpc", hystrix.CommandConfig{
		Timeout:               10,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	for {
		err := hystrix.Do("get_baidu", func() error {
			//talk to other service
			_, err := http.Get("https://www.baidu.com/")
			if err != nil {
				fmt.Println("get error")
				return err
			}

			return nil
		}, func(err error) error {
			// 应急预案
			fmt.Printf("get an error, handle it, err : %v\n", err)
			return err
		})
		if err == nil {
			fmt.Println("request succ")
		}
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Second * 2)
}
