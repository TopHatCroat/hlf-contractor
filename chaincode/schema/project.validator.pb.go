// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: project.proto

package schema

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Project) Validate() error {
	for _, item := range this.ApplicationsIds {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("ApplicationsIds", err)
			}
		}
	}
	if this.OpenDate != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.OpenDate); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("OpenDate", err)
		}
	}
	if this.StartDate != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.StartDate); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("StartDate", err)
		}
	}
	if this.EndDate != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.EndDate); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("EndDate", err)
		}
	}
	for _, item := range this.Applications {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Applications", err)
			}
		}
	}
	return nil
}
func (this *ProjectId) Validate() error {
	return nil
}
func (this *ProjectList) Validate() error {
	for _, item := range this.Items {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Items", err)
			}
		}
	}
	return nil
}
func (this *PublishProject) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if nil == this.StartDate {
		return github_com_mwitkow_go_proto_validators.FieldError("StartDate", fmt.Errorf("message must exist"))
	}
	if this.StartDate != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.StartDate); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("StartDate", err)
		}
	}
	if nil == this.EndDate {
		return github_com_mwitkow_go_proto_validators.FieldError("EndDate", fmt.Errorf("message must exist"))
	}
	if this.EndDate != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.EndDate); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("EndDate", err)
		}
	}
	if !(this.EstimatedValue > 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("EstimatedValue", fmt.Errorf(`value '%v' must be greater than '0'`, this.EstimatedValue))
	}
	if this.Description == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Description", fmt.Errorf(`value '%v' must not be an empty string`, this.Description))
	}
	return nil
}
