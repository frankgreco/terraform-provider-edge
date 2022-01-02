package provider

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/frankgreco/edge-sdk-go/types"

// 	"github.com/hashicorp/terraform-plugin-framework/diag"
// 	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
// 	tfftypes "github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-go/tftypes"
// )

// type resourcePortGroupType struct{}

// func (r resourcePortGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
// 	return tfsdk.Schema{
// 		Description: "",
// 		Attributes: map[string]tfsdk.Attribute{
// 			"name": {
// 				Type:          tfftypes.StringType,
// 				Required:      true,
// 				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
// 				Description:   "",
// 			},
// 			"description": {
// 				Type:        tfftypes.StringType,
// 				Optional:    true,
// 				Description: "",
// 			},
// 			"port_ranges": {
// 				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
// 					"from": {
// 						Type:     types.StringType,
// 						Required: true,
// 					},
// 					"to": {
// 						Type:     types.StringType,
// 						Required: true,
// 					},
// 				}, tfsdk.ListNestedAttributesOptions{}),
// 				Optional:    true,
// 				Description: "",
// 			},
// 			"ports": {
// 				Type: tfftypes.ListType{
// 					ElemType: tfftypes.StringType,
// 				},
// 				Optional:    true,
// 				Description: "",
// 			},
// 		},
// 	}, nil
// }

// // New resource instance
// func (r resourcePortGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
// 	return resourcePortGroup{
// 		p: *(p.(*provider)),
// 	}, nil
// }

// type resourcePortGroup struct {
// 	p provider
// }

// func (r resourcePortGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

// }

// func (r resourcePortGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

// }

// func (r resourcePortGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

// }

// func (r resourcePortGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

// }
