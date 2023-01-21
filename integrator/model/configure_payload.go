package model

// ConfigPayload defines model for ConfigPayload.
type ConfigPayload struct {
	// Action Action to be performed
	Action string `json:"action"`

	// Data Data for the action
	Data map[string]interface{} `json:"data"`

	// PluginName Plugin name for the operation
	PluginName string `json:"plugin_name"`

	// Resource Resource to be configured
	Resource string `json:"resource"`
}
