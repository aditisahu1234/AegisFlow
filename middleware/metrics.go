package middleware

import "api-gateway/models"

var GlobalMetrics = models.Metrics{}		//global metrcis store
/* above creates: 
Metrics{
	TotalRequests: 0,
	AllowedRequests: 0,
	BlockedRequests: 0,
}  .......................when server starts 
*/