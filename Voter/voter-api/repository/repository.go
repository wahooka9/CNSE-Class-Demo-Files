package repository

import (
    "context"
    "fmt"
    "os"
    "github.com/redis/go-redis/v9"
    "errors"
)


func createReddisClient() (*redis.Client, error) {
    redisURL := os.Getenv("REDIS_URL")
    fmt.Println(redisURL)
    if redisURL == "" {
        redisURL = "localhost:6379"
    }
	LocalReddisCilent := redis.NewClient(&redis.Options{
        Addr:     redisURL, // Replace with the appropriate host and port
        Password: "",              // No password if you haven't set one
        DB:       0,               // Default DB
    })

    ctx := context.Background()
    _, err := LocalReddisCilent.Ping(ctx).Result()
    if err != nil {
        fmt.Println("Error pinging Redis:", err)
        return LocalReddisCilent, errors.New("Error Saving to DB")
    }
    return LocalReddisCilent, nil
}

func SetValueForKey(key string, data []byte) bool {
	LocalReddisCilent, err := createReddisClient()
	if err != nil {
		fmt.Println(err)
		return false
	}

    ctx := context.Background()
	err = LocalReddisCilent.Set(ctx, key, data, 0).Err()
    if err != nil {
        fmt.Println("Error setting key:", err)
        return false
    }
    return true
}

func GetValueForKey(key string) (string, error) {
    fmt.Println("start work")
    LocalReddisCilent, err := createReddisClient()
    if err != nil {
        fmt.Println(err)
    }
fmt.Println("halfway work" , key)
    ctx := context.Background()
    val, err := LocalReddisCilent.Get(ctx, key).Result()
    if err != nil {
        fmt.Println(err)
        //panic(err)
        
    }
    fmt.Println(key, val)
    fmt.Println("work work")
    return val, nil
}