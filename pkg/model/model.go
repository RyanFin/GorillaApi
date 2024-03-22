package model

type Car struct {
	UUID  string `json:"uuid"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
	Color string `json:"color"`
	Price int    `json:"price"`
}
