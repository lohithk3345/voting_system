package types

// import (
// 	"time"

// 	"github.com/google/uuid"
// )

// type CartID = ID

// type Cart struct {
// 	Id        CartID  `bson:"_id,omitempty" json:"cartId"`
// 	Product   Product `bson:"product,omitempty" json:"product"`
// 	UserId    UserID  `bson:"userId" json:"userId"`
// 	UpdatedAt time.Time
// 	CreatedAt time.Time
// }

// func (c *Cart) SetID() {
// 	c.Id = uuid.New().String()
// }

// func (c *Cart) SetUpdatedTimestamp() {
// 	c.UpdatedAt = time.Now().UTC()
// }

// func (c *Cart) SetCreatedTimestamp() {
// 	c.CreatedAt = time.Now().UTC()
// }
