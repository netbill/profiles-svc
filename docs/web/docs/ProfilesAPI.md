# \ProfilesAPI

All URIs are relative to *http://localhost:8002*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProfilesSvcV1ProfilesAccountIdGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesAccountIdGet) | **Get** /profiles-svc/v1/profiles/{account_id} | Get profile by account id
[**ProfilesSvcV1ProfilesAccountIdOfficialPatch**](ProfilesAPI.md#ProfilesSvcV1ProfilesAccountIdOfficialPatch) | **Patch** /profiles-svc/v1/profiles/{account_id}/official | Update profile official status
[**ProfilesSvcV1ProfilesGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesGet) | **Get** /profiles-svc/v1/profiles/ | Filter profiles
[**ProfilesSvcV1ProfilesMeGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeGet) | **Get** /profiles-svc/v1/profiles/me/ | Get my profile
[**ProfilesSvcV1ProfilesMeUpdateSessionAvatarDelete**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeUpdateSessionAvatarDelete) | **Delete** /profiles-svc/v1/profiles/me/update-session/avatar/ | Delete uploaded avatar in session
[**ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch) | **Patch** /profiles-svc/v1/profiles/me/update-session/confirm/ | Update my profile
[**ProfilesSvcV1ProfilesMeUpdateSessionPost**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeUpdateSessionPost) | **Post** /profiles-svc/v1/profiles/me/update-session/ | Open profile update session
[**ProfilesSvcV1ProfilesUUsernameGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesUUsernameGet) | **Get** /profiles-svc/v1/profiles/u/{username} | Get profile by username



## ProfilesSvcV1ProfilesAccountIdGet

> Profile ProfilesSvcV1ProfilesAccountIdGet(ctx, accountId).Execute()

Get profile by account id



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Account id (UUID).

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesAccountIdGet(context.Background(), accountId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesAccountIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesAccountIdGet`: Profile
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesAccountIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **uuid.UUID** | Account id (UUID). | 

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesAccountIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Profile**](Profile.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesAccountIdOfficialPatch

> Profile ProfilesSvcV1ProfilesAccountIdOfficialPatch(ctx, accountId).UpdateProfileOfficial(updateProfileOfficial).Execute()

Update profile official status



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	accountId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // uuid.UUID | Account id (UUID).
	updateProfileOfficial := *openapiclient.NewUpdateProfileOfficial(*openapiclient.NewUpdateProfileOfficialData("TODO", "Type_example", *openapiclient.NewUpdateProfileOfficialDataAttributes(false))) // UpdateProfileOfficial | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesAccountIdOfficialPatch(context.Background(), accountId).UpdateProfileOfficial(updateProfileOfficial).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesAccountIdOfficialPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesAccountIdOfficialPatch`: Profile
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesAccountIdOfficialPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**accountId** | **uuid.UUID** | Account id (UUID). | 

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesAccountIdOfficialPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **updateProfileOfficial** | [**UpdateProfileOfficial**](UpdateProfileOfficial.md) |  | 

### Return type

[**Profile**](Profile.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesGet

> ProfilesCollection ProfilesSvcV1ProfilesGet(ctx).UsernameLike(usernameLike).Pseudonym(pseudonym).Limit(limit).Offset(offset).Execute()

Filter profiles



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	usernameLike := "usernameLike_example" // string | Prefix filter for username. (optional)
	pseudonym := "pseudonym_example" // string | Prefix filter for pseudonym. (optional)
	limit := int32(56) // int32 | Maximum number of items to return. (optional)
	offset := int32(56) // int32 | Number of items to skip. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesGet(context.Background()).UsernameLike(usernameLike).Pseudonym(pseudonym).Limit(limit).Offset(offset).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesGet`: ProfilesCollection
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **usernameLike** | **string** | Prefix filter for username. | 
 **pseudonym** | **string** | Prefix filter for pseudonym. | 
 **limit** | **int32** | Maximum number of items to return. | 
 **offset** | **int32** | Number of items to skip. | 

### Return type

[**ProfilesCollection**](ProfilesCollection.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesMeGet

> Profile ProfilesSvcV1ProfilesMeGet(ctx).Execute()

Get my profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMeGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMeGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesMeGet`: Profile
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesMeGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMeGetRequest struct via the builder pattern


### Return type

[**Profile**](Profile.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesMeUpdateSessionAvatarDelete

> ProfilesSvcV1ProfilesMeUpdateSessionAvatarDelete(ctx).Execute()

Delete uploaded avatar in session



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionAvatarDelete(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionAvatarDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMeUpdateSessionAvatarDeleteRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch

> Profile ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch(ctx).UpdateProfile(updateProfile).Execute()

Update my profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	updateProfile := *openapiclient.NewUpdateProfile(*openapiclient.NewUpdateProfileData("TODO", "Type_example", *openapiclient.NewUpdateProfileDataAttributes())) // UpdateProfile | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch(context.Background()).UpdateProfile(updateProfile).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch`: Profile
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionConfirmPatch`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMeUpdateSessionConfirmPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateProfile** | [**UpdateProfile**](UpdateProfile.md) |  | 

### Return type

[**Profile**](Profile.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesMeUpdateSessionPost

> UploadProfileMediaLinks ProfilesSvcV1ProfilesMeUpdateSessionPost(ctx).Execute()

Open profile update session



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesMeUpdateSessionPost`: UploadProfileMediaLinks
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesMeUpdateSessionPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMeUpdateSessionPostRequest struct via the builder pattern


### Return type

[**UploadProfileMediaLinks**](UploadProfileMediaLinks.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesUUsernameGet

> Profile ProfilesSvcV1ProfilesUUsernameGet(ctx, username).Execute()

Get profile by username



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	username := "username_example" // string | Username of the profile owner.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesUUsernameGet(context.Background(), username).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesUUsernameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesUUsernameGet`: Profile
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesUUsernameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**username** | **string** | Username of the profile owner. | 

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesUUsernameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Profile**](Profile.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json, application/problem+json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

