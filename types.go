package main

type ViewData struct {
	Feedback  string
	ViewModel interface{}
}

type Transaction struct {
	Description string
	From string
	To string
	Amount float32 
}