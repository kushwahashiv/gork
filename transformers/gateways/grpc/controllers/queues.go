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

func (ctrl *Queues) List(ctx context.Context, request *proto.QueuesCmds_List_Request) (response *proto.QueuesCmds_List_Response, err error) {

	records, info, err := ctrl.queuesSvc.List(ctx, models.NewCollectionParams(
		request.Params.GetCursor(),
		uint8(request.Params.GetLimit()),
	))
	if err != nil {
		return nil, err
	}

	response = &proto.QueuesCmds_List_Response{
		Info: &proto.Collection_Info{
			Cursor: info.Cursor,
			Total:  info.Total,
		},
	}
	for _, record := range records {
		response.Records = append(response.Records)
	}

	return
}

// Create creates a new queue.
func (ctrl *Queues) Create(ctx context.Context, request *proto.QueuesCmds_Create_Request) (response *proto.QueuesCmds_Create_Response, err error) {

	var settings []models.QueueSetting

	if request.Settings.GetRateLimit() != nil {
		settings = append(settings, models.QueueWithRateLimit(
			request.Settings.RateLimit.Tokens,
			time.Duration(time.Duration(request.Settings.RateLimit.Duration)*time.Second),
		))
	}

	record, err := ctrl.queuesSvc.Create(ctx, request.Name, settings...)
	if err != nil {
		err = errors.Wrap(err, "create failed")
		return nil, err
	}

	response.Record.Id = record.Id
	response.Record.Name = record.Name

	return
}

func (ctrl *Queues) Read(ctx context.Context, request *proto.QueuesCmds_Read_Request) (response *proto.QueuesCmds_Read_Response, err error) {
	return
}

func (ctrl *Queues) Update(ctx context.Context, request *proto.QueuesCmds_Update_Request) (response *proto.QueuesCmds_Update_Response, err error) {
	return
}

func (ctrl *Queues) Delete(ctx context.Context, request *proto.QueuesCmds_Delete_Request) (response *proto.QueuesCmds_Delete_Response, err error) {
	return
}
