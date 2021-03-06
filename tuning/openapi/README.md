# Go API client for openapi

anke-to API

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 1.0.0-oas3
- Package version: 1.0.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [https://github.com/traPtitech/anke-to](https://github.com/traPtitech/anke-to)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/oauth2
go get golang.org/x/net/context
go get github.com/antihax/optional
```

Put the package under your project folder and add the following in import:

```golang
import "./openapi"
```

## Documentation for API Endpoints

All URIs are relative to *https://anke-to.trap.jp/api*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*GroupApi* | [**GetGroups**](docs/GroupApi.md#getgroups) | **Get** /groups | 未実装
*QuestionApi* | [**DeleteQuestion**](docs/QuestionApi.md#deletequestion) | **Delete** /questions/{questionID} | 
*QuestionApi* | [**PatchQuestion**](docs/QuestionApi.md#patchquestion) | **Patch** /questions/{questionID} | 
*QuestionApi* | [**PostQuestion**](docs/QuestionApi.md#postquestion) | **Post** /questions | 
*QuestionnaireApi* | [**DelteQuestionnaire**](docs/QuestionnaireApi.md#deltequestionnaire) | **Delete** /questionnaires/{questionnaireID} | 
*QuestionnaireApi* | [**GetQuestionnaire**](docs/QuestionnaireApi.md#getquestionnaire) | **Get** /questionnaires/{questionnaireID} | 
*QuestionnaireApi* | [**GetQuestionnaires**](docs/QuestionnaireApi.md#getquestionnaires) | **Get** /questionnaires | 
*QuestionnaireApi* | [**GetQuestions**](docs/QuestionnaireApi.md#getquestions) | **Get** /questionnaires/{questionnaireID}/questions | 
*QuestionnaireApi* | [**PatchQuestionnaire**](docs/QuestionnaireApi.md#patchquestionnaire) | **Patch** /questionnaires/{questionnaireID} | 
*QuestionnaireApi* | [**PostQuestionnaire**](docs/QuestionnaireApi.md#postquestionnaire) | **Post** /questionnaires | 
*ResponseApi* | [**DeleteResponse**](docs/ResponseApi.md#deleteresponse) | **Delete** /responses/{responseID} | 
*ResponseApi* | [**GetResponses**](docs/ResponseApi.md#getresponses) | **Get** /responses/{responseID} | 
*ResponseApi* | [**PatchResponse**](docs/ResponseApi.md#patchresponse) | **Patch** /responses/{responseID} | 
*ResponseApi* | [**PostResponse**](docs/ResponseApi.md#postresponse) | **Post** /responses | 
*ResultApi* | [**GetResults**](docs/ResultApi.md#getresults) | **Get** /results/{questionnaireID} | 
*UserApi* | [**GetMe**](docs/UserApi.md#getme) | **Get** /users/me | 
*UserApi* | [**GetMyAdministrates**](docs/UserApi.md#getmyadministrates) | **Get** /users/me/administrates | 
*UserApi* | [**GetMyResponses**](docs/UserApi.md#getmyresponses) | **Get** /users/me/responses | 
*UserApi* | [**GetMyTargeted**](docs/UserApi.md#getmytargeted) | **Get** /users/me/targeted | 
*UserApi* | [**GetResponsesToQuestionnaire**](docs/UserApi.md#getresponsestoquestionnaire) | **Get** /users/me/responses/{questionnaireID} | 
*UserApi* | [**GetUsers**](docs/UserApi.md#getusers) | **Get** /users | 未実装


## Documentation For Models

 - [Group](docs/Group.md)
 - [Me](docs/Me.md)
 - [NewQuestion](docs/NewQuestion.md)
 - [NewQuestionnaire](docs/NewQuestionnaire.md)
 - [NewQuestionnaireResponse](docs/NewQuestionnaireResponse.md)
 - [NewResponse](docs/NewResponse.md)
 - [Question](docs/Question.md)
 - [QuestionAllOf](docs/QuestionAllOf.md)
 - [QuestionDetails](docs/QuestionDetails.md)
 - [QuestionDetailsAllOf](docs/QuestionDetailsAllOf.md)
 - [Questionnaire](docs/Questionnaire.md)
 - [QuestionnaireById](docs/QuestionnaireById.md)
 - [QuestionnaireByIdAllOf](docs/QuestionnaireByIdAllOf.md)
 - [QuestionnaireForList](docs/QuestionnaireForList.md)
 - [QuestionnaireForListAllOf](docs/QuestionnaireForListAllOf.md)
 - [QuestionnaireMyAdministrates](docs/QuestionnaireMyAdministrates.md)
 - [QuestionnaireMyAdministratesAllOf](docs/QuestionnaireMyAdministratesAllOf.md)
 - [QuestionnaireMyTargeted](docs/QuestionnaireMyTargeted.md)
 - [QuestionnaireMyTargetedAllOf](docs/QuestionnaireMyTargetedAllOf.md)
 - [QuestionnaireUser](docs/QuestionnaireUser.md)
 - [QuestionnaireUserAllOf](docs/QuestionnaireUserAllOf.md)
 - [Response](docs/Response.md)
 - [ResponseAllOf](docs/ResponseAllOf.md)
 - [ResponseBody](docs/ResponseBody.md)
 - [ResponseDetails](docs/ResponseDetails.md)
 - [ResponseDetailsAllOf](docs/ResponseDetailsAllOf.md)
 - [ResponseResult](docs/ResponseResult.md)
 - [ResponseResultAllOf](docs/ResponseResultAllOf.md)
 - [ResponseSummary](docs/ResponseSummary.md)
 - [User](docs/User.md)


## Documentation For Authorization



## application


- **Type**: OAuth
- **Flow**: application
- **Authorization URL**: 
- **Scopes**: 
 - **write**: allows modifying resources
 - **read**: allows reading resources

Example

```golang
auth := context.WithValue(context.Background(), sw.ContextAccessToken, "ACCESSTOKENSTRING")
r, err := client.Service.Operation(auth, args)
```

Or via OAuth2 module to automatically refresh tokens and perform user authentication.

```golang
import "golang.org/x/oauth2"

/* Perform OAuth2 round trip request and obtain a token */

tokenSource := oauth2cfg.TokenSource(createContext(httpClient), &token)
auth := context.WithValue(oauth2.NoContext, sw.ContextOAuth2, tokenSource)
r, err := client.Service.Operation(auth, args)
```



## Author



