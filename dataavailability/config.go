package dataavailability

// DABackendType is the data availability protocol for the CDK
type DABackendType string

const (
	// DataAvailabilityCommittee is the DAC protocol backend
	DataAvailabilityCommittee DABackendType = "DataAvailabilityCommittee"
	Nubit                     DABackendType = "Nubit"
)

type Config struct {
	NodeRpc   string `mapstructure:"NodeRpc"`
	AuthToken string `mapstructure:"AuthToken"`
	Namespace string `mapstructure:"Namespace"`
}
