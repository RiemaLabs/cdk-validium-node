package dataavailability

// DABackendType is the data availability protocol for the CDK
type DABackendType string

const (
	// DataAvailabilityCommittee is the DAC protocol backend
	DataAvailabilityCommittee DABackendType = "DataAvailabilityCommittee"
	Nubit                     DABackendType = "Nubit"
)

type Config struct {
	NubitNodeRpcURL string `mapstructure:"NubitNodeRpcURL"`
	AuthKey         string `mapstructure:"AuthKey"`
	AppName         string `mapstructure:"AppName"`
}
