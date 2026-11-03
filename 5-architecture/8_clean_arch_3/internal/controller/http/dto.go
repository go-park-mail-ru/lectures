package httpapi

import "lectures/cleanarch3/internal/domain"

type userRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func toResponse(u domain.User) userResponse {
	return userResponse{ID: u.ID(), Name: u.Name(), Email: u.Email()}
}
