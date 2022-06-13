package dockerregistry

type registryAPIError struct {
	code string
	message string
	detail string
}

var (
	ErrBlobUnknown   = &registryAPIError         { "BLOB_UNKNOWN", "blob unknown to registry", ""  }
	ErrBlobUploadUnknown = &registryAPIError     { "BLOB_UPLOAD_UNKNOWN", "blob upload unknown to registry", ""  }
	ErrDigestInvalid = &registryAPIError         { "DIGEST_INVALID", "provided digest did not match uploaded content", "" }
	ErrManifestBlobUnknown = &registryAPIError   { "MANIFEST_BLOB_UNKNOWN", "blob unknown to registry", ""  }
	ErrManifestInvalid = &registryAPIError       { "MANIFEST_INVALID", "manifest invalid", "" }
	ErrManifestUnknown = &registryAPIError       { "MANIFEST_UNKNOWN", "manifest unknown", "" }
	ErrNameInvalid = &registryAPIError           { "NAME_INVALID", "repository name not known to registry", "" }
	ErrBlobUploadInvalidSize = &registryAPIError { "SIZE_INVALID", "provided length did not match content length", "" }
	ErrTagInvalid = &registryAPIError            { "TAG_INVALID", "manifest tag did not match URI", "" }
	ErrUnauthorized = &registryAPIError          { "UNAUTHORIZED", "authentication required", "" }
	ErrDenied = &registryAPIError                { "DENIED", "requested access to the resource is denied", "" }
)

/*
	Endpoint								Error cases									Error Codes


	GET /v2/									On Failure: Authentication Required			UNAUTHORIZED
												On Failure: Too Many Requests				TOOMANYREQUESTS

	GET /v2/<name>/tags/list					On failure: Bad request						400 Bad Request, (NAME_INVALID, invalid repository name)
												On Failure: Authentication Required			401 Unauthorized UNAUTHORIZED, authentication required
												On Failure: No Such Repository Error		404, NAME_UNKNOWN
												On Failure: Access Denied					403 Forbidden, DENIED, requested access to the resource is denied
												On Failure: Too Many Requests				429, TOOMANYREQUESTS


	GET/HEAD /v2/<name>/manifests/<reference>	On Failure: Bad request						400 bad Request (NAME_INVALID, TAG_INVALID)
												On Failure: Authentication Required			401 Unauthorized UNAUTHORIZED, authentication required
												On Failure: Access Denied					403 Forbidden, DENIED, requested access to the resource is denied
												On Failure: No Such Repository Error		404 Not Found, NAME_UNKNOWN, repository name not known to registry

	PUT /v2/<name>/manifests/<reference>		On Failure: Invalid Manifest				400 Bad request (NAME_INVALID, TAG_INVALID, MANIFEST_INVALID)
												On Failure: Missing Layer(s)				400 Bad request (BLOB_UNKNOWN, Missing layers references in manifest)
												On Failure: Access Denied					403 Forbidden, DENIED, requested access to the resource is denied
												On Failure: No Such Repository Error		404 Not Found (NAME_UNKNOWN	repository name not known to registry)

	DELETE /v2/<name>manifests/<reference>		On Failure: Invalid Name or Reference		400 Bad Request (NAME_INVALID, TAG_INVALID)
												On Failure: Authentication Required			401 Unauthorized UNAUTHORIZED, authentication required
												On Failure: No Such Repository Error		404 Not Found, NAME_UNKNOWN, repository name not known to registry
												On Failure: Access Denied					403 Forbidden, DENIED, requested access to the resource is denied
												On Failure: Too Many Requests				429, TOOMANYREQUESTS
												On Failure: Unknown Manifest				404 Not Found (NAME_UNKNOWN, MANIFEST_UNKNOWN)

	GET /v2/<name>/blobs/<digest>				On Success: Temporary Redirect
												On Failure: Bad Request						400 Bad Request (NAME_INVALID, DIGEST_INVALID)
												On Failure: Not Found						404 Not Found (NAME_UNKNOWN, BLOB_UNKNOWN)
												On Failure: Authentication Required			401 Unauthorized UNAUTHORIZED, authentication required
												On Failure: No Such Repository Error		404 Not Found (NAME_UNKNOWN)
												On Failure: Access Denied					403 Forbidden, DENIED, requested access to the resource is denied
												On Failure: Too Many Requests				429, TOOMANYREQUESTS

	GET /v2/<name>/blobs/<digest>				On Success: Partial Content
												On Failure: Bad Request						400 Bad Request (NAME_INVALID, DIGEST_INVALID)
												On Failure: Not Found						404 Not Found (NAME_UNKNOWN, BLOB_UNKNOWN)
												On Failure: Authentication Required			401 Unauthorized UNAUTHORIZED, authentication required
												On Failure: No Such Repository Error		404 Not Found (NAME_UNKNOWN)
												On Failure: Access Denied					403 Forbidden, DENIED, requested access to the resource is denied
												On Failure: Too Many Requests				429, TOOMANYREQUESTS

	DELETE /v2/<name>/blobs/<digest>			On Success: Accepted
												On Failure: Not Found
												On Failure: Method Not Allowed
												On Failure: Authentication Required
												On Failure: No Such Repository Error
												On Failure: Access Denied
												On Failure: Too Many Requests

*/