# \ProfilesAPI

All URIs are relative to *http://localhost:8002*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProfilesSvcV1ProfilesAccountIdGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesAccountIdGet) | **Get** /profiles-svc/v1/profiles/{account_id} | Get profile by account id
[**ProfilesSvcV1ProfilesGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesGet) | **Get** /profiles-svc/v1/profiles/ | Filter profiles
[**ProfilesSvcV1ProfilesMeGet**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeGet) | **Get** /profiles-svc/v1/profiles/me/ | Get my profile
[**ProfilesSvcV1ProfilesMeMediaDelete**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeMediaDelete) | **Delete** /profiles-svc/v1/profiles/me/media/ | Delete profile upload media
[**ProfilesSvcV1ProfilesMeMediaPost**](ProfilesAPI.md#ProfilesSvcV1ProfilesMeMediaPost) | **Post** /profiles-svc/v1/profiles/me/media/ | Create profile upload media link
[**ProfilesSvcV1ProfilesMePatch**](ProfilesAPI.md#ProfilesSvcV1ProfilesMePatch) | **Patch** /profiles-svc/v1/profiles/me/ | Update my profile
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


## ProfilesSvcV1ProfilesGet

> ProfilesCollection ProfilesSvcV1ProfilesGet(ctx).Text(text).PageLimit(pageLimit).PageOffset(pageOffset).Execute()

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
	text := "text_example" // string | Text to filter profiles by. Matches against `username` and `pseudonym` fields using prefix-based filtering.  (optional)
	pageLimit := int32(56) // int32 | Max number of items to return (1-100). (optional)
	pageOffset := int32(56) // int32 | Number of items to skip. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesGet(context.Background()).Text(text).PageLimit(pageLimit).PageOffset(pageOffset).Execute()
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
 **text** | **string** | Text to filter profiles by. Matches against &#x60;username&#x60; and &#x60;pseudonym&#x60; fields using prefix-based filtering.  | 
 **pageLimit** | **int32** | Max number of items to return (1-100). | 
 **pageOffset** | **int32** | Number of items to skip. | 

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


## ProfilesSvcV1ProfilesMeMediaDelete

> ProfilesSvcV1ProfilesMeMediaDelete(ctx).DeleteUploadProfileAvatar(deleteUploadProfileAvatar).Execute()

Delete profile upload media



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
	deleteUploadProfileAvatar := *openapiclient.NewDeleteUploadProfileAvatar(*openapiclient.NewDeleteUploadProfileAvatarData("TODO", "Type_example", *openapiclient.NewDeleteUploadProfileAvatarDataAttributes())) // DeleteUploadProfileAvatar | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMeMediaDelete(context.Background()).DeleteUploadProfileAvatar(deleteUploadProfileAvatar).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMeMediaDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMeMediaDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deleteUploadProfileAvatar** | [**DeleteUploadProfileAvatar**](DeleteUploadProfileAvatar.md) |  | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesMeMediaPost

> UploadProfileMediaLinks ProfilesSvcV1ProfilesMeMediaPost(ctx).Execute()

Create profile upload media link



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
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMeMediaPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMeMediaPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesMeMediaPost`: UploadProfileMediaLinks
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesMeMediaPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMeMediaPostRequest struct via the builder pattern


### Return type

[**UploadProfileMediaLinks**](UploadProfileMediaLinks.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProfilesSvcV1ProfilesMePatch

> Profile ProfilesSvcV1ProfilesMePatch(ctx).UpdateProfile(updateProfile).Execute()

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
	resp, r, err := apiClient.ProfilesAPI.ProfilesSvcV1ProfilesMePatch(context.Background()).UpdateProfile(updateProfile).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProfilesAPI.ProfilesSvcV1ProfilesMePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ProfilesSvcV1ProfilesMePatch`: Profile
	fmt.Fprintf(os.Stdout, "Response from `ProfilesAPI.ProfilesSvcV1ProfilesMePatch`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProfilesSvcV1ProfilesMePatchRequest struct via the builder pattern


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

