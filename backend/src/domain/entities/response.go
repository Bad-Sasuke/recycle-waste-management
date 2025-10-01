package entities

type ResponseMessage struct {
	Message string `json:"message" bson:"message,omitempty"`
}

type ResponseModel struct {
	Message string      `json:"message" bson:"message,omitempty"`
	Data    interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Status  int         `json:"status,omitempty" bson:"status,omitempty"`
}

type ResponsePaginationModel struct {
	Message    string      `json:"message" bson:"message,omitempty"`
	Data       interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Page       int         `json:"page" bson:"page,omitempty"`
	Limit      int         `json:"limit" bson:"limit,omitempty"`
	TotalPages int         `json:"total_pages" bson:"total_pages,omitempty"`
	TotalItems int64       `json:"total_items" bson:"total_items,omitempty"`
	Status     int         `json:"status,omitempty" bson:"status,omitempty"`
}

type ResponseBool struct {
	Message string `json:"message" bson:"message,omitempty"`
	IsTrue  bool   `json:"istrue,omitempty" bson:"istrue,omitempty"`
}
