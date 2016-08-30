package main

func requestTest() {
	req, _ := httptesting.NewRequest("GET", "/", nil)
	w := httptesting.NewRecorder()
	Home(w, req)
}
