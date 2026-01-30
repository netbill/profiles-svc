# UpdateProfileSessionDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UploadToken** | **string** | JWT upload session token | 
**UploadUrl** | **string** | Pre-signed PUT URL for avatar upload | 
**GetUrl** | **string** | Pre-signed GET URL to read uploaded avatar | 

## Methods

### NewUpdateProfileSessionDataAttributes

`func NewUpdateProfileSessionDataAttributes(uploadToken string, uploadUrl string, getUrl string, ) *UpdateProfileSessionDataAttributes`

NewUpdateProfileSessionDataAttributes instantiates a new UpdateProfileSessionDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateProfileSessionDataAttributesWithDefaults

`func NewUpdateProfileSessionDataAttributesWithDefaults() *UpdateProfileSessionDataAttributes`

NewUpdateProfileSessionDataAttributesWithDefaults instantiates a new UpdateProfileSessionDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUploadToken

`func (o *UpdateProfileSessionDataAttributes) GetUploadToken() string`

GetUploadToken returns the UploadToken field if non-nil, zero value otherwise.

### GetUploadTokenOk

`func (o *UpdateProfileSessionDataAttributes) GetUploadTokenOk() (*string, bool)`

GetUploadTokenOk returns a tuple with the UploadToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadToken

`func (o *UpdateProfileSessionDataAttributes) SetUploadToken(v string)`

SetUploadToken sets UploadToken field to given value.


### GetUploadUrl

`func (o *UpdateProfileSessionDataAttributes) GetUploadUrl() string`

GetUploadUrl returns the UploadUrl field if non-nil, zero value otherwise.

### GetUploadUrlOk

`func (o *UpdateProfileSessionDataAttributes) GetUploadUrlOk() (*string, bool)`

GetUploadUrlOk returns a tuple with the UploadUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadUrl

`func (o *UpdateProfileSessionDataAttributes) SetUploadUrl(v string)`

SetUploadUrl sets UploadUrl field to given value.


### GetGetUrl

`func (o *UpdateProfileSessionDataAttributes) GetGetUrl() string`

GetGetUrl returns the GetUrl field if non-nil, zero value otherwise.

### GetGetUrlOk

`func (o *UpdateProfileSessionDataAttributes) GetGetUrlOk() (*string, bool)`

GetGetUrlOk returns a tuple with the GetUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGetUrl

`func (o *UpdateProfileSessionDataAttributes) SetGetUrl(v string)`

SetGetUrl sets GetUrl field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


