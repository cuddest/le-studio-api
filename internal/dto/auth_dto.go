package dto

// RegisterDTO defines user registration payload.
type RegisterDTO struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender" validate:"omitempty,oneof=male female"`
}

// RegisterAdminDTO defines admin registration payload.
type RegisterAdminDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=80"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginDTO defines login payload.
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RefreshDTO defines refresh token payload.
type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// GuestPurchaseDTO defines guest purchase payload.
type GuestPurchaseDTO struct {
	FirstName      string `json:"first_name" validate:"required,min=2,max=50"`
	LastName       string `json:"last_name" validate:"required,min=2,max=50"`
	Email          string `json:"email" validate:"required,email"`
	Phone          string `json:"phone"`
	Gender         string `json:"gender" validate:"omitempty,oneof=male female"`
	PackTemplateID uint   `json:"pack_template_id" validate:"required"`
}
