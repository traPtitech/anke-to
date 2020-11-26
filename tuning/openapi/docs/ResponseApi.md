# \ResponseApi

All URIs are relative to *https://anke-to.trap.jp/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteResponse**](ResponseApi.md#DeleteResponse) | **Delete** /responses/{responseID} | 
[**GetResponses**](ResponseApi.md#GetResponses) | **Get** /responses/{responseID} | 
[**PatchResponse**](ResponseApi.md#PatchResponse) | **Patch** /responses/{responseID} | 
[**PostResponse**](ResponseApi.md#PostResponse) | **Post** /responses | 



## DeleteResponse

> DeleteResponse(ctx, responseID)



回答を削除します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**responseID** | **int32**| 回答ID  | 

### Return type

 (empty response body)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetResponses

> Response GetResponses(ctx, responseID)



あるresponseIDを持つ回答に含まれる全ての質問に対する自分の回答を取得します

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**responseID** | **int32**| 回答ID  | 

### Return type

[**Response**](Response.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchResponse

> PatchResponse(ctx, responseID, newResponse)



回答を変更します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**responseID** | **int32**| 回答ID  | 
**newResponse** | [**NewResponse**](NewResponse.md)|  | 

### Return type

 (empty response body)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostResponse

> ResponseDetails PostResponse(ctx, newResponse)



新しい回答を作成します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**newResponse** | [**NewResponse**](NewResponse.md)|  | 

### Return type

[**ResponseDetails**](ResponseDetails.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

