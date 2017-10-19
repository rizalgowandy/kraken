package blobserver

import (
	"hash"
	"time"

	"code.uber.internal/infra/kraken/lib/hrw"
	"code.uber.internal/infra/kraken/utils/memsize"
	"github.com/spaolacci/murmur3"
)

// Config defines the configuration used by Origin cluster for hashing blob digests.
type Config struct {
	NumReplica int          `yaml:"num_replica"`
	HashNodes  HashNodeMap  `yaml:"hash_nodes"`
	Repair     RepairConfig `yaml:"repair"`
}

// HashNodeMap defines a map from address of HashNodeConfig
// NOTE: Address should be in the form of hostname:port or dns name.
// The reason behind this is that we do not want to make assumption of the port
// each origin is running on.
type HashNodeMap map[string]HashNodeConfig

// HashNodeConfig defines the config for a single origin node
type HashNodeConfig struct {
	Label  string `yaml:"label"`
	Weight int    `yaml:"weight"`
}

// RepairConfig defines the configuration used by Origin cluster for
// running concurrent repairs.
type RepairConfig struct {
	NumWorkers    int           `yaml:"num_workers"`
	MaxRetries    uint          `yaml:"max_retries"`
	RetryDelay    time.Duration `yaml:"retry_delay"`
	MaxRetryDelay time.Duration `yaml:"max_retry_delay"`
	ConnTimeout   time.Duration `yaml:"conn_timeout"`
}

// ClientConfig defines configuration for blobserver HTTP Client.
type ClientConfig struct {
	UploadChunkSize int64 `yaml:"upload_chunk_size"`
}

func (c ClientConfig) applyDefaults() ClientConfig {
	if c.UploadChunkSize == 0 {
		return ClientConfig{
			UploadChunkSize: int64(50 * memsize.MB),
		}
	}
	return c
}

// LabelToAddress generates a reverse mapping of HashNodes by label to hostname.
func (c Config) LabelToAddress() map[string]string {
	m := make(map[string]string, len(c.HashNodes))
	for addr, node := range c.HashNodes {
		m[node.Label] = addr
	}
	return m
}

// HashState initializes hash state from the configured hash nodes.
func (c Config) HashState() *hrw.RendezvousHash {
	h := hrw.NewRendezvousHash(
		func() hash.Hash { return murmur3.New64() },
		hrw.UInt64ToFloat64)
	for _, node := range c.HashNodes {
		h.AddNode(node.Label, node.Weight)
	}
	return h
}