package main

import (
	"encoding/json"
)

type TlsData struct {
	ProtocolVersion              string `json:"v,omitempty"`
	ClientRandom                 string `json:"cr,omitempty"`
	SessionKey                   string `json:"ssk,omitempty"`
	ClientEarlyTrafficSecret     string `json:"cets,omitempty"`
	ClientHandshakeTrafficSecret string `json:"chts,omitempty"`
	ServerHandshakeTrafficSecret string `json:"shts,omitempty"`
	ClientTrafficSecret0         string `json:"cts0,omitempty"`
	ServerTrafficSecret0         string `json:"sts0,omitempty"`
	EarlyExporterSecret          string `json:"ees,omitempty"`
	ExporterSecret               string `json:"es,omitempty"`
}

func (td *TlsData) Sprint() string {
	var res string
	switch td.ProtocolVersion {
	case "TLSv1.3":
		if td.ClientEarlyTrafficSecret != "" {
			res += "CLIENT_EARLY_TRAFFIC_SECRET " + td.ClientRandom + " " + td.ClientEarlyTrafficSecret + "\n"
		}
		if td.ClientHandshakeTrafficSecret != "" {
			res += "CLIENT_HANDSHAKE_TRAFFIC_SECRET " + td.ClientRandom + " " + td.ClientHandshakeTrafficSecret + "\n"
		}
		if td.ServerHandshakeTrafficSecret != "" {
			res += "SERVER_HANDSHAKE_TRAFFIC_SECRET " + td.ClientRandom + " " + td.ServerHandshakeTrafficSecret + "\n"
		}
		if td.ClientTrafficSecret0 != "" {
			res += "CLIENT_TRAFFIC_SECRET_0 " + td.ClientRandom + " " + td.ClientTrafficSecret0 + "\n"
		}
		if td.ServerTrafficSecret0 != "" {
			res += "SERVER_TRAFFIC_SECRET_0 " + td.ClientRandom + " " + td.ServerTrafficSecret0 + "\n"
		}
		if td.EarlyExporterSecret != "" {
			res += "EARLY_EXPORTER_SECRET " + td.ClientRandom + " " + td.EarlyExporterSecret + "\n"
		}
		if td.ExporterSecret != "" {
			res += "EXPORTER_SECRET " + td.ClientRandom + " " + td.ExporterSecret + "\n"
		}
		return res
	case "TLSv1.2", "TLSv1.1", "TLSv1.0":
		return "CLIENT_RANDOM " + td.ClientRandom + " " + td.SessionKey + "\n"
	default:
		return ""
	}
}

func NewTlsData(req []byte) (*TlsData, error) {
	var err error
	var data *TlsData = new(TlsData)

	if err = json.Unmarshal(req, data); err != nil {
		return nil, err
	}

	return data, nil
}
