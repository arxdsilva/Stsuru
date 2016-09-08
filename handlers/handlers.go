package handlers

import (
  "net/http"
	"net/url"
)

func AddLink(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  link := r.Form["user_link"][0]

  // checking the URL
  isURL := govalidator.IsURL(link)
  if isURL != true {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }
  validateURL := govalidator.IsRequestURL(link)
  if validateURL != true {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }

  path := "http://localhost:8080/"
  // URL hashing
  linkShort, dbHash := hash(link, path)

  l := &lines{Link: link, Short: linkShort, Hash: dbHash}

  _, err := findOne(dbHash)
  if err == nil {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }

  err = insert(l)
  checkError(err)

  http.Redirect(w, r, "/", http.StatusFound)
}

func Home(w http.ResponseWriter, r *http.Request)  {
  Data := []lines{}
  session, err := mgo.Dial("localhost")
  defer session.Close()
  checkError(err)
  c := session.DB("tsuru").C("links")
  err = c.Find(bson.M{}).All(&Data)
  checkError(err)
  t, err := template.ParseFiles("tmpl/index.html")
  checkError(err)

  t.Execute(w, Data)
}

func css(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/out/", http.FileServer(http.Dir("out/")))
}

// RemoveLink searches db for a certain link & removes It if It exists
func RemoveLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]

	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": idInfo})
	checkError(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

// LinkSolver ...
func LinkSolver(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInfo := id["id"]

	l, err := findOne(idInfo)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.Redirect(w, r, l, http.StatusFound)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

func hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}
