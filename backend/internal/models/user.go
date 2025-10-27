package models

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleAccountant UserRole = "accountant"
	RoleViewer     UserRole = "viewer"
)

type User struct {
	BaseModel
	Email     string   `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string   `gorm:"not null;size:255" json:"-"`
	FullName  string   `gorm:"not null;size:255" json:"full_name"`
	Role      UserRole `gorm:"type:varchar(20);not null;default:'viewer'" json:"role"`
	CompanyID uint     `gorm:"index;default:null" json:"company_id"`
	Company   Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	IsActive  bool     `gorm:"default:true" json:"is_active"`
}

type UserResponse struct {
	ID        uint     `json:"id"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	Role      UserRole `json:"role"`
	CompanyID uint     `json:"company_id"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
}