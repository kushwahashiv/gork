// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: queries.proto

package proto

import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto1 "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Queue) Validate() error {
	if this.Settings != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Settings); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Settings", err)
		}
	}
	return nil
}
func (this *Queue_Settings) Validate() error {
	if this.RateLimit != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.RateLimit); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("RateLimit", err)
		}
	}
	return nil
}
func (this *Queue_Settings_RateLimit) Validate() error {
	return nil
}
func (this *QueuesCmds) Validate() error {
	return nil
}
func (this *QueuesCmds_List) Validate() error {
	return nil
}
func (this *QueuesCmds_List_Request) Validate() error {
	if this.Params != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Params); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Params", err)
		}
	}
	return nil
}
func (this *QueuesCmds_List_Response) Validate() error {
	if this.Info != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Info); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Info", err)
		}
	}
	for _, item := range this.Records {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Records", err)
			}
		}
	}
	return nil
}
func (this *QueuesCmds_Create) Validate() error {
	return nil
}
func (this *QueuesCmds_Create_Request) Validate() error {
	if this.Settings != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Settings); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Settings", err)
		}
	}
	return nil
}
func (this *QueuesCmds_Create_Response) Validate() error {
	if this.Record != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Record); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Record", err)
		}
	}
	return nil
}
func (this *QueuesCmds_Read) Validate() error {
	return nil
}
func (this *QueuesCmds_Read_Request) Validate() error {
	return nil
}
func (this *QueuesCmds_Read_Response) Validate() error {
	if this.Record != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Record); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Record", err)
		}
	}
	return nil
}
func (this *QueuesCmds_Update) Validate() error {
	return nil
}
func (this *QueuesCmds_Update_Request) Validate() error {
	if this.Data != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Data); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Data", err)
		}
	}
	return nil
}
func (this *QueuesCmds_Update_Request_Data) Validate() error {
	if this.Settings != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Settings); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Settings", err)
		}
	}
	return nil
}
func (this *QueuesCmds_Update_Response) Validate() error {
	if this.Record != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Record); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Record", err)
		}
	}
	return nil
}
func (this *QueuesCmds_Delete) Validate() error {
	return nil
}
func (this *QueuesCmds_Delete_Request) Validate() error {
	return nil
}
func (this *QueuesCmds_Delete_Response) Validate() error {
	return nil
}
