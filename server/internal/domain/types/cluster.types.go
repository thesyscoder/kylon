/**
 * @File: cluster.types.go
 * @Title: Cluster API Types
 * @Description: Defines data structures for API requests and responses related to Cluster management.
 */

package types

// RegisterClusterRequest represents the request body for registering a new cluster.
type RegisterClusterRequest struct {
	Name       string `json:"name" binding:"required"`       // Name of the cluster.
	Kubeconfig string `json:"kubeconfig" binding:"required"` // Kubeconfig content for accessing the cluster.
}
