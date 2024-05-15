// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
}

type CreateDeploymentRequest struct {
	Name    string                 `json:"name"`
	Image   string                 `json:"image"`
	EnvVars map[string]interface{} `json:"env_vars,omitempty"`
}

type DeleteDeploymentRequest struct {
	ID string `json:"id"`
}

type Deployment struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Image       string                 `json:"image"`
	Status      *string                `json:"status,omitempty"`
	ContainerID *string                `json:"container_id,omitempty"`
	EnvVars     map[string]interface{} `json:"env_vars,omitempty"`
}

type GetDeploymentRequest struct {
	ID string `json:"id"`
}

type InviteCompanyUsersResponse struct {
	Message string  `json:"message"`
	Users   []*User `json:"users,omitempty"`
}

type ListDeploymentsRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LoginRequest struct {
	Email string `json:"email"`
}

type Mutation struct {
}

type Query struct {
}

type Response struct {
	Message string `json:"message"`
}

type SignUpRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	IsVerified bool   `json:"is_verified"`
	Email      string `json:"email"`
}

type VerifyOtpRequest struct {
	Otp string `json:"otp"`
}
