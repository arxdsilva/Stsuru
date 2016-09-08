package mgo



func Insert() {
  session, err := mgo.Dial("localhost")
	defer session.Close()
	err = session.DB("tsuru").C("links").Insert(l)
	return err
}

func Delete() {

}

func Update() {

}

func FindOne() {
  dbData := lines{}
	session, err := mgo.Dial("localhost")
	checkError(err)
	defer session.Close()
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": dbHash}).One(&dbData)
	return dbData.Link, err
}

func FindAll() {

}
