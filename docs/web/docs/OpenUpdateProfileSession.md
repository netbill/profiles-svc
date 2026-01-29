# OpenUpdateProfileSession

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**OpenUpdateProfileSessionData**](OpenUpdateProfileSessionData.md) |  | 
**Included** | [**[]ProfileData**](ProfileData.md) |  | 

## Methods

### NewOpenUpdateProfileSession

`func NewOpenUpdateProfileSession(data OpenUpdateProfileSessionData, included []ProfileData, ) *OpenUpdateProfileSession`

NewOpenUpdateProfileSession instantiates a new OpenUpdateProfileSession object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenUpdateProfileSessionWithDefaults

`func NewOpenUpdateProfileSessionWithDefaults() *OpenUpdateProfileSession`

NewOpenUpdateProfileSessionWithDefaults instantiates a new OpenUpdateProfileSession object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *OpenUpdateProfileSession) GetData() OpenUpdateProfileSessionData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *OpenUpdateProfileSession) GetDataOk() (*OpenUpdateProfileSessionData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *OpenUpdateProfileSession) SetData(v OpenUpdateProfileSessionData)`

SetData sets Data field to given value.


### GetIncluded

`func (o *OpenUpdateProfileSession) GetIncluded() []ProfileData`

GetIncluded returns the Included field if non-nil, zero value otherwise.

### GetIncludedOk

`func (o *OpenUpdateProfileSession) GetIncludedOk() (*[]ProfileData, bool)`

GetIncludedOk returns a tuple with the Included field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncluded

`func (o *OpenUpdateProfileSession) SetIncluded(v []ProfileData)`

SetIncluded sets Included field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


