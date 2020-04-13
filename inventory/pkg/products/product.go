package products

import "time"

type Product struct {
	Id         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `json:"name"`
	Categories []string  `json:"categories"`
	Added      time.Time `json:"added"`
}

type Metadata struct {
	ExpiresOn time.Time
	Peer      string
}

type RemoteProduct struct {
	Product Product
	Meta    Metadata
}

func OnlyProducts(remotes []RemoteProduct) []Product {
	var res []Product

	for _, remote := range remotes {
		res = append(res, remote.Product)
	}

	return res
}
