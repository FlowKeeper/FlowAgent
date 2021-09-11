package cache

import (
	"github.com/FlowKeeper/FlowUtils/v2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//This file is used to cache information we fetched from the remote server

//RemoteAgent stores the struct returned from the server when calling /api/v1/config
var RemoteAgent models.Agent

//CurrentItems maps item ids to the actual item structs
//This is used by the scheduler to load the current item struct into the go thread
var CurrentItems map[primitive.ObjectID]models.Item
