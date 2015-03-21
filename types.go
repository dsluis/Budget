package main

//todo: store password hash instead
type User struct {
	Username string
	Password string
}

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