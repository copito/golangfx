package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv1 "github.com/copito/runner/idl_gen/common/v1"
)

func EncodeToBase64(data any) (string, error) {
	// Serialize the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", nil
	}

	// Encode the JSON data to Base64
	base64String := base64.StdEncoding.EncodeToString(jsonData)
	return base64String, nil
}

func DecodeFromBase64(base64String string, data any) error {
	// Decode the Base64 string
	jsonData, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, data)
}

func EncodePaginationToBase64(pagination Pagination) (string, error) {
	data, err := json.Marshal(pagination)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pagination: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func DecodePaginationFromBase64(pageToken string) (Pagination, error) {
	// Decode the base64-encided token
	decoded, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return Pagination{}, fmt.Errorf("invalid page token: %w", err)
	}

	// Unmarshal the JSON into a Pagination object
	var pagination Pagination
	err = json.Unmarshal(decoded, &pagination)
	if err != nil {
		return Pagination{}, fmt.Errorf("invalid page token format: %w", err)
	}

	return pagination, nil
}

func handleDefaults(defaultPageSize int, defaultPageNumber int) (defaultPageSizeResult int, defaultPageNumberResult int) {
	if defaultPageSize <= 0 {
		defaultPageSizeResult = int(DEFAULT_LIMIT)
	}
	if defaultPageSize > int(MAX_PAGE_SIZE) {
		defaultPageSizeResult = int(MAX_PAGE_SIZE)
	}
	if defaultPageNumber < 0 {
		defaultPageNumberResult = 1
	}

	return defaultPageSizeResult, defaultPageNumberResult
}

func ParsePaginationRequest(paginationRequest *commonv1.PaginationRequest, defaultPageSize int, defaultPageNumber int) (Pagination, error) {
	defaultPageSize, defaultPageNumber = handleDefaults(defaultPageSize, defaultPageNumber)

	// Handle nul case
	if paginationRequest == nil {
		return Pagination{
			PageSize:   int32(defaultPageSize),
			PageNumber: int32(defaultPageNumber),
		}, nil
	}

	pageToken := paginationRequest.GetPageToken()
	// Decode page token if provided
	if pageToken == "" {
		// No page token - return defaults
		return Pagination{
			PageSize:   int32(defaultPageSize),
			PageNumber: int32(defaultPageNumber),
		}, nil
	}

	decodedToken, err := DecodePaginationFromBase64(pageToken)
	if err != nil {
		return Pagination{}, fmt.Errorf("failed to decode page token: %w", err)
	}

	if decodedToken.PageSize <= 0 {
		return Pagination{}, status.Error(codes.InvalidArgument, "page_size must be greater than 0")
	}

	if decodedToken.PageSize > MAX_PAGE_SIZE {
		return Pagination{}, status.Errorf(codes.InvalidArgument, "page_size must be less than or equal to %d", MAX_PAGE_SIZE)
	}

	if decodedToken.PageNumber < 0 {
		return Pagination{}, status.Error(codes.InvalidArgument, "page_number must be greater than or equal to 0")
	}

	return decodedToken, nil
}
