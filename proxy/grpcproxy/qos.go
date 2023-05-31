package grpcproxy

import (
	"encoding/json"
	"github.com/cj1105/etcd-proxy/proxy/grpcproxy/qos"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	PathQosRuleConfig = "/proxy/qos"
)

func HandleQosRuleConfig(lg *zap.Logger, mux *http.ServeMux, c *clientv3.Client) {
	if lg == nil {
		lg = zap.NewNop()
	}
	mux.HandleFunc(PathQosRuleConfig, func(w http.ResponseWriter, r *http.Request) {
		qosStore := qos.NewQoSStore(lg, c, "client")
		switch r.Method {
		case "POST":
			var err error
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			rule := &qos.QoSRule{}
			err = json.Unmarshal(body, rule)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			err = qosStore.QoSRuleAdd(rule)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			lg.Debug("/health OK", zap.Int("status-code", http.StatusOK))
		case "GET":
			if name := r.URL.Query().Get("name"); name != "" {
				rule, err := qosStore.QoSRuleGet(name)
				if err != nil {
					http.Error(w, "", http.StatusInternalServerError)
				}
				respBody, _ := json.Marshal(rule)
				w.Write(respBody)
			} else {
				rules, err := qosStore.QoSRuleList()
				if err != nil {
					http.Error(w, "", http.StatusInternalServerError)
				}
				respBody, _ := json.Marshal(rules)
				w.Write(respBody)
			}
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	})
}
