package main

import "time"

type Database struct {
    Name      string    `json:"name"`
    running   bool      `json:"running"`
    started   time.Time `json:"due"`
    port      int       `json:"port"`
    host      string    `json:"host"`
    size      int       `json:size`

}

type Databases []Database
