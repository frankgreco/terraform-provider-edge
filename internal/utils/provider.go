package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mattbaird/jsonpatch"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type api interface {
	Read(context.Context, string) (interface{}, error)
	Create(context.Context, interface{}) (interface{}, error)
	Update(context.Context, interface{}, []jsonpatch.JsonPatchOperation) (interface{}, error)
	Delete(context.Context, string) error
}

type Resource struct {
	Type         interface{}
	Name         string
	IsConfigured bool
	Attribute    string
	Api          api
}

func (r Resource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	CreateFunc(
		ctx,
		req,
		resp,
		r.Type,
		r.Name,
		r.IsConfigured,
		r.Api.Create,
	)
}

func (r Resource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	ReadFunc(
		ctx,
		req,
		resp,
		r.Attribute,
		r.Name,
		r.IsConfigured,
		r.Api.Read,
	)
}

func (r Resource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	UpdateFunc(
		ctx,
		req,
		resp,
		r.Name,
		r.Type,
		r.Api.Update,
	)
}

func (r Resource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	DeleteFunc(
		ctx,
		req,
		resp,
		r.Attribute,
		r.Name,
		r.IsConfigured,
		r.Api.Delete,
	)
}

func (r Resource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	ImportFunc(
		ctx,
		req,
		resp,
		r.Name,
		r.IsConfigured,
		r.Api.Read,
	)
}

func ImportFunc(
	ctx context.Context,
	req tfsdk.ImportResourceStateRequest,
	resp *tfsdk.ImportResourceStateResponse,
	resource string,
	configured bool,
	f func(context.Context, string) (interface{}, error),
) {
	log.Printf("[TRACE] begin import func for %s", resource)
	defer func() {
		log.Printf("[TRACE] end import func for %s", resource)
	}()

	if !configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	actual, err := f(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue retrieving the %s %s.", resource, req.ID),
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, actual)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func DeleteFunc(
	ctx context.Context,
	req tfsdk.DeleteResourceRequest,
	resp *tfsdk.DeleteResourceResponse,
	idAttribute string,
	resource string,
	configured bool,
	f func(context.Context, string) error,
) {
	log.Printf("[TRACE] begin delete func for %s", resource)
	defer func() {
		log.Printf("[TRACE] end delete func for %s", resource)
	}()

	if !configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var id string
	{
		diagnostics := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName(idAttribute), &id)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if err := f(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue deleting the %s.", resource),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func CreateFunc(
	ctx context.Context,
	req tfsdk.CreateResourceRequest,
	resp *tfsdk.CreateResourceResponse,
	resourceType interface{},
	resource string,
	configured bool,
	f func(context.Context, interface{}) (interface{}, error),
) {
	log.Printf("[TRACE] begin create func for %s", resource)
	defer func() {
		log.Printf("[TRACE] end create func for %s", resource)
	}()

	if !configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	retrieved, diags := retrieve(ctx, req.Plan, resourceType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	actual, err := f(ctx, retrieved)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue creating the %s.", resource),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, actual)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ReadFunc(
	ctx context.Context,
	req tfsdk.ReadResourceRequest,
	resp *tfsdk.ReadResourceResponse,
	idAttribute string,
	resource string,
	configured bool,
	f func(context.Context, string) (interface{}, error),
) {
	log.Printf("[TRACE] begin read func for %s", resource)
	defer func() {
		log.Printf("[TRACE] end read func for %s", resource)
	}()

	if !configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var id string
	{
		diagnostics := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName(idAttribute), &id)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	actual, err := f(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue retrieving the %s %s.", resource, id),
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, actual)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func UpdateFunc(
	ctx context.Context,
	req tfsdk.UpdateResourceRequest,
	resp *tfsdk.UpdateResourceResponse,
	name string,
	resourceType interface{},
	f func(context.Context, interface{}, []jsonpatch.JsonPatchOperation) (interface{}, error),
) {
	log.Printf("[TRACE] begin update func for %s", name)
	defer func() {
		log.Printf("[TRACE] end update func for %s", name)
	}()

	current, diags := retrieve(ctx, req.State, resourceType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("[TRACE] current %s struct: %+v", name, current)

	desired, diags := retrieve(ctx, req.Plan, resourceType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	log.Printf("[TRACE] desired %s struct: %+v", name, desired)

	var patches []jsonpatch.JsonPatchOperation
	{
		cData, err := json.Marshal(&current)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not marshal %s from state", name),
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] current %s json: %s", name, string(cData))

		dData, err := json.Marshal(&desired)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not marshal %s from plan", name),
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] desired %s json: %s", name, string(dData))

		p, err := jsonpatch.CreatePatch(cData, dData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not create patch document.",
				err.Error(),
			)
			return
		}
		patches = p
		log.Printf("[DEBUG] patch document: %+v", patches)
	}

	updated, err := f(ctx, current, patches)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue updating the.", name),
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, updated)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}
