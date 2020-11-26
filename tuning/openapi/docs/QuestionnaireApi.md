# \QuestionnaireApi

All URIs are relative to *https://anke-to.trap.jp/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DelteQuestionnaire**](QuestionnaireApi.md#DelteQuestionnaire) | **Delete** /questionnaires/{questionnaireID} | 
[**GetQuestionnaire**](QuestionnaireApi.md#GetQuestionnaire) | **Get** /questionnaires/{questionnaireID} | 
[**GetQuestionnaires**](QuestionnaireApi.md#GetQuestionnaires) | **Get** /questionnaires | 
[**GetQuestions**](QuestionnaireApi.md#GetQuestions) | **Get** /questionnaires/{questionnaireID}/questions | 
[**PatchQuestionnaire**](QuestionnaireApi.md#PatchQuestionnaire) | **Patch** /questionnaires/{questionnaireID} | 
[**PostQuestionnaire**](QuestionnaireApi.md#PostQuestionnaire) | **Post** /questionnaires | 



## DelteQuestionnaire

> DelteQuestionnaire(ctx, questionnaireID)



アンケートを削除します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionnaireID** | **int32**| アンケートID  | 

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


## GetQuestionnaire

> QuestionnaireById GetQuestionnaire(ctx, questionnaireID)



アンケートの情報を取得します。

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionnaireID** | **int32**| アンケートID  | 

### Return type

[**QuestionnaireById**](QuestionnaireByID.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetQuestionnaires

> []QuestionnaireForList GetQuestionnaires(ctx, sort, page, nontargeted)



与えられた条件を満たす20件以下のアンケートのリストを取得します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**sort** | **string**| 並び順 (作成日時が新しい \&quot;created_at\&quot;, 作成日時が古い \&quot;-created_at\&quot;, タイトルの昇順 \&quot;title\&quot;, タイトルの降順 \&quot;-title\&quot;, 更新日時が新しい \&quot;modified_at\&quot;, 更新日時が古い \&quot;-modified_at\&quot; )  | 
**page** | **int32**| 何ページ目か (未定義の場合は1ページ目) | 
**nontargeted** | **bool**| 自分がターゲットになっていないもののみ取得 (true), ターゲットになっているものも含めてすべて取得 (false)  | 

### Return type

[**[]QuestionnaireForList**](QuestionnaireForList.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetQuestions

> []QuestionDetails GetQuestions(ctx, questionnaireID)



アンケートに含まれる質問のリストを取得します。

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionnaireID** | **int32**| アンケートID  | 

### Return type

[**[]QuestionDetails**](QuestionDetails.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchQuestionnaire

> PatchQuestionnaire(ctx, questionnaireID, newQuestionnaire)



アンケートの情報を変更します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**questionnaireID** | **int32**| アンケートID  | 
**newQuestionnaire** | [**NewQuestionnaire**](NewQuestionnaire.md)|  | 

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


## PostQuestionnaire

> NewQuestionnaireResponse PostQuestionnaire(ctx, newQuestionnaire)



新しいアンケートを作成します．

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**newQuestionnaire** | [**NewQuestionnaire**](NewQuestionnaire.md)|  | 

### Return type

[**NewQuestionnaireResponse**](NewQuestionnaireResponse.md)

### Authorization

[application](../README.md#application)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

