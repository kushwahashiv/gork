package controllers

import (
	"context"

	"time"

	"github.com/gork-io/gork/models"
	"github.com/gork-io/gork/services"
	"github.com/gork-io/gork/transformers/gateways/grpc/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Queues controller
type Queues struct {
	queuesSvc services.Queues // queues service
}

func (ctrl *Queues) Register(server *grpc.Server) {
	proto.RegisterQueuesServer(server, ctrl)
}

// List returns a subset of the queries, based on collection params given.
func (ctrl *Queues) List(ctx context.Context, request *proto.QueuesCmds_List_Request) (response *proto.QueuesCmds_List_Response, err error) {

	// Fetch queries
	records, info, err := ctrl.queuesSvc.List(ctx, models.NewCollectionParams(
		request.Params.Cursor,
		request.Params.Limit,
	))
	if err != nil {
		return nil, errors.Wrap(err, "list failed")
	}

	// Return response
	response = &proto.QueuesCmds_List_Response{
		Info: marshalCollectionInfo(info),
	}
	for _, record := range records {
		response.Records = append(response.Records, marshalQueue(record))
	}

	return
}

// Create creates a new queue.
func (ctrl *Queues) Create(ctx context.Context, request *proto.QueuesCmds_Create_Request) (response *proto.QueuesCmds_Create_Response, err error) {

	var settings []models.QueueSetting
	if request.Settings.RateLimit != nil {
		settings = append(settings, models.QueueWithRateLimit(
			request.Settings.RateLimit.Tokens,
			time.Duration(time.Duration(request.Settings.RateLimit.Duration)*time.Second),
		))
	}

	record, err := ctrl.queuesSvc.Create(ctx, request.Name, settings...)
	if err != nil {
		return nil, errors.Wrap(err, "create failed")
	}

	response = &proto.QueuesCmds_Create_Response{
		Record: marshalQueue(record),
	}

	return
}

// Read returns query by its id.
func (ctrl *Queues) Read(ctx context.Context, request *proto.QueuesCmds_Read_Request) (response *proto.QueuesCmds_Read_Response, err error) {

	record, err := ctrl.queuesSvc.Read(ctx, request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	response = &proto.QueuesCmds_Read_Response{
		Record: marshalQueue(record),
	}

	return
}

func (ctrl *Queues) Update(ctx context.Context, request *proto.QueuesCmds_Update_Request) (response *proto.QueuesCmds_Update_Response, err error) {
	return
}

// Delete removes queue with given ID.
func (ctrl *Queues) Delete(ctx context.Context, request *proto.QueuesCmds_Delete_Request) (response *proto.QueuesCmds_Delete_Response, err error) {

	response = &proto.QueuesCmds_Delete_Response{}

	err = ctrl.queuesSvc.Delete(ctx, request.Id)
	if err == nil {
		response.Result = true
	}

	return response, errors.Wrap(err, "delete failed")
}

// marshalQueue is a helper function that marshals domain model of the queue into GRCP model.
func marshalQueue(input *models.Queue) (output *proto.Queue) {

	if input == nil {
		return nil
	}

	return &proto.Queue{
		Id:   input.Id,
		Name: input.Name,
		Settings: &proto.Queue_Settings{
			RateLimit: &proto.Queue_Settings_RateLimit{
				Tokens:   input.Settings.RateLimit.Tokens,
				Duration: uint32(input.Settings.RateLimit.Duration.Seconds()),
			},
		},
		CreatedAt: input.CreatedAt.Format(time.RFC3339Nano),
	}
}
