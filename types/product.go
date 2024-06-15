package types

// import (
// 	// buffers "ecommerce/buffers/productpb/protobuffs"

// 	"github.com/google/uuid"
// )

// type ProductID = ID

// type Product struct {
// 	Id          ProductID  `bson:"_id,omitempty"`
// 	Name        string     `bson:"name"`
// 	Price       float32    `bson:"price"`
// 	Description string     `bson:"description"`
// 	DealerId    DealerID   `bson:"dealerId"`
// 	CategoryId  CategoryID `bson:"categoryId"`
// 	Stock       int32      `bson:"stock"`
// }

// func (p *Product) SetID() {
// 	p.Id = uuid.New().String()
// }

// func ConvertProductRPCRequest(req *buffers.GetProductByIdRequest) *Product {
// 	return &Product{
// 		Id: req.ProductId,
// 	}
// }

// func ConvertProductRPCIncRequest(req *buffers.UpdateStockByIdRequest) *Product {
// 	return &Product{
// 		Id: req.ProductId,
// 	}
// }
