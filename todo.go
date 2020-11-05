package main

import (
    "crypto/sha256"
    "fmt"
    "time"
)

type Todo struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Done        bool      `json:"done"`
    Created     time.Time `json:"created"`
}

func NewTodo(name, description string) *Todo {
    t := Todo{
        Name:        name,
        Description: description,
        Done:        false,
        Created:     time.Now(),
    }
    t.ID = t.HashCode()

    return &t
}

func (t Todo) HashCode() string {
    inline := fmt.Sprintf("%s-%s-%v-%d", t.Name, t.Description, t.Done, t.Created.Unix())
    return fmt.Sprintf("%x", sha256.Sum256([]byte(inline)))[:20]
}