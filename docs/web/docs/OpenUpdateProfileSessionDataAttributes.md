# OpenUpdateProfileSessionDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UploadToken** | **string** | JWT upload session token | 
**UploadUrl** | **string** | Pre-signed PUT URL for avatar upload | 
**GetUrl** | **string** | Pre-signed GET URL to read uploaded avatar | 
**ExpiresAt** | **time.Time** | Upload session expiration time | 

## Methods

### NewOpenUpdateProfileSessionDataAttributes

`func NewOpenUpdateProfileSessionDataAttributes(uploadToken string, uploadUrl string, getUrl string, expiresAt time.Time, ) *OpenUpdateProfileSessionDataAttributes`

NewOpenUpdateProfileSessionDataAttributes instantiates a new OpenUpdateProfileSessionDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenUpdateProfileSessionDataAttributesWithDefaults

`func NewOpenUpdateProfileSessionDataAttributesWithDefaults() *OpenUpdateProfileSessionDataAttributes`

NewOpenUpdateProfileSessionDataAttributesWithDefaults instantiates a new OpenUpdateProfileSessionDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUploadToken

`func (o *OpenUpdateProfileSessionDataAttributes) GetUploadToken() string`

GetUploadToken returns the UploadToken field if non-nil, zero value otherwise.

### GetUploadTokenOk

`func (o *OpenUpdateProfileSessionDataAttributes) GetUploadTokenOk() (*string, bool)`

GetUploadTokenOk returns a tuple with the UploadToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadToken

`func (o *OpenUpdateProfileSessionDataAttributes) SetUploadToken(v string)`

SetUploadToken sets UploadToken field to given value.


### GetUploadUrl

`func (o *OpenUpdateProfileSessionDataAttributes) GetUploadUrl() string`

GetUploadUrl returns the UploadUrl field if non-nil, zero value otherwise.

### GetUploadUrlOk

`func (o *OpenUpdateProfileSessionDataAttributes) GetUploadUrlOk() (*string, bool)`

GetUploadUrlOk returns a tuple with the UploadUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadUrl

`func (o *OpenUpdateProfileSessionDataAttributes) SetUploadUrl(v string)`

SetUploadUrl sets UploadUrl field to given value.


### GetGetUrl

`func (o *OpenUpdateProfileSessionDataAttributes) GetGetUrl() string`

GetGetUrl returns the GetUrl field if non-nil, zero value otherwise.

### GetGetUrlOk

`func (o *OpenUpdateProfileSessionDataAttributes) GetGetUrlOk() (*string, bool)`

GetGetUrlOk returns a tuple with the GetUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGetUrl

`func (o *OpenUpdateProfileSessionDataAttributes) SetGetUrl(v string)`

SetGetUrl sets GetUrl field to given value.


### GetExpiresAt

`func (o *OpenUpdateProfileSessionDataAttributes) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *OpenUpdateProfileSessionDataAttributes) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *OpenUpdateProfileSessionDataAttributes) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


