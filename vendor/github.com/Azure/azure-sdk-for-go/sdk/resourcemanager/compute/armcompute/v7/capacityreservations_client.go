// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// CapacityReservationsClient contains the methods for the CapacityReservations group.
// Don't use this type directly, use NewCapacityReservationsClient() instead.
type CapacityReservationsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCapacityReservationsClient creates a new instance of CapacityReservationsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCapacityReservationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CapacityReservationsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CapacityReservationsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - The operation to create or update a capacity reservation. Please note some properties can be set
// only during capacity reservation creation. Please refer to https://aka.ms/CapacityReservation for more
// details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - capacityReservationGroupName - The name of the capacity reservation group.
//   - capacityReservationName - The name of the capacity reservation.
//   - parameters - Parameters supplied to the Create capacity reservation.
//   - options - CapacityReservationsClientBeginCreateOrUpdateOptions contains the optional parameters for the CapacityReservationsClient.BeginCreateOrUpdate
//     method.
func (client *CapacityReservationsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters CapacityReservation, options *CapacityReservationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[CapacityReservationsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CapacityReservationsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CapacityReservationsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - The operation to create or update a capacity reservation. Please note some properties can be set only
// during capacity reservation creation. Please refer to https://aka.ms/CapacityReservation for more
// details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
func (client *CapacityReservationsClient) createOrUpdate(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters CapacityReservation, options *CapacityReservationsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "CapacityReservationsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *CapacityReservationsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters CapacityReservation, _ *CapacityReservationsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}/capacityReservations/{capacityReservationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if capacityReservationGroupName == "" {
		return nil, errors.New("parameter capacityReservationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationGroupName}", url.PathEscape(capacityReservationGroupName))
	if capacityReservationName == "" {
		return nil, errors.New("parameter capacityReservationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationName}", url.PathEscape(capacityReservationName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - The operation to delete a capacity reservation. This operation is allowed only when all the associated resources
// are disassociated from the capacity reservation. Please refer to
// https://aka.ms/CapacityReservation for more details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - capacityReservationGroupName - The name of the capacity reservation group.
//   - capacityReservationName - The name of the capacity reservation.
//   - options - CapacityReservationsClientBeginDeleteOptions contains the optional parameters for the CapacityReservationsClient.BeginDelete
//     method.
func (client *CapacityReservationsClient) BeginDelete(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, options *CapacityReservationsClientBeginDeleteOptions) (*runtime.Poller[CapacityReservationsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CapacityReservationsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CapacityReservationsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - The operation to delete a capacity reservation. This operation is allowed only when all the associated resources
// are disassociated from the capacity reservation. Please refer to
// https://aka.ms/CapacityReservation for more details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
func (client *CapacityReservationsClient) deleteOperation(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, options *CapacityReservationsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "CapacityReservationsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *CapacityReservationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, _ *CapacityReservationsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}/capacityReservations/{capacityReservationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if capacityReservationGroupName == "" {
		return nil, errors.New("parameter capacityReservationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationGroupName}", url.PathEscape(capacityReservationGroupName))
	if capacityReservationName == "" {
		return nil, errors.New("parameter capacityReservationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationName}", url.PathEscape(capacityReservationName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - The operation that retrieves information about the capacity reservation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - capacityReservationGroupName - The name of the capacity reservation group.
//   - capacityReservationName - The name of the capacity reservation.
//   - options - CapacityReservationsClientGetOptions contains the optional parameters for the CapacityReservationsClient.Get
//     method.
func (client *CapacityReservationsClient) Get(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, options *CapacityReservationsClientGetOptions) (CapacityReservationsClientGetResponse, error) {
	var err error
	const operationName = "CapacityReservationsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, options)
	if err != nil {
		return CapacityReservationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CapacityReservationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CapacityReservationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *CapacityReservationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, options *CapacityReservationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}/capacityReservations/{capacityReservationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if capacityReservationGroupName == "" {
		return nil, errors.New("parameter capacityReservationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationGroupName}", url.PathEscape(capacityReservationGroupName))
	if capacityReservationName == "" {
		return nil, errors.New("parameter capacityReservationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationName}", url.PathEscape(capacityReservationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", string(*options.Expand))
	}
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CapacityReservationsClient) getHandleResponse(resp *http.Response) (CapacityReservationsClientGetResponse, error) {
	result := CapacityReservationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CapacityReservation); err != nil {
		return CapacityReservationsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByCapacityReservationGroupPager - Lists all of the capacity reservations in the specified capacity reservation group.
// Use the nextLink property in the response to get the next page of capacity reservations.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - capacityReservationGroupName - The name of the capacity reservation group.
//   - options - CapacityReservationsClientListByCapacityReservationGroupOptions contains the optional parameters for the CapacityReservationsClient.NewListByCapacityReservationGroupPager
//     method.
func (client *CapacityReservationsClient) NewListByCapacityReservationGroupPager(resourceGroupName string, capacityReservationGroupName string, options *CapacityReservationsClientListByCapacityReservationGroupOptions) *runtime.Pager[CapacityReservationsClientListByCapacityReservationGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[CapacityReservationsClientListByCapacityReservationGroupResponse]{
		More: func(page CapacityReservationsClientListByCapacityReservationGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CapacityReservationsClientListByCapacityReservationGroupResponse) (CapacityReservationsClientListByCapacityReservationGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "CapacityReservationsClient.NewListByCapacityReservationGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByCapacityReservationGroupCreateRequest(ctx, resourceGroupName, capacityReservationGroupName, options)
			}, nil)
			if err != nil {
				return CapacityReservationsClientListByCapacityReservationGroupResponse{}, err
			}
			return client.listByCapacityReservationGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByCapacityReservationGroupCreateRequest creates the ListByCapacityReservationGroup request.
func (client *CapacityReservationsClient) listByCapacityReservationGroupCreateRequest(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, _ *CapacityReservationsClientListByCapacityReservationGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}/capacityReservations"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if capacityReservationGroupName == "" {
		return nil, errors.New("parameter capacityReservationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationGroupName}", url.PathEscape(capacityReservationGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByCapacityReservationGroupHandleResponse handles the ListByCapacityReservationGroup response.
func (client *CapacityReservationsClient) listByCapacityReservationGroupHandleResponse(resp *http.Response) (CapacityReservationsClientListByCapacityReservationGroupResponse, error) {
	result := CapacityReservationsClientListByCapacityReservationGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CapacityReservationListResult); err != nil {
		return CapacityReservationsClientListByCapacityReservationGroupResponse{}, err
	}
	return result, nil
}

// BeginUpdate - The operation to update a capacity reservation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - capacityReservationGroupName - The name of the capacity reservation group.
//   - capacityReservationName - The name of the capacity reservation.
//   - parameters - Parameters supplied to the Update capacity reservation operation.
//   - options - CapacityReservationsClientBeginUpdateOptions contains the optional parameters for the CapacityReservationsClient.BeginUpdate
//     method.
func (client *CapacityReservationsClient) BeginUpdate(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters CapacityReservationUpdate, options *CapacityReservationsClientBeginUpdateOptions) (*runtime.Poller[CapacityReservationsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CapacityReservationsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CapacityReservationsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - The operation to update a capacity reservation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-11-01
func (client *CapacityReservationsClient) update(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters CapacityReservationUpdate, options *CapacityReservationsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "CapacityReservationsClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, capacityReservationGroupName, capacityReservationName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *CapacityReservationsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string, parameters CapacityReservationUpdate, _ *CapacityReservationsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}/capacityReservations/{capacityReservationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if capacityReservationGroupName == "" {
		return nil, errors.New("parameter capacityReservationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationGroupName}", url.PathEscape(capacityReservationGroupName))
	if capacityReservationName == "" {
		return nil, errors.New("parameter capacityReservationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{capacityReservationName}", url.PathEscape(capacityReservationName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}
