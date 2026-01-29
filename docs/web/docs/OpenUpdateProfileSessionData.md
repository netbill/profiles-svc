# OpenUpdateProfileSessionData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | Upload session id | 
**Type** | **string** |  | 
**Attributes** | [**OpenUpdateProfileSessionDataAttributes**](OpenUpdateProfileSessionDataAttributes.md) |  | 
**Relationships** | [**OpenUpdateProfileSessionDataRelationships**](OpenUpdateProfileSessionDataRelationships.md) |  | 

## Methods

### NewOpenUpdateProfileSessionData

`func NewOpenUpdateProfileSessionData(id uuid.UUID, type_ string, attributes OpenUpdateProfileSessionDataAttributes, relationships OpenUpdateProfileSessionDataRelationships, ) *OpenUpdateProfileSessionData`

NewOpenUpdateProfileSessionData instantiates a new OpenUpdateProfileSessionData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenUpdateProfileSessionDataWithDefaults

`func NewOpenUpdateProfileSessionDataWithDefaults() *OpenUpdateProfileSessionData`

NewOpenUpdateProfileSessionDataWithDefaults instantiates a new OpenUpdateProfileSessionData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *OpenUpdateProfileSessionData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *OpenUpdateProfileSessionData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *OpenUpdateProfileSessionData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *OpenUpdateProfileSessionData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *OpenUpdateProfileSessionData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *OpenUpdateProfileSessionData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *OpenUpdateProfileSessionData) GetAttributes() OpenUpdateProfileSessionDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *OpenUpdateProfileSessionData) GetAttributesOk() (*OpenUpdateProfileSessionDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *OpenUpdateProfileSessionData) SetAttributes(v OpenUpdateProfileSessionDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *OpenUpdateProfileSessionData) GetRelationships() OpenUpdateProfileSessionDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *OpenUpdateProfileSessionData) GetRelationshipsOk() (*OpenUpdateProfileSessionDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *OpenUpdateProfileSessionData) SetRelationships(v OpenUpdateProfileSessionDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


