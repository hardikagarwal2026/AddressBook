package storage

import "addressbook/models"

//slice of contact. containing all the contacts, of type COntact in the contact struct
// ham basically saare contacts isi mae store karwaenge

var Contacts = []models.Contact{
	{ID: 1, Name: "Hardik", Email: "hardikagarwaljpr@gmail.com", Phone: "7424940418", Address: "Jaipur"},
	{ID: 2, Name: "Saumya Sharma", Email: "saumyasharma@gmail.com", Phone: "9352829218", Address: "Jaipur"},
}

var NextID = 3;




