package products

import "time"

type Product struct {
	Id         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `json:"name"`
	Categories []string  `json:"categories"`
	Added      time.Time `json:"added"`
}
