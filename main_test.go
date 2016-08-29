package main

func requestTest() {
	req, _ := http.NewRequest("GET", "/", nil)
	w := http.NewRecorder()
	Home(w, req)
}
