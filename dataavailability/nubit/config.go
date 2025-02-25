package nubit

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// NubitNamespaceBytesLength is the fixed-size bytes array.
const NubitNamespaceBytesLength = 58

// NubitMinCommitTime is the minimum commit time interval between blob submissions to NubitDA.
const NubitMinCommitTime time.Duration = 12 * time.Second

// Config is the NubitDA backend configurations
type Config struct {
	NubitRpcURL             string         `mapstructure:"NubitRpcURL"`
	NubitAuthKey            string         `mapstructure:"NubitAuthKey"`
	NubitNamespace          string         `mapstructure:"NubitNamespace"`
	NubitGetProofMaxRetry   uint64         `mapstructure:"NubitGetProofMaxRetry"`
	NubitGetProofWaitPeriod types.Duration `mapstructure:"NubitGetProofWaitPeriod"`
}
