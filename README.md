## Address Book API (Go)

I built a small but production-minded REST API to manage contacts using Go’s standard library only. It demonstrates clean separation of concerns, JSON encoding/decoding, basic validation, middleware (logging + CORS), and in-memory storage. I avoided external frameworks on purpose to understand how `net/http` works under the hood.

### Why I built it this way
- **Minimal dependencies**: I used only the standard library to keep the surface area small and highlight fundamentals (routing, handlers, middleware, JSON I/O).
- **Separation of concerns**: I split the code into `handlers`, `models`, `storage`, and `utils` to keep responsibilities clear.
- **Explicit routing**: I registered routes with `http.HandleFunc` and handled methods explicitly to be clear about the HTTP surface area.
- **Demonstrate middleware**: I wrapped the server with a request logger and CORS middleware to show cross-cutting concerns.

### Tech stack
- **Language**: Go (`go 1.24.x`) using the standard library (`net/http`, `encoding/json`, etc.)
- **Persistence**: In-memory slice for simplicity (no external DB)

### Project structure
```
AddressBook/
  main.go                   # HTTP server, route registration, middleware wiring
  applicationapi.json       # Postman collection for manual testing
  go.mod
  handlers/
    contacts.go             # All contact CRUD + search handlers + hello
    middleware.go           # Logger and CORS middleware
  models/
    contact.go              # Contact struct + JSON tags
  storage/
    memory.go               # In-memory data store and ID generator
  utils/
    response.go             # Reusable JSON error writer
```

### Data model
- **`models.Contact`**: `id`, `name`, `email`, `phone`, `address` with JSON tags for clean API payloads.

### Storage strategy
- **In-memory slice**: `storage.Contacts` holds all contacts; `storage.NextID` auto-increments IDs.
- This keeps onboarding fast for demos. It’s intentionally not concurrent-safe or persistent; see improvements below for how I’d evolve it.

### HTTP server and routing (how it works)
- **Server bootstrap (`main.go`)**
  - Registers routes with `http.HandleFunc` using the default mux.
  - Wraps the mux with `MiddlewareLogger` and then `MiddlewareCORS` before starting on `:8080`.
  - Routes:
    - `/` → `HelloHandler`
    - `/contacts` → `GET` list, `POST` create
    - `/contacts/` → `GET` by id, `PUT` update, `DELETE` delete (path parsing via `strings.Split`)
    - `/contacts/search` → `GET` search by name/email/phone (query `q`)

- **Handlers (`handlers/contacts.go`)**
  - `HelloHandler`: simple health/demo endpoint.
  - `GetcontactsHandler` (GET `/contacts`): returns all contacts as JSON.
  - `CreateContactHandler` (POST `/contacts`):
    - Decodes JSON body into `Contact`.
    - Validates `name` and `email` are provided.
    - Assigns a new ID, appends to in-memory store, returns the created contact.
  - `GetContactByIDHandler` (GET `/contacts/{id}`):
    - Parses `id` from path, finds the contact, returns 404 if missing.
  - `UpdateContactHandler` (PUT `/contacts/{id}`):
    - Validates input, replaces the contact at index; returns 404 if missing.
  - `DeleteContactHandler` (DELETE `/contacts/{id}`):
    - Removes the contact from the slice; returns 204 No Content on success.
  - `SearchContactHandler` (GET `/contacts/search?q=...`):
    - Case-insensitive substring match on `name`, `email`, `phone`.

- **Middleware (`handlers/middleware.go`)**
  - `MiddlewareLogger`: logs method, path, and request duration.
  - `MiddlewareCORS`: sets permissive CORS headers and short-circuits `OPTIONS` with 204.

- **Utilities (`utils/response.go`)**
  - `WriteJSONError`: standardizes error JSON shape and headers.

### API reference (what the API does)
- **GET `/contacts`**: list all contacts
  - Response: `200 OK` with `[]Contact`.

- **POST `/contacts`**: create a contact
  - Body JSON:
    ```json
    {"name":"Test User","email":"test@example.com","phone":"1234567890","address":"Somewhere"}
    ```
  - Validates `name` and `email`.
  - Response: `202 Accepted` with created `Contact` (see notes on improvements for status code).

- **GET `/contacts/{id}`**: fetch a contact by numeric `id`
  - Errors: `400` if `id` isn’t numeric; `404` if not found.

- **PUT `/contacts/{id}`**: update a contact
  - Requires `name` and `email`.
  - Errors: `400` invalid input; `404` if not found.

- **DELETE `/contacts/{id}`**: delete a contact
  - Response: `204 No Content` on success; `404` if not found.

- **GET `/contacts/search?q=term`**: search by name/email/phone (substring, case-insensitive)
  - Errors: `400` when `q` missing.

### Example requests
```bash
# List
curl -s http://localhost:8080/contacts | jq

# Create
curl -s -X POST http://localhost:8080/contacts \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","phone":"1234567890","address":"Somewhere"}' | jq

# Get by id
curl -s http://localhost:8080/contacts/1 | jq

# Update
curl -s -X PUT http://localhost:8080/contacts/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated","email":"updated@example.com","phone":"000","address":"New"}' | jq

# Delete
curl -i -X DELETE http://localhost:8080/contacts/1

# Search
curl -s "http://localhost:8080/contacts/search?q=hard" | jq
```

### Error handling and validation
- Consistent error envelope: `{ "error": "message" }` via `utils.WriteJSONError`.
- Input validation: `name` and `email` required on create/update.
- Path/ID parsing: numeric `id` enforced with appropriate `400` responses.
- Not founds: `404` with clear message.

### Logging and CORS
- All requests are timed and logged: method, path, and latency.
- CORS is open (`*`) with `GET, POST, PUT, DELETE` allowed and `Content-Type` header; `OPTIONS` is short-circuited with `204` to support browsers.

### How to run locally
1. Install Go 1.24+
2. From the project root, run:
   ```bash
   go run main.go
   ```
3. Server listens on `http://localhost:8080`.
4. Optional: import `applicationapi.json` into Postman to test endpoints quickly.

### Why This?
- To design a clean, minimal web API in Go using only the standard library.
- To apply proper separation of concerns and introduce middleware for cross-cutting functionalities.
- To apply proper separation of concerns and introduce middleware for cross-cutting functionalities.


