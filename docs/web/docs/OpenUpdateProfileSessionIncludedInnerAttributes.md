# OpenUpdateProfileSessionIncludedInnerAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Username** | **string** | Username | 
**Pseudonym** | Pointer to **string** | Pseudonym | [optional] 
**Description** | Pointer to **string** | Description | [optional] 
**Official** | **bool** | Is official account | 
**Avatar** | Pointer to **string** | Avatar URL | [optional] 
**CreatedAt** | **time.Time** | Created at | 
**UpdatedAt** | **time.Time** | Updated at | 

## Methods

### NewOpenUpdateProfileSessionIncludedInnerAttributes

`func NewOpenUpdateProfileSessionIncludedInnerAttributes(username string, official bool, createdAt time.Time, updatedAt time.Time, ) *OpenUpdateProfileSessionIncludedInnerAttributes`

NewOpenUpdateProfileSessionIncludedInnerAttributes instantiates a new OpenUpdateProfileSessionIncludedInnerAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenUpdateProfileSessionIncludedInnerAttributesWithDefaults

`func NewOpenUpdateProfileSessionIncludedInnerAttributesWithDefaults() *OpenUpdateProfileSessionIncludedInnerAttributes`

NewOpenUpdateProfileSessionIncludedInnerAttributesWithDefaults instantiates a new OpenUpdateProfileSessionIncludedInnerAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsername

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetUsername(v string)`

SetUsername sets Username field to given value.


### GetPseudonym

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetPseudonym() string`

GetPseudonym returns the Pseudonym field if non-nil, zero value otherwise.

### GetPseudonymOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetPseudonymOk() (*string, bool)`

GetPseudonymOk returns a tuple with the Pseudonym field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPseudonym

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetPseudonym(v string)`

SetPseudonym sets Pseudonym field to given value.

### HasPseudonym

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) HasPseudonym() bool`

HasPseudonym returns a boolean if a field has been set.

### GetDescription

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetOfficial

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetOfficial() bool`

GetOfficial returns the Official field if non-nil, zero value otherwise.

### GetOfficialOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetOfficialOk() (*bool, bool)`

GetOfficialOk returns a tuple with the Official field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOfficial

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetOfficial(v bool)`

SetOfficial sets Official field to given value.


### GetAvatar

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetAvatar() string`

GetAvatar returns the Avatar field if non-nil, zero value otherwise.

### GetAvatarOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetAvatarOk() (*string, bool)`

GetAvatarOk returns a tuple with the Avatar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatar

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetAvatar(v string)`

SetAvatar sets Avatar field to given value.

### HasAvatar

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) HasAvatar() bool`

HasAvatar returns a boolean if a field has been set.

### GetCreatedAt

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.


### GetUpdatedAt

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *OpenUpdateProfileSessionIncludedInnerAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


