package entities

// Store connected clients (maps userID -> socket UUID)
var SocketClients = make(map[string]string)
