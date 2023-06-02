package grpcproxy

import (
	"encoding/json"
	"github.com/cj1105/etcd-proxy/proxy/grpcproxy/namespacequota"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	PathNamespaceQuotaConfig = "/proxy/nsq"
)

func HandleNamespaceQuotaConfig(lg *zap.Logger, mux *http.ServeMux, c *clientv3.Client) {
	if lg == nil {
		lg = zap.NewNop()
	}
	mux.HandleFunc(PathNamespaceQuotaConfig, func(w http.ResponseWriter, r *http.Request) {
		nsqStore := namespacequota.NewNamespaceQuotaManager(lg, c, namespacequota.NamespaceQuotaManagerConfig{NamespaceQuotaEnforcement: 2}, "client")
		switch r.Method {
		case "POST":
			var err error
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			quota := &namespacequota.NamespaceQuota{}
			err = json.Unmarshal(body, quota)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			_, err = nsqStore.SetNamespaceQuota(quota.Key, quota.QuotaKeyCount, quota.QuotaByteCount)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			lg.Debug("/health OK", zap.Int("status-code", http.StatusOK))
		case "GET":
			if name := r.URL.Query().Get("name"); name != "" {
				rule, err := nsqStore.GetNamespaceQuota(name)
				if err != nil {
					http.Error(w, "", http.StatusInternalServerError)
				}
				respBody, _ := json.Marshal(rule)
				w.Write(respBody)
			} else {
				quotas := nsqStore.ListNamespaceQuotas()
				respBody, _ := json.Marshal(quotas)
				w.Write(respBody)
			}
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	})
}
