# \QuestionApi

All URIs are relative to *https://anke-to.trap.jp/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteQuestion**](QuestionApi.md#DeleteQuestion) | **Delete** /questions/{questionID} | 
[**PatchQuestion**](QuestionApi.md#PatchQuestion) | **Patch** /questions/{questionID} | 
[**PostQuestion**](QuestionApi.md#PostQuestion) | **Post** /questions | 



## DeleteQuestion

> DeleteQuestion(ctx, questionID)



質問を削除します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionID** | **int32**| 質問ID  | 

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


## PatchQuestion

> PatchQuestion(ctx, questionID, newQuestion)



質問を変更します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionID** | **int32**| 質問ID  | 
**newQuestion** | [**NewQuestion**](NewQuestion.md)|  | 

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


## PostQuestion

> Question PostQuestion(ctx, newQuestion)



新しい質問を作成します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**newQuestion** | [**NewQuestion**](NewQuestion.md)|  | 

### Return type

[**Question**](Question.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

