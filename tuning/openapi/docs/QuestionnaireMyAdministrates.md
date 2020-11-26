# QuestionnaireMyAdministrates

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**QuestionnaireID** | **int32** |  | 
**Title** | **string** |  | 
**Description** | **string** |  | 
**ResTimeLimit** | [**time.Time**](time.Time.md) |  | 
**CreatedAt** | [**time.Time**](time.Time.md) |  | 
**ModifiedAt** | [**time.Time**](time.Time.md) |  | 
**ResSharedTo** | **string** | アンケートの結果を, 運営は見られる (\&quot;administrators\&quot;), 回答済みの人は見られる (\&quot;respondents\&quot;) 誰でも見られる (\&quot;public\&quot;)  | 
**Targets** | **[]string** |  | 
**Administrators** | **[]string** |  | 
**AllResponded** | **bool** | 回答必須でない場合、またはすべてのターゲットが回答済みの場合、true を返す。それ以外はfalseを返す。  | 
**Respondents** | **[]string** |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


