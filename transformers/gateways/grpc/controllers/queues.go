package controllers

import (
	"time"

	"github.com/gork-io/gork/models"
	"github.com/gork-io/gork/services/resources"
	"github.com/gork-io/gork/transformers/gateways/grpc/proto"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewQueues creates a new instance of Queues.
func NewQueues(queuesSvc *resources.Queues) (ctrl *Queues) {
	return &Queues{
		queuesSvc: queuesSvc,
	}
}

// Queues controller is a proxy that links GRPC gateway with service layer.
type Queues struct {
	queuesSvc *resources.Queues // queues service
}

// Register registers this controller as a GRPC service implementation.
func (ctrl *Queues) Register(server *grpc.Server) {
	proto.RegisterQueuesServer(server, ctrl)
}

// List returns a subset of the queries, based on collection params given.
func (ctrl *Queues) List(ctx context.Context, request *proto.QueuesCmds_List_Request) (response *proto.QueuesCmds_List_Response, err error) {

	// Fetch records
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

	// Convert settings
	settings := make(map[models.QueueSetting]string)
	for _, setting := range request.Settings {
		settings[models.QueueSetting(setting.Key)] = setting.Value
	}

	// Create record
	record, err := ctrl.queuesSvc.Create(ctx, request.Name, settings)
	if err != nil {
		return nil, errors.Wrap(err, "create failed")
	}

	// Return response
	response = &proto.QueuesCmds_Create_Response{
		Record: marshalQueue(record),
	}

	return
}

// Read returns query by its id.
func (ctrl *Queues) Read(ctx context.Context, request *proto.QueuesCmds_Read_Request) (response *proto.QueuesCmds_Read_Response, err error) {

	// Fetch record
	record, err := ctrl.queuesSvc.Read(ctx, request.Id)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	// Return response
	response = &proto.QueuesCmds_Read_Response{
		Record: marshalQueue(record),
	}

	return
}

// Delete removes queue with given ID.
func (ctrl *Queues) Delete(ctx context.Context, request *proto.QueuesCmds_Delete_Request) (response *proto.QueuesCmds_Delete_Response, err error) {

	response = &proto.QueuesCmds_Delete_Response{}

	// Delete record
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

	output = &proto.Queue{
		Id:        input.Id,
		Name:      input.Name,
		CreatedAt: input.CreatedAt.Format(time.RFC3339Nano),
	}
	for key, value := range input.Settings {
		output.Settings = append(output.Settings, &proto.Queue_Setting{
			Key:   string(key),
			Value: value,
		})
	}

	return
}
