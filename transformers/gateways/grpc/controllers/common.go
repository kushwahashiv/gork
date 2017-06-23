package controllers

import (
	"github.com/gork-io/gork/models"
	"github.com/gork-io/gork/transformers/gateways/grpc/proto"
)

// marshalCollectionInfo is a helper function that marshals domain model of the collection info into GRCP model.
func marshalCollectionInfo(input *models.CollectionInfo) (output *proto.Collection_Info) {

	if input == nil {
		return nil
	}

	return &proto.Collection_Info{
		Cursor: input.Cursor,
		Total:  input.Total,
	}
}
