# UpdateProfileSessionData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Upload session id | 
**Type** | **string** |  | 
**Attributes** | [**UpdateProfileSessionDataAttributes**](UpdateProfileSessionDataAttributes.md) |  | 
**Relationships** | [**UpdateProfileSessionDataRelationships**](UpdateProfileSessionDataRelationships.md) |  | 

## Methods

### NewUpdateProfileSessionData

`func NewUpdateProfileSessionData(id uuid.UUID, type_ string, attributes UpdateProfileSessionDataAttributes, relationships UpdateProfileSessionDataRelationships, ) *UpdateProfileSessionData`

NewUpdateProfileSessionData instantiates a new UpdateProfileSessionData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateProfileSessionDataWithDefaults

`func NewUpdateProfileSessionDataWithDefaults() *UpdateProfileSessionData`

NewUpdateProfileSessionDataWithDefaults instantiates a new UpdateProfileSessionData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UpdateProfileSessionData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateProfileSessionData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateProfileSessionData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UpdateProfileSessionData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UpdateProfileSessionData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UpdateProfileSessionData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UpdateProfileSessionData) GetAttributes() UpdateProfileSessionDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UpdateProfileSessionData) GetAttributesOk() (*UpdateProfileSessionDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UpdateProfileSessionData) SetAttributes(v UpdateProfileSessionDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *UpdateProfileSessionData) GetRelationships() UpdateProfileSessionDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *UpdateProfileSessionData) GetRelationshipsOk() (*UpdateProfileSessionDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *UpdateProfileSessionData) SetRelationships(v UpdateProfileSessionDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


