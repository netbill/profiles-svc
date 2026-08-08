# ProfilesCollection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | [**[]ProfileData**](ProfileData.md) |  | 
**Links** | [**PaginationData**](PaginationData.md) |  | 

## Methods

### NewProfilesCollection

`func NewProfilesCollection(data []ProfileData, links PaginationData, ) *ProfilesCollection`

NewProfilesCollection instantiates a new ProfilesCollection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProfilesCollectionWithDefaults

`func NewProfilesCollectionWithDefaults() *ProfilesCollection`

NewProfilesCollectionWithDefaults instantiates a new ProfilesCollection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *ProfilesCollection) GetData() []ProfileData`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *ProfilesCollection) GetDataOk() (*[]ProfileData, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *ProfilesCollection) SetData(v []ProfileData)`

SetData sets Data field to given value.


### GetLinks

`func (o *ProfilesCollection) GetLinks() PaginationData`

GetLinks returns the Links field if non-nil, zero value otherwise.

### GetLinksOk

`func (o *ProfilesCollection) GetLinksOk() (*PaginationData, bool)`

GetLinksOk returns a tuple with the Links field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinks

`func (o *ProfilesCollection) SetLinks(v PaginationData)`

SetLinks sets Links field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


