package category

import (
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/database"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/pkg/image"
	"github.com/sika365/admin-tools/pkg/node"
)

type LocalCategories []*LocalCategory

type LocalCategory struct {
	models.CommonTableFields
	Title      string            `json:"title,omitempty"`
	Alias      string            `json:"alias,omitempty" gorm:"index"`
	Slug       string            `json:"slug,omitempty" gorm:"index"`
	Content    string            `json:"content,omitempty"`
	CoverID    database.NullPID  `json:"cover_id,omitempty"`
	CategoryID database.PID      `json:"category_id,omitempty"`
	Cover      *image.LocalImage `json:"cover,omitempty"`
	Nodes      node.Nodes        `json:"nodes,omitempty" gorm:"polymorphic:Owner"`
	Category   *models.Category  `json:"category,omitempty"`
}

func (LocalCategory) TableName() string {
	return "local_categories"
}

func (n *LocalCategory) Key() string {
	return n.ID.String()
}
