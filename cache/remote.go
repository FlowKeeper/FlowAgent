package cache

import (
	"github.com/FlowKeeper/FlowUtils/v2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//This file is used to cache information we fetched from the remote server

var RemoteAgent models.Agent
var CurrentItems map[primitive.ObjectID]models.Item
