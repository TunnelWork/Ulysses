package main

// Form: submitted by directly to the API server with POST method
// Directly generated through: c.BindJSON(&form)
type (
	FormAuthorize struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)
