package main

import (
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis/v7"
    "time"
)

type TodoCache struct {
    Client *redis.Client
}

func CreateTodoCache() *TodoCache {
    redisHost := GetEnvOfDefault(RedisHostKey, "localhost")
    redisPort := GetEnvOfDefault(RedisPortKey, "6379")
    redisPass := GetEnvOfDefault(RedisPassKey, "")

    client := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
        Password: redisPass,
        DB: 0,
    })

    return &TodoCache{
        Client: client,
    }
}

func (cache TodoCache) Get(k string) (*Todo, error) {
    payload, err := cache.Client.Get(k).Bytes()
    if err != nil {
        return nil, err
    }

    todo := &Todo{}
    if err := json.Unmarshal(payload, todo); err != nil {
        return nil, err
    }

    return todo, nil
}

func (cache TodoCache) GetAll() ([]Todo, error) {
    var (
    	cursor uint64
    	keys []string
    	err error
        todos []Todo
    )

    keys, cursor, err = cache.Client.Scan(cursor, "*", 100).Result()
    if err != nil {
        return todos, err
    }

    for _, key := range keys {
        todo, _ := cache.Get(key)
        todos = append(todos, *todo)
    }

    return todos, nil
}

func (cache TodoCache) Set(k string, todo Todo) error {
    payload, err := json.Marshal(todo)
    if err != nil {
        return err
    }

    if err := cache.Client.Set(k, string(payload), 24 * time.Hour).Err(); err != nil {
        return err
    }

    return nil
}

func (cache TodoCache) Del(k string) error {
    return cache.Client.Del(k).Err()
}