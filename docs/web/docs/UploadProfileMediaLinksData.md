# UploadProfileMediaLinksData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | [**uuid.UUID**](uuid.UUID.md) | account id | 
**Type** | **string** |  | 
**Attributes** | [**UploadProfileMediaLinksDataAttributes**](UploadProfileMediaLinksDataAttributes.md) |  | 
**Relationships** | [**UploadProfileMediaLinksDataRelationships**](UploadProfileMediaLinksDataRelationships.md) |  | 

## Methods

### NewUploadProfileMediaLinksData

`func NewUploadProfileMediaLinksData(id uuid.UUID, type_ string, attributes UploadProfileMediaLinksDataAttributes, relationships UploadProfileMediaLinksDataRelationships, ) *UploadProfileMediaLinksData`

NewUploadProfileMediaLinksData instantiates a new UploadProfileMediaLinksData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadProfileMediaLinksDataWithDefaults

`func NewUploadProfileMediaLinksDataWithDefaults() *UploadProfileMediaLinksData`

NewUploadProfileMediaLinksDataWithDefaults instantiates a new UploadProfileMediaLinksData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *UploadProfileMediaLinksData) GetId() uuid.UUID`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UploadProfileMediaLinksData) GetIdOk() (*uuid.UUID, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UploadProfileMediaLinksData) SetId(v uuid.UUID)`

SetId sets Id field to given value.


### GetType

`func (o *UploadProfileMediaLinksData) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *UploadProfileMediaLinksData) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *UploadProfileMediaLinksData) SetType(v string)`

SetType sets Type field to given value.


### GetAttributes

`func (o *UploadProfileMediaLinksData) GetAttributes() UploadProfileMediaLinksDataAttributes`

GetAttributes returns the Attributes field if non-nil, zero value otherwise.

### GetAttributesOk

`func (o *UploadProfileMediaLinksData) GetAttributesOk() (*UploadProfileMediaLinksDataAttributes, bool)`

GetAttributesOk returns a tuple with the Attributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributes

`func (o *UploadProfileMediaLinksData) SetAttributes(v UploadProfileMediaLinksDataAttributes)`

SetAttributes sets Attributes field to given value.


### GetRelationships

`func (o *UploadProfileMediaLinksData) GetRelationships() UploadProfileMediaLinksDataRelationships`

GetRelationships returns the Relationships field if non-nil, zero value otherwise.

### GetRelationshipsOk

`func (o *UploadProfileMediaLinksData) GetRelationshipsOk() (*UploadProfileMediaLinksDataRelationships, bool)`

GetRelationshipsOk returns a tuple with the Relationships field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRelationships

`func (o *UploadProfileMediaLinksData) SetRelationships(v UploadProfileMediaLinksDataRelationships)`

SetRelationships sets Relationships field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


