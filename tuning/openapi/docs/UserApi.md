# \UserApi

All URIs are relative to *https://anke-to.trap.jp/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetMe**](UserApi.md#GetMe) | **Get** /users/me | 
[**GetMyAdministrates**](UserApi.md#GetMyAdministrates) | **Get** /users/me/administrates | 
[**GetMyResponses**](UserApi.md#GetMyResponses) | **Get** /users/me/responses | 
[**GetMyTargeted**](UserApi.md#GetMyTargeted) | **Get** /users/me/targeted | 
[**GetResponsesToQuestionnaire**](UserApi.md#GetResponsesToQuestionnaire) | **Get** /users/me/responses/{questionnaireID} | 
[**GetUsers**](UserApi.md#GetUsers) | **Get** /users | 未実装



## GetMe

> Me GetMe(ctx, )



自分のユーザー情報を取得します

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**Me**](Me.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMyAdministrates

> []QuestionnaireMyAdministrates GetMyAdministrates(ctx, )



自分が管理者になっているアンケートのリストを取得します。

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]QuestionnaireMyAdministrates**](QuestionnaireMyAdministrates.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMyResponses

> []ResponseSummary GetMyResponses(ctx, )



自分のすべての回答のリストを取得します。

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]ResponseSummary**](ResponseSummary.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMyTargeted

> []QuestionnaireMyTargeted GetMyTargeted(ctx, )



自分が対象になっている アンケートのリストを取得します。

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]QuestionnaireMyTargeted**](QuestionnaireMyTargeted.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetResponsesToQuestionnaire

> []ResponseSummary GetResponsesToQuestionnaire(ctx, questionnaireID)



特定のquestionnaireIdを持つアンケートに対する自分のすべての回答のリストを取得します。

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionnaireID** | **int32**| アンケートID  | 

### Return type

[**[]ResponseSummary**](ResponseSummary.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUsers

> []User GetUsers(ctx, )

未実装

(botおよび除名されたユーザーを除く、全ての) ユーザーのtraQIDのリストを取得します。

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]User**](User.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

