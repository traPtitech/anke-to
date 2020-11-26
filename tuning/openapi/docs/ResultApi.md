# \ResultApi

All URIs are relative to *https://anke-to.trap.jp/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetResults**](ResultApi.md#GetResults) | **Get** /results/{questionnaireID} | 



## GetResults

> []ResponseResult GetResults(ctx, questionnaireID)



あるquestionnaireIDを持つアンケートの結果をすべて取得します。

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionnaireID** | **int32**| アンケートID  | 

### Return type

[**[]ResponseResult**](ResponseResult.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

