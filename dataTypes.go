package main

type DataNested struct {
	Hour      int     `json:"hour"`
	Price_sek float32 `json:"price_eur"`
	Price_eur float32 `json:"price_sek"`
	Kmeans    int     `json:"kmeans"`
}

type Data struct {
	Date string `json:"date"`
	Se1  []DataNested
	Se2  []DataNested
	Se3  []DataNested
	Se4  []DataNested
}
