package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)

type Note struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
}

type NoteResource struct {
	Note Note `json:"note"`
}

type NotesResource struct {
	Notes []Note `json:"notes"`
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {

	var noteResource NoteResource

	err := json.NewDecoder(r.Body).Decode(&noteResource)
	if err != nil {
		panic(err)
	}

	note := noteResource.Note
	// get a new id
	obj_id := bson.NewObjectId()
	note.Id = obj_id
	//insert into document collection
	err = collection.Insert(&note)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Inserted new Note %s with name %s", note.Id, note.Name)
	}

	j, err := json.Marshal(NoteResource{Note: note})
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func NotesHandler(w http.ResponseWriter, r *http.Request) {

	var notes []Note

	iter := collection.Find(nil).Iter()
	result := Note{}
	for iter.Next(&result) {
		notes = append(notes, result)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(NotesResource{Notes: notes})
	if err != nil {
		panic(err)
	}
	w.Write(j)
	log.Println("Provided json")

}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	// Get id from the incoming url
	vars := mux.Vars(r)
	id := bson.ObjectIdHex(vars["id"])

	// Decode the incoming note json
	var noteResource NoteResource
	err = json.NewDecoder(r.Body).Decode(&noteResource)
	if err != nil {
		panic(err)
	}

	// Update the document
	err = collection.Update(bson.M{"_id": id},
		bson.M{"name": noteResource.Note.Name,
			"_id":         id,
			"description": noteResource.Note.Description,
		})
	if err == nil {
		log.Printf("Updated Note %s name to %s", id, noteResource.Note.Name)
	} else {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id := vars["id"]

	// Remove it from database
	err = collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Printf("Could not find Note %s to delete", id)
	}
	w.WriteHeader(http.StatusNoContent)
}
func main() {
	log.Println("Starting Server 2")

	r := mux.NewRouter()
	r.HandleFunc("/api/notes", NotesHandler).Methods("GET")
	r.HandleFunc("/api/notes", CreateNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", UpdateNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")
	http.Handle("/api/", r)

	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Starting mongo db session")
	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("NotesDB").C("notes")

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
