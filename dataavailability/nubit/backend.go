package nubit

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	share "github.com/RiemaLabs/nubit-node/da"
	client "github.com/RiemaLabs/nubit-node/rpc/rpc/client"
	nodeBlob "github.com/RiemaLabs/nubit-node/strucs/btx"
	"github.com/RiemaLabs/nubit-validator/da/namespace"
	"github.com/ethereum/go-ethereum/common"
)

const (
	DefaultFetchTimeout  = time.Minute
	DefaultSubmitTimeout = time.Minute
)

type NubitDABackend struct {
	ctx            context.Context
	ns             namespace.Namespace
	client         *client.Client
	SubmintTimeout time.Duration
	FetchTimeout   time.Duration
}

func NewNubitDABackend(node_rpc, auth_token, np string) (*NubitDABackend, error) {

	cn, err := client.NewClient(context.TODO(), node_rpc, auth_token)
	if err != nil {
		return nil, err
	}
	name := namespace.MustNewV0([]byte(np))

	log.Infof("⚙️     Nubit Namespace : %s ", string(name.ID))
	return &NubitDABackend{
		ns:             name,
		client:         cn,
		ctx:            context.Background(),
		SubmintTimeout: DefaultSubmitTimeout,
		FetchTimeout:   DefaultFetchTimeout,
	}, nil
}

func (a *NubitDABackend) Init() error {
	return nil
}

// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
// as expected by the contract
func (a *NubitDABackend) PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error) {
	encodedData, err := MarshalBatchData(batchesData)
	if err != nil {
		log.Errorf("🏆    NubitDABackend.MarshalBatchData:%s", err)
		return encodedData, err
	}

	nsp, err := share.NamespaceFromBytes(a.ns.Bytes())
	if nil != err {
		log.Errorf("🏆    NubitDABackend.NamespaceFromBytes:%s", err)
		return nil, err
	}

	body, err := nodeBlob.NewBlobV0(nsp, encodedData)
	if nil != err {
		log.Errorf("🏆    NubitDABackend.NewBlobV0:%s", err)
		return nil, err
	}

	log.Infof("🏆     Nubit send data:%+v", body)

	ctx, cancel := context.WithTimeout(ctx, a.SubmintTimeout)
	blockNumber, err := a.client.Blob.Submit(ctx, []*nodeBlob.Blob{body}, 0.01)
	cancel()
	if err != nil {
		log.Errorf("🏆    NubitDABackend.Submit:%s", err)
		return nil, err
	}

	var batchDAData BatchDAData
	copy(batchDAData.Commitment[:], body.Commitment)

	batchDAData.BlockNumber = int64(blockNumber)
	log.Infof("🏆  Nubit prepared DA data:%+v", batchDAData)

	// todo: use bridge API data
	returnData, err := batchDAData.Encode()
	if err != nil {
		return nil, fmt.Errorf("🏆  Nubit cannot encode batch data:%w", err)
	}

	log.Infof("🏆  Nubit Data submitted by sequencer:%d bytes against namespace %v sent with height %#x", len(encodedData), a.ns, blockNumber)

	return returnData, nil
}

func (a *NubitDABackend) GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error) {
	var batchDAData BatchDAData
	err := batchDAData.Decode(dataAvailabilityMessage)
	if err != nil {
		log.Errorf("🏆    NubitDABackend.GetSequence.Decode:%s", err)
		return nil, err
	}
	log.Infof("🏆     Nubit GetSequence batchDAData:%+v", batchDAData)
	ctx, cancel := context.WithTimeout(ctx, a.FetchTimeout)
	blob, err := a.client.Blob.Get(context.TODO(), uint64(batchDAData.BlockNumber), a.ns.Bytes(), batchDAData.Commitment[:])
	cancel()
	if err != nil {
		log.Errorf("🏆    NubitDABackend.GetSequence.Blob.Get:%s", err)
		return nil, err
	}
	log.Infof("🏆     Nubit GetSequence blob.data:%+v", blob.GetData())
	return UnmarshalBatchData(blob.GetData())
}
