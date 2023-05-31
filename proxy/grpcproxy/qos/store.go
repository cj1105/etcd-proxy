// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "qs IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qos

import (
	"context"
	"encoding/json"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"

	"go.uber.org/zap"
)

type Subject struct {
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

type QoSRule struct {
	RuleName    string   `json:"rule_name,omitempty"`
	RuleType    string   `json:"rule_type,omitempty"`
	Subject     *Subject `json:"subject,omitempty"`
	Qps         uint64   `json:"qps,omitempty"`
	Threshold   uint64   `json:"threshold,omitempty"`
	Condition   string   `json:"condition,omitempty"`
	Priority    uint64   `json:"priority,omitempty"`
	Ratelimiter string   `json:"ratelimiter,omitempty"`
}

var (
	enableFlagKey = "qosEnabled"
	qosEnabled    = "1"
	qosDisabled   = "0"

	ErrQoSNotEnabled       = errors.New("qos: qos is not enabled")
	ErrQoSRuleAlreadyExist = errors.New("qos: QoSRule already exists")
	ErrQoSRuleEmpty        = errors.New("qos: QoSRule name is empty")
	ErrQoSRuleNotFound     = errors.New("qos: QoSRule not found")
	ErrQoSRateExceeded     = errors.New("qos: rate exceeded ")
	ErrInvalidPriorityRule = errors.New("qos: invalid rule priority")
	ErrNoRuleMatched       = errors.New("qos: no rules are matched")
	ErrEnvLackOfKey        = errors.New("qos: lack of env key")

	qosBucketName     = "qos"
	qosRuleBucketName = "qosRules"
)

// QoSStore defines qos storage interface.
type QoSStore interface {

	// Recover recovers the state of qos store from the given backend
	Recover()

	// QoSEnable turns on the qos feature
	QoSEnable() error

	// QoSDisable turns on the qos feature
	QoSDisable() error

	// IsQoSEnabled returns true if the qos feature is enabled.
	IsQoSEnabled() bool

	// QoSRuleAdd adds a new QoSRule into the backend store
	QoSRuleAdd(rule *QoSRule) error

	// QoSRuleDelete deletes a QoSRule from the backend store
	QoSRuleDelete(ruleName string) error

	// QoSRuleUpdate updates a QoSRule into the backend store
	QoSRuleUpdate(rule *QoSRule) error

	// QoSRuleGet gets the detailed information of a QoSRules
	QoSRuleGet(ruleName string) error

	// QoSRuleList list the detailed information of all QoSRules
	QoSRuleList() ([]*QoSRule, error)

	// GetToken returns true if rate do not exceed,otherwise return false
	GetToken(r *RequestContext) bool

	// PutToken returns token to rate limiter, return true if returns succeed,otherwise return false
	PutToken(r *RequestContext) bool

	// Close does cleanup of QoSStore
	Close() error
}

type Qos struct {
	lg        *zap.Logger
	enabled   bool
	enabledMu sync.RWMutex
	enforcer  Enforcer
	c         *clientv3.Client
}

func (qs *Qos) QoSEnable() error {
	qs.enabledMu.Lock()
	defer qs.enabledMu.Unlock()
	if qs.enabled {
		qs.lg.Info("qos is already enabled; ignored qos enable request")
		return nil
	}
	if _, err := qs.c.Put(context.TODO(), enableFlagKey, qosEnabled); err != nil {
		return err
	}
	qs.lg.Info("enabled qos")
	return nil
}

func (qs *Qos) QoSDisable() error {
	qs.enabledMu.Lock()
	defer qs.enabledMu.Unlock()
	if !qs.enabled {
		return ErrQoSNotEnabled
	}
	if _, err := qs.c.Put(context.TODO(), enableFlagKey, qosDisabled); err != nil {
		return err
	}
	qs.enabled = false
	qs.lg.Info("disabled qos")
	return nil
}

func (qs *Qos) Close() error {
	qs.enabledMu.Lock()
	defer qs.enabledMu.Unlock()
	if !qs.enabled {
		return nil
	}
	return nil
}

func (qs *Qos) Recover() error {
	enabled := false
	if resp, err := qs.c.Get(context.TODO(), enableFlagKey); err != nil {
		return err
	} else {
		for _, ev := range resp.Kvs {
			if string(ev.Key) == enableFlagKey {
				enabled = true
			}
		}
	}
	qs.enabledMu.Lock()
	qs.enabled = enabled
	qs.enabledMu.Unlock()
	return nil
}

func (qs *Qos) QoSRuleAdd(rule *QoSRule) error {
	QoSRule := getQoSRule(qs.lg, qs.c, rule.RuleName)
	if QoSRule != nil {
		return ErrQoSRuleAlreadyExist
	}

	putQoSRule(qs.lg, qs.c, rule)
	qs.enforcer.SyncRule(rule)

	qs.lg.Info("added a QoSRule", zap.String("name", rule.RuleName))
	return nil
}

func (qs *Qos) QoSRuleUpdate(r *QoSRule) error {
	rule := getQoSRule(qs.lg, qs.c, r.RuleName)
	if rule == nil {
		return ErrQoSRuleNotFound
	}

	putQoSRule(qs.lg, qs.c, r)
	qs.enforcer.SyncRule(r)

	qs.lg.Info("updated a QoSRule", zap.String("name", r.RuleName))
	return nil
}

func (qs *Qos) QoSRuleDelete(r *QoSRule) error {
	rule := getQoSRule(qs.lg, qs.c, r.RuleName)
	if rule == nil {
		return ErrQoSRuleNotFound
	}

	delQoSRule(qs.c, r.RuleName)

	qs.enforcer.DeleteRule(r)

	qs.lg.Info(
		"deleted a QoSRule",
		zap.String("QoSRule-name", r.RuleName),
	)
	return nil
}

func (qs *Qos) QoSRuleGet(ruleName string) (*QoSRule, error) {
	rule := getQoSRule(qs.lg, qs.c, ruleName)

	if rule == nil {
		return nil, ErrQoSRuleNotFound
	}
	return rule, nil
}

func (qs *Qos) QoSRuleList() ([]*QoSRule, error) {
	if rules, err := getAllQoSRules(qs.lg, qs.c); err != nil {
		return nil, err
	} else {
		return rules, nil
	}
}

func getQoSRule(lg *zap.Logger, c *clientv3.Client, ruleName string) *QoSRule {
	var rule QoSRule
	if resp, err := c.Get(context.TODO(), ruleName, clientv3.WithLimit(1)); err != nil {
		return nil
	} else {
		err = json.Unmarshal(resp.Kvs[0].Value, &rule)
		if err != nil {
			return nil
		}
		return &rule
	}
}

func getAllQoSRules(lg *zap.Logger, c *clientv3.Client) ([]*QoSRule, error) {
	resp, err := c.Get(context.TODO(), qosRuleBucketName, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	QoSRules := make([]*QoSRule, len(resp.Kvs))
	for i := range resp.Kvs {
		rule := &QoSRule{}
		err := json.Unmarshal(resp.Kvs[i].Value, rule)
		if err != nil {
			lg.Panic("failed to unmarshal 'qospb.QoSRule'", zap.Error(err))
		}
		QoSRules[i] = rule
	}
	return QoSRules, nil
}

func putQoSRule(lg *zap.Logger, c *clientv3.Client, rule *QoSRule) {
	ruleStr, err := json.Marshal(rule)
	if err != nil {
		return
	}
	if err != nil {
		lg.Panic("failed to unmarshal 'qospb.QoSRule'", zap.Error(err))
	}
	c.Put(context.TODO(), rule.RuleName, string(ruleStr))
}

func delQoSRule(c *clientv3.Client, ruleName string) {
	c.Delete(context.TODO(), ruleName)
}

func (qs *Qos) IsQoSEnabled() bool {
	qs.enabledMu.RLock()
	defer qs.enabledMu.RUnlock()
	return qs.enabled
}

// NewQoSStore creates a new qosStore.
func NewQoSStore(lg *zap.Logger, c *clientv3.Client) *Qos {
	if lg == nil {
		lg = zap.NewNop()
	}
	enabled := false
	resp, err := c.Get(context.TODO(), enableFlagKey)
	if err != nil {
		return nil
	}
	if len(resp.Kvs) == 1 {
		if string(resp.Kvs[0].Value) == qosEnabled {
			enabled = true
		}
	}

	qs := &Qos{
		lg:      lg,
		enabled: enabled,
		c:       c,
	}

	qs.enforcer = NewQoSEnforcer(lg)
	qs.buildQoSRules()

	return qs
}

func (qs *Qos) buildQoSRules() {
	rules, _ := getAllQoSRules(qs.lg, qs.c)
	for i := 0; i < len(rules); i++ {
		qs.enforcer.SyncRule(rules[i])
	}
}

func (qs *Qos) GetToken(r *RequestContext) bool {
	if !qs.IsQoSEnabled() {
		return true
	}
	return qs.enforcer.GetToken(r)
}

func (qs *Qos) PutToken(r *RequestContext) bool {
	if !qs.IsQoSEnabled() {
		return true
	}
	return qs.enforcer.PutToken(r)
}
