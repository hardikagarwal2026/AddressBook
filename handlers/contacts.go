package handlers

import(
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
	"addressbook/models"
	"addressbook/storage"
)

//http.Responsewriter , response ko apne hisab se control krne deta hai
// w - response writer(how we reply!!)
// r - incoming request (method, path, headers, body)
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"Hello World")
}

func GetcontactsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.Contacts)
}

func CreateContactHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newContact models.Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	if newContact.Name == "" || newContact.Email == "" {
		http.Error(w, "Name and Email Required", http.StatusBadRequest)
		return
	}

	newContact.ID = storage.NextID
	storage.NextID++;
	storage.Contacts = append(storage.Contacts, newContact)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(newContact)
}

func GetContactByIDHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil{
		http.Error(w, "ID must be a number", http.StatusBadRequest)
	}

	for _, c := range storage.Contacts {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}


func UpdateContactHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	var updatedContact models.Contact
	err = json.NewDecoder(r.Body).Decode(&updatedContact)
	if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	if updatedContact.Name == "" || updatedContact.Email == "" {
        http.Error(w, "Name and Email are required", http.StatusBadRequest)
        return
    }

	for i, c := range storage.Contacts {
		if c.ID == id {
			updatedContact.ID = id
			storage.Contacts[i] = updatedContact
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedContact)
			return

		}
	}

	http.Error(w, "Contact not found", http.StatusNotFound)
}


func DeleteContactHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be integer", http.StatusBadRequest)
		return
	}

	for i, c := range storage.Contacts {
		if c.ID == id {
			storage.Contacts = append(storage.Contacts[:i],storage.Contacts[i+1:]... )
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "No Contact Found", http.StatusNotFound)

}


