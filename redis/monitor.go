//redis health monitor
//contionuosly check and updates the health 
package redis

import (
	"context"
	"log"
	"time"
	"api-gateway/middleware"
)

func StartHealthMonitor() {		

	go func() {

		for {

			err := Client.Ping(
				context.Background(),
			).Err()

			if err != nil {

				if RedisHealthy.Load() {

					middleware.GlobalReliabilityMetrics.RedisFailureCount.Add(1)
					log.Println(
						"Redis became unavailable",
					)
				}

				RedisHealthy.Store(false)

			} else {

				if !RedisHealthy.Load() {

					middleware.GlobalReliabilityMetrics.RedisRecoveryCount.Add(1)
					log.Println(
						"Redis recovered",
					)
				}

				RedisHealthy.Store(true)
			}

			time.Sleep(
				5 * time.Second,		//every 5 seconds, monitor
			)							//does PING redis
		}
	}()
}