package middleware

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	constants "github.com/lakhansamani/cloud-container/internal/contants"
	messages "github.com/lakhansamani/cloud-container/internal/messages"
)

// Define a struct to parse the GraphQL request
type GraphQLRequest struct {
	Query         interface{} `json:"query"`
	Variables     interface{} `json:"variables"`
	OperationName string      `json:"operationName"`
}

// List of queries/mutations that require authentication along with associated roles
var protectedOperations = map[string][]string{
	"create_deployment": {constants.RoleTypeDeploymentAdmin},
	"delete_deployment": {constants.RoleTypeDeploymentAdmin},
	"session":           {constants.RoleTypeDeploymentAdmin},
}

// List of public queries/mutations
var publicOperations = []string{
	"login",
	"logout",
	"verify_otp",
}

// AuthorizationMiddleware is a gin middleware that performs authorization for incoming requests
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the path and check if it's "/v1/graphql"
		path := c.Request.URL.Path
		if path != "/v1/graphql" {
			c.Next()
			return
		}
		// Read the request body
		var gqlRequest GraphQLRequest
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Debug().Err(err).Msg("error reading request body")
			c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
			return
		}
		err = json.Unmarshal(jsonData, &gqlRequest)
		if err != nil {
			log.Debug().Err(err).Msg("error unmarshalling json")
			c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
			return
		}
		// Bind the request body to GraphQLRequest struct
		c.Request.Body = io.NopCloser(strings.NewReader(string(jsonData)))
		operationName := strings.ToLower(gqlRequest.OperationName)
		// Check if the operation is IntrospectionQuery and allow if it is, required for graphql playground docs
		if isIntrospectionQuery(gqlRequest.Query.(string), operationName) {
			c.Next()
			return
		}
		// Check if operationName matches the actual name in the mutation or query
		if !doesOperationNameMatchQuery(gqlRequest.Query.(string), operationName) {
			log.Debug().Str("operationName", operationName).Msg("operation name mismatch")
			c.AbortWithStatusJSON(400, gin.H{"message": messages.InvalidGraphqlOperationMessage})
			return
		}
		// Check if the operation is a public mutation
		if isPublicOperation(operationName) {
			c.Next()
			return
		}
		// If the operation is protected, perform authentication
		// if allowedRoles, ok := protectedOperations[operationName]; ok {
		// 	// Get access token from authorization header
		// 	accessToken, err := token.GetAccessToken(c)
		// 	if err != nil || accessToken == "" {
		// 		log.Debug().Err(err).Msg("error getting access token from authorization header")
		// 		c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
		// 		return
		// 	}
		// 	// Validate access token
		// 	accessTokenData, err := token.ValidateAccessToken(c, accessToken)
		// 	if err != nil || accessTokenData == nil {
		// 		log.Debug().Err(err).Msg("error validating access token")
		// 		c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
		// 		return
		// 	}
		// 	// Check if user has the required roles
		// 	if !hasRequiredRole(accessTokenData.Roles, allowedRoles) {
		// 		log.Debug().Str("roles", accessTokenData.Roles).Err(err).Msg("user does not have the required role")
		// 		c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
		// 		return
		// 	}
		// 	c.Set("user_id", accessTokenData.Subject)
		// 	c.Set("user_roles", accessTokenData.Roles)
		// 	c.Set("company_id", accessTokenData.CompanyID)
		// 	if val, ok := accessTokenData.HasuraClaims["x-hasura-company-role-id"]; ok {
		// 		c.Set("company_role_id", val.(string))
		// 	}
		// 	c.Next()
		// 	return
		// }
		c.Next()
		// Operation not found, abort
		// c.AbortWithStatusJSON(400, gin.H{"message": messages.InvalidGraphqlOperationMessage})
	}
}

// hasRequiredRole checks if one of the allowed role is present in the user roles present as part of jwt tokens
func hasRequiredRole(userRoles string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if strings.Contains(userRoles, role) {
			return true
		}
	}
	return false
}

// isPublicOperation checks if the operation is a public mutation/query
func isPublicOperation(operationName string) bool {
	for _, mutation := range publicOperations {
		if operationName == mutation {
			return true
		}
	}
	return false
}

// doesOperationNameMatchQuery checks if the operationName matches any of the actual names in the mutation or query
func doesOperationNameMatchQuery(query string, operationName string) bool {
	// Regex to extract all the names of the mutations or queries after the opening brace '{'
	r := regexp.MustCompile(`\{\s*(\w+)`)
	matches := r.FindAllStringSubmatch(query, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		actualName := match[1]
		if strings.ToLower(actualName) == operationName {
			return true
		}
	}
	return false
}

// Check if the request is for IntrospectionQuery
func isIntrospectionQuery(query string, operationName string) bool {
	// Check if the operation name is "IntrospectionQuery"
	if operationName != "introspectionquery" {
		return false
	}
	// Further check a substring from the query to confirm
	return strings.Contains(query, "__schema")
}
