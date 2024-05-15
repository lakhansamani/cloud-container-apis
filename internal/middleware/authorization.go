package middleware

import (
	"encoding/json"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	constants "github.com/lakhansamani/cloud-container/internal/contants"
	messages "github.com/lakhansamani/cloud-container/internal/messages"
	"github.com/lakhansamani/cloud-container/internal/session"
)

// Define a struct to parse the GraphQL request
type GraphQLRequest struct {
	Query         interface{} `json:"query"`
	Variables     interface{} `json:"variables"`
	OperationName string      `json:"operationName"`
}

// List of queries/mutations that require authentication
var protectedOperations = []string{
	"create_deployment",
	"delete_deployment",
	"logout",
	"session",
	"deployments",
	"deployment",
}

// List of public queries/mutations
var publicOperations = []string{
	"signup",
	"login",
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
		// If the operation is protected, check if operation exists in protectedOperations
		// Get the cookie and decrypt the session token
		// Add user_id to the context
		if isProtectedOperation(operationName) {
			sessionCookie, err := c.Request.Cookie(constants.SessionCookieName)
			if err != nil {
				log.Debug().Err(err).Msg("error getting session cookie")
				c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
				return
			}
			sessionValue, err := url.PathUnescape(sessionCookie.Value)
			if err != nil {
				log.Debug().Err(err).Msg("error unescaping session value")
				c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
				return
			}
			// Decrypt session token
			userID, _, err := session.DecryptSession(sessionValue)
			if err != nil {
				log.Debug().Err(err).Msg("error decrypting session token")
				c.AbortWithStatusJSON(400, gin.H{"message": messages.UnauthorizedError})
				return
			}
			c.Set("user_id", userID)
			c.Next()
			return
		} else {
			c.AbortWithStatusJSON(400, gin.H{"message": messages.InvalidGraphqlOperationMessage})
			return
		}
	}
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

// isProtectedOperation checks if the operation is a public mutation/query
func isProtectedOperation(operationName string) bool {
	for _, mutation := range protectedOperations {
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
