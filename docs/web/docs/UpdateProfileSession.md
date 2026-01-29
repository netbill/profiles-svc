# UpdateProfileSession

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**UpdateProfileSessionData**](UpdateProfileSessionData.md) |  | 
**Included** | [**[]ProfileData**](ProfileData.md) |  | 

## Methods

### NewUpdateProfileSession

`func NewUpdateProfileSession(data UpdateProfileSessionData, included []ProfileData, ) *UpdateProfileSession`

NewUpdateProfileSession instantiates a new UpdateProfileSession object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateProfileSessionWithDefaults

`func NewUpdateProfileSessionWithDefaults() *UpdateProfileSession`

NewUpdateProfileSessionWithDefaults instantiates a new UpdateProfileSession object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *UpdateProfileSession) GetData() UpdateProfileSessionData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *UpdateProfileSession) GetDataOk() (*UpdateProfileSessionData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *UpdateProfileSession) SetData(v UpdateProfileSessionData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *UpdateProfileSession) GetIncluded() []ProfileData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *UpdateProfileSession) GetIncludedOk() (*[]ProfileData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *UpdateProfileSession) SetIncluded(v []ProfileData)`

SetIncluded sets Included field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


