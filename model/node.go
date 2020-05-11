package model


import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A node is an IoT SBC. A node can have actuators and sensors attached to it. A node implements
// the state management and protocol of iot/protocol 
type Node struct {
	Id primative.ObjectId `bson:"_id,omitempty"`
	SerialNumber string `bson:"serialno"`
	Firmware string `bson:"firmware"`
	
}
