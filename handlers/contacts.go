package handlers

import(
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
	"addressbook/models"
	"addressbook/storage"
	"addressbook/utils"
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
		utils.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newContact models.Contact
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		utils.WriteJSONError(w, "Invalid JSON", http.StatusBadRequest)
	}

	if newContact.Name == "" || newContact.Email == "" {
		utils.WriteJSONError(w, "Name and Email Required", http.StatusBadRequest)
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
		utils.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		utils.WriteJSONError(w, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil{
		utils.WriteJSONError(w, "ID must be a number", http.StatusBadRequest)
	}

	for _, c := range storage.Contacts {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	utils.WriteJSONError(w, "Contact not found", http.StatusNotFound)
}


func UpdateContactHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut {
		utils.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		utils.WriteJSONError(w, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, "ID must be a number", http.StatusBadRequest)
		return
	}

	var updatedContact models.Contact
	err = json.NewDecoder(r.Body).Decode(&updatedContact)
	if err != nil {
        utils.WriteJSONError(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	if updatedContact.Name == "" || updatedContact.Email == "" {
        utils.WriteJSONError(w, "Name and Email are required", http.StatusBadRequest)
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

	utils.WriteJSONError(w, "Contact not found", http.StatusNotFound)
}


func DeleteContactHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete{
		utils.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		utils.WriteJSONError(w, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, "ID must be integer", http.StatusBadRequest)
		return
	}

	for i, c := range storage.Contacts {
		if c.ID == id {
			storage.Contacts = append(storage.Contacts[:i],storage.Contacts[i+1:]... )
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	utils.WriteJSONError(w, "No Contact Found", http.StatusNotFound)

}

func SearchContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteJSONError(w, "Query paramater q required", http.StatusBadRequest)
	}

	query = strings.ToLower(query)
	var results []models.Contact
	for _, c := range storage.Contacts {
		if strings.Contains(strings.ToLower(c.Name), query) || strings.Contains(strings.ToLower(c.Email), query) || strings.Contains(strings.ToLower(c.Phone), query) {
			results = append(results, c)
		}
	}

	w.Header().Set("Content-Typr", "application/json")
	json.NewEncoder(w).Encode(results)
}

