// Copyright 2021 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package namespacequota

import (
	"context"
	"encoding/json"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"math"
	"strings"
	"sync"
	"time"
)

type NamespaceQuotaEnforcement int8

const (
	DISABLED NamespaceQuotaEnforcement = 0
	SOFTMODE NamespaceQuotaEnforcement = 1
	HARDMODE NamespaceQuotaEnforcement = 2
)

var (
	// namespaceQuotaBucketName defines the name of bucket in the backend
	namespaceQuotaBucketName       = "namespacequota"
	ErrNamespaceQuotaExceeded      = errors.New("namespace quota exceeded")
	ErrNamespaceQuotaNotFound      = errors.New("namespace quota not found")
	ErrNamespaceQuotaRestoreFailed = errors.New("namespace quota restore failed")
)

// NamespaceQuota represents a namespace quota
type NamespaceQuota struct {
	// Key represents the namespace key
	Key string `json:"key"`
	// QuotaByteCount quota byte count represents the byte quota for respective Key
	QuotaByteCount uint64 `json:"quotaByteCount"`
	// QuotaKeyCount quota key count represents the quota key quota for respective Key
	QuotaKeyCount uint64 `json:"quotaKeyCount"`
	// UsageByteCount represents the current usage of bytes for respective Key
	UsageByteCount uint64 `json:"usageByteCount"`
	// UsageKeyCount represents the current usage of bytes for respective Key
	UsageKeyCount uint64 `json:"usageKeyCount"`
}

// ReadView is a ReadView that only permits reads. Defined here
// to avoid circular dependency with mvcc.
type ReadView interface {
	RangeValueSize(tenant string) (keySize, valueSize int64, err error)
	GetValueSize(key string) (valueSize int, found bool)
}

type readView struct{ kv clientv3.KV }

func (rv *readView) RangeValueSize(tenant string) (keySize, valueSize int64, err error) {
	resp, err := rv.kv.Get(context.TODO(), tenant, clientv3.WithPrefix())
	if err != nil {
		return 0, 0, err
	}
	keySize = int64(0)
	valueSize = int64(0)
	for i := range resp.Kvs {
		keySize = keySize + int64(len(resp.Kvs[i].Key))
		valueSize = valueSize + int64(len(resp.Kvs[i].Value))
	}
	return keySize, valueSize, nil
}

func (rv *readView) GetValueSize(key string) (valueSize int, found bool) {
	resp, err := rv.kv.Get(context.TODO(), key)
	if err != nil {
		return 0, false
	}
	if len(resp.Kvs) == 0 {
		return 0, false
	}
	return len(resp.Kvs[0].Value), true
}

// NamespaceQuotaManagerConfig represents the namespace quota manager configuration
type NamespaceQuotaManagerConfig struct {
	NamespaceQuotaEnforcement int
}

// NamespaceQuotaManager owns NamespaceQuota. NamespaceQuotaManager can create, update, read, and delete NamespaceQuota.
type NamespaceQuotaManager interface {

	// SetReadView lets the namespace quota manager create ReadView to the store.
	// NamespaceQuotaManager ranges over the items
	SetReadView(rv ReadView)

	// SetNamespaceQuota creates or updates NamespaceQuota
	// If quota exists, it will be updated
	// Else created
	SetNamespaceQuota(key string, quotaByteCount uint64, quotaKeyCount uint64) (quota *NamespaceQuota, err error)

	// GetNamespaceQuota returns the NamespaceQuota on the key if it exists, else responds with an error
	GetNamespaceQuota(key string) (quota *NamespaceQuota, err error)

	// ListNamespaceQuotas returns all the NamespaceQuota available in the tree
	ListNamespaceQuotas() (quotas []*NamespaceQuota)

	// DeleteNamespaceQuota deletes the NamespaceQuota if exists
	DeleteNamespaceQuota(key string) (quota *NamespaceQuota, err error)

	// UpdateNamespaceQuotaUsage updates the usage of a key
	UpdateNamespaceQuotaUsage(key string, byteDiff int, keyDiff int)

	// DeleteRangeUpdateUsage will update usage in the quota tree when an etcd delete range operation is performed
	DeleteRangeUpdateUsage(keys [][]byte, valueSizes []int)

	// DiffNamespaceQuotaUsage diffs namespace quota usage, if incoming new key exceeds quota, error is returned
	DiffNamespaceQuotaUsage(key, newValue []byte) (byteDiff int, keyDiff int)

	// IsQuotaExceeded checks the quota based on the byteDiff/keyDiff, and enforcement status
	IsQuotaExceeded(key string, byteDiff int, keyDiff int) bool

	// Recover recovers the namespace quota manager
	Recover(c *clientv3.Client, rv ReadView)

	// Promote promotes the namespace quota manager to be primary
	Promote()

	// Demote demotes the namespace quota manager from being the primary.
	Demote()
}

// namespaceQuotaManager implements NamespaceQuotaManager interface.
type namespaceQuotaManager struct {
	mu sync.RWMutex

	// rv get value sizes over a range
	rv ReadView

	// namespaceQuotaTree stores the map of a key to the Namespace struct
	namespaceQuotaTree QuotaIndex

	// namespaceQuotaEnforcement namespace quota enforcement
	namespaceQuotaEnforcement NamespaceQuotaEnforcement

	// namespaceQuotaExceedMap map of quotas that have exceeded map
	namespaceQuotaExceedMap map[string]string

	// b backend to persist namespace quotas
	// persists the quota and usage along with the ID
	c *clientv3.Client

	// demotec is set when the namespace quota manager is the primary.
	// demotec will be closed if the namespace quota manager is demoted.
	demotec chan struct{}

	// lg zap logger
	lg *zap.Logger
}

func (nqm *namespaceQuotaManager) SetReadView(rv ReadView) {
	nqm.mu.Lock()
	defer nqm.mu.Unlock()
	nqm.rv = rv
}

func (nqm *namespaceQuotaManager) SetNamespaceQuota(key string, quotaByteCount uint64, quotaKeyCount uint64) (*NamespaceQuota, error) {
	// Step 1: Check if the quota exist already, if so this is an update, else a create
	treeItem := nqm.namespaceQuotaTree.Get(key)

	namespaceQuota := &NamespaceQuota{Key: key}

	// Step 2a: If quota already exist, just update the newly set fields
	if treeItem != nil {
		if nqm.lg != nil {
			nqm.lg.Debug("updating a quota, current quota values",
				zap.String("key", treeItem.Key),
				zap.Uint64("quota byte count", treeItem.QuotaByteCount),
				zap.Uint64("quota key count", treeItem.QuotaKeyCount),
				zap.Uint64("usage byte count", treeItem.UsageByteCount),
				zap.Uint64("usage key count", treeItem.UsageKeyCount))
		}

		// Populate the existing fields that can't change in a set operation
		namespaceQuota.UsageByteCount = treeItem.UsageByteCount
		namespaceQuota.UsageKeyCount = treeItem.UsageKeyCount

		// 0 means the quota is not set by the user via the client.
		// So if any of the incoming fields are 0, just pull them from the existing quota
		// If the user is setting the quota to 0, then they should simply delete the quota, since at that point
		// the namespace owner wouldn't be able to add anything.
		// The user can set:
		// 1. Byte Quota Count
		if quotaByteCount == 0 {
			quotaByteCount = treeItem.QuotaByteCount
		}

		// 2. Key Quota Count
		if quotaKeyCount == 0 {
			quotaKeyCount = treeItem.QuotaKeyCount
		}
	} else {
		// 2b. value is empty, create a new quota
		if nqm.lg != nil {
			nqm.lg.Debug("creating new quota")
		}

		// If either of the quotas are absent, default them to having a very large value
		if quotaByteCount == 0 {
			quotaByteCount = math.MaxUint64
		}
		if quotaKeyCount == 0 {
			quotaKeyCount = math.MaxUint64
		}

		// Fetch the size of the values stored, and add them up
		// The range value size will fetch everything across the range and return the list of values
		keySize, valueSize := nqm.RangeValueSize(key)

		// Fill up the missing fields
		namespaceQuota.UsageByteCount = uint64(keySize + valueSize)
		// TODO: Does the key quota count really need to be uint64?
		namespaceQuota.UsageKeyCount = 0
	}

	namespaceQuota.QuotaByteCount = quotaByteCount
	namespaceQuota.QuotaKeyCount = quotaKeyCount

	if namespaceQuota.QuotaByteCount < namespaceQuota.UsageByteCount || namespaceQuota.QuotaKeyCount < namespaceQuota.UsageKeyCount {
		if nqm.lg != nil {
			nqm.lg.Warn("quota exceeded at the time of setting quota", zap.ByteString("key", []byte(namespaceQuota.Key)))
		}
		nqm.namespaceQuotaExceedMap[string(namespaceQuota.Key)] = string(namespaceQuota.Key)
	}

	// override the existing quota with updated values
	_, isUpdated := nqm.namespaceQuotaTree.Put(namespaceQuota)
	if nqm.lg != nil {
		nqm.lg.Debug("new quota values",
			zap.ByteString("key", []byte(namespaceQuota.Key)),
			zap.Uint64("quota byte count", namespaceQuota.QuotaByteCount),
			zap.Uint64("quota key count", namespaceQuota.QuotaKeyCount),
			zap.Uint64("usage byte count", namespaceQuota.UsageByteCount),
			zap.Uint64("usage key count", namespaceQuota.UsageKeyCount))
	}

	nqm.mu.Lock()
	namespaceQuota.persistTo(nqm.c)
	nqm.mu.Unlock()
	if isUpdated {
		// TODO: update with a namespace quota updated metric instead of create
	} else {
		namespaceQuotaCreatedCounter.Inc()
	}
	return namespaceQuota, nil
}

func (nqm *namespaceQuotaManager) GetNamespaceQuota(key string) (*NamespaceQuota, error) {
	// Check if quota exist, else return error
	currQuota := nqm.namespaceQuotaTree.Get(key)
	if currQuota != nil {
		return currQuota, nil
	}
	return nil, ErrNamespaceQuotaNotFound
}

func (nqm *namespaceQuotaManager) ListNamespaceQuotas() []*NamespaceQuota {
	return nqm.namespaceQuotaTree.GetAllNamespaceQuotas()
}

func (nqm *namespaceQuotaManager) DiffNamespaceQuotaUsage(key, newValue []byte) (byteDiff int, keyDiff int) {
	if nqm.namespaceQuotaEnforcement == DISABLED {
		return 0, 0
	}
	start := time.Now()
	// Step 0: get keys and value sizes
	// a nil response indicates the absence of key in the key tree, and is a new addition to the key tree
	valueSize, isFound := nqm.rv.GetValueSize(string(key))

	// Step 1a: set defaults as if a new key was supposed to be
	// increment key usage count by 1
	// increment byte usage by len(newValue)
	keyDiff = 1
	byteDiff = len(newValue)

	// Step 1b: If key exist, no key count changes
	// update the byteCount
	if isFound {
		// Step 1b: the key exist, no changes to key count, check diff for value size
		keyDiff = 0
		// remove the current size and add the new incoming valueSize
		byteDiff = len(newValue) - valueSize
		return byteDiff, keyDiff
	}
	// Step 1c: if the key does not exist, a new key will be created
	namespaceQuotaGetDiffValueDurationSec.Observe(time.Since(start).Seconds())
	return byteDiff, keyDiff
}

func (nqm *namespaceQuotaManager) IsQuotaExceeded(key string, byteDiff int, keyDiff int) bool {
	if nqm.namespaceQuotaEnforcement == DISABLED {
		return false
	}

	// Step 1: Check quota path for any keys exceeding quotas and return an error if any quotas are exceeded
	exceededQuota := nqm.namespaceQuotaTree.CheckQuotaPath(key, uint64(byteDiff), uint64(keyDiff))
	if exceededQuota == nil {
		return false
	}

	if nqm.namespaceQuotaEnforcement == SOFTMODE {
		if nqm.lg != nil {
			value, ok := nqm.namespaceQuotaExceedMap[string(exceededQuota.Key)]
			// 3a. If key not present and quota exceeded, add to map, and log
			if !ok && (exceededQuota.QuotaByteCount < exceededQuota.UsageByteCount || exceededQuota.QuotaKeyCount < exceededQuota.UsageKeyCount) {
				nqm.lg.Warn("quota exceeded but allowed", zap.String("key", exceededQuota.Key))
				nqm.namespaceQuotaExceedMap[string(exceededQuota.Key)] = key
			} else if ok && (exceededQuota.QuotaByteCount >= exceededQuota.UsageByteCount && exceededQuota.QuotaKeyCount >= exceededQuota.UsageKeyCount) {
				// 3b. If key present and quota reduced, remove from map and log
				nqm.lg.Warn("quota back in limits", zap.String("key", exceededQuota.Key))
				delete(nqm.namespaceQuotaExceedMap, value)
				delete(nqm.namespaceQuotaExceedMap, exceededQuota.Key)
			}
		}
		return false
	}

	if nqm.lg != nil {
		nqm.lg.Debug("quota violated", zap.String("key", key))
	}

	return true
}

func (nqm *namespaceQuotaManager) UpdateNamespaceQuotaUsage(key string, byteDiff int, keyDiff int) {
	if nqm.namespaceQuotaEnforcement == DISABLED {
		return
	}
	nqm.namespaceQuotaTree.UpdateQuotaPath(key, byteDiff, keyDiff)
}

func (nqm *namespaceQuotaManager) DeleteNamespaceQuota(key string) (quota *NamespaceQuota, err error) {
	deletedQuota := nqm.namespaceQuotaTree.Delete(key)
	if deletedQuota != nil {
		namespaceQuotaDeletedCounter.Inc()
		return deletedQuota, nil
	}
	return nil, ErrNamespaceQuotaNotFound
}

func (nqm *namespaceQuotaManager) DeleteRangeUpdateUsage(keys [][]byte, valueSizes []int) {
	if nqm.namespaceQuotaEnforcement == DISABLED {
		return
	}
	// Walk through each key, and decrease the usage for each key
	for i := 0; i < len(keys); i++ {
		nqm.namespaceQuotaTree.UpdateQuotaPath(string(keys[i]), -1*(len(keys[i])+valueSizes[i]), -1)
	}
	nqm.lg.Debug("delete range, updating usage complete")
}

func (nqm *namespaceQuotaManager) Recover(c *clientv3.Client, rv ReadView) {
	nqm.mu.Lock()
	defer nqm.mu.Unlock()

	nqm.c = c
	nqm.rv = rv
	// TODO: Restore the tree
	nqm.initAndRecover()
}

func (nqm *namespaceQuotaManager) Promote() {
	nqm.mu.Lock()
	defer nqm.mu.Unlock()
	nqm.demotec = make(chan struct{})
}

func (nqm *namespaceQuotaManager) Demote() {
	nqm.mu.Lock()
	defer nqm.mu.Unlock()
	if nqm.demotec != nil {
		close(nqm.demotec)
		nqm.demotec = nil
	}
}

// persistTo persists data to backend
func (nq *NamespaceQuota) persistTo(c *clientv3.Client) {
	nqpb := NamespaceQuota{
		Key:            nq.Key,
		QuotaByteCount: nq.QuotaByteCount,
		QuotaKeyCount:  nq.QuotaKeyCount,
		UsageByteCount: nq.UsageByteCount,
		UsageKeyCount:  nq.UsageKeyCount,
	}
	nqpdBytes, err := json.Marshal(nqpb)
	if err != nil {
		return
	}

	c.Put(context.TODO(), strings.Join([]string{namespaceQuotaBucketName, nqpb.Key}, "/"), string(nqpdBytes))
}

// NewNamespaceQuotaManager returns a new NamespaceQuotaManager
func NewNamespaceQuotaManager(lg *zap.Logger, c *clientv3.Client, cfg NamespaceQuotaManagerConfig, mode string) NamespaceQuotaManager {
	return newNamespaceQuotaManager(lg, c, cfg, mode)
}

// newNamespaceQuotaManager initializes/recovers NamespaceQuotaManager
func newNamespaceQuotaManager(lg *zap.Logger, c *clientv3.Client, cfg NamespaceQuotaManagerConfig, mode string) *namespaceQuotaManager {
	nqm := &namespaceQuotaManager{
		namespaceQuotaTree:        NewQuotaTreeIndex(lg),
		rv:                        &readView{kv: c.KV},
		c:                         c,
		lg:                        lg,
		namespaceQuotaExceedMap:   make(map[string]string),
		namespaceQuotaEnforcement: SOFTMODE,
	}
	if mode != "client" {
		nqm.initAndRecover()
	}
	lg.Warn("namespace quota manager config", zap.Int64("enforcement status", int64(nqm.namespaceQuotaEnforcement)))
	namespaceQuotaEnforcementStatus.Set(float64(nqm.namespaceQuotaEnforcement))
	return nqm
}

// initAndRecover init and recover the NamespaceQuotaManager
func (nqm *namespaceQuotaManager) initAndRecover() {
	quotas := mustUnsafeGetAllNamespaceQuotas(nqm.c)
	for _, q := range quotas {
		keySize, valueSize, err := nqm.rv.RangeValueSize(q.Key)
		if err != nil {
			return
		}
		nqm.namespaceQuotaTree.Put(&NamespaceQuota{
			Key:            q.Key,
			QuotaKeyCount:  q.QuotaKeyCount,
			QuotaByteCount: q.QuotaByteCount,
			UsageKeyCount:  q.UsageKeyCount,
			UsageByteCount: uint64(keySize + valueSize),
		})
	}
}

func mustUnsafeGetAllNamespaceQuotas(c *clientv3.Client) []*NamespaceQuota {
	resp, _ := c.Get(context.TODO(), namespaceQuotaBucketName, clientv3.WithPrefix())
	nqs := make([]*NamespaceQuota, 0, len(resp.Kvs))
	for i := range resp.Kvs {
		nqpb := &NamespaceQuota{}
		err := json.Unmarshal(resp.Kvs[i].Value, nqpb)
		if err != nil {
			panic("failed to unmarshal namespace quota proto item")
		}
		nqs = append(nqs, nqpb)
	}
	return nqs
}

// RangeValueSize ranges the value size
func (nqm *namespaceQuotaManager) RangeValueSize(tenant string) (keySize, valueSize int64) {
	keySize, valueSize, err := nqm.rv.RangeValueSize(tenant)
	if err != nil {
		return 0, 0
	}
	return keySize, valueSize
}
