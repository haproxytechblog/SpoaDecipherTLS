package main

import (
	"errors"
	"github.com/negasus/haproxy-spoe-go/message"
	"strings"
)

const (
	ClientRnd                    string = "client-rnd"
	ProtocolVersion              string = "version"
	SessionKey                   string = "ssl-session-key"
	ClientEarlyTrafficSecret     string = "Client-Early-Traffic-Secret"
	ClientHandshakeTrafficSecret string = "Client-Handshake-Traffic-Secret"
	ServerHandshakeTrafficSecret string = "Server-Handshake-Traffic-Secret"
	ClientTrafficSecret0         string = "Client-Traffic-Secret-0"
	ServerTrafficSecret0         string = "Server-Traffic-Secret-0"
	EarlyExporterSecret          string = "Early-Exporter-Secret"
	ExporterSecret               string = "Exporter-Secret"
)

type TlsData struct {
	ProtocolVersion              string
	ClientRandom                 string
	SessionKey                   string
	ClientEarlyTrafficSecret     string
	ClientHandshakeTrafficSecret string
	ServerHandshakeTrafficSecret string
	ClientTrafficSecret0         string
	ServerTrafficSecret0         string
	EarlyExporterSecret          string
	ExporterSecret               string
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

func NewTlsData(msg *message.Message) (*TlsData, error) {
	var exists bool
	var data *TlsData = new(TlsData)
	var buffer any
	var fieldsMap map[string]*string = map[string]*string{
		SessionKey:                   &data.SessionKey,
		ClientEarlyTrafficSecret:     &data.ClientEarlyTrafficSecret,
		ClientHandshakeTrafficSecret: &data.ClientHandshakeTrafficSecret,
		ServerHandshakeTrafficSecret: &data.ServerHandshakeTrafficSecret,
		ClientTrafficSecret0:         &data.ClientTrafficSecret0,
		ServerTrafficSecret0:         &data.ServerTrafficSecret0,
		EarlyExporterSecret:          &data.EarlyExporterSecret,
		ExporterSecret:               &data.ExporterSecret,
	}

	if buffer, exists = msg.KV.Get(ProtocolVersion); !exists || buffer == nil {
		return nil, errors.New("tls version not found")
	}
	data.ProtocolVersion = buffer.(string)

	switch data.ProtocolVersion {
	case "TLSv1.3":
		if buffer, exists = msg.KV.Get(ClientRnd); !exists {
			return nil, errors.New("tls client-rnd not found")
		}
		data.ClientRandom = strings.ToLower(buffer.(string))

		for k, v := range fieldsMap {
			if buffer, exists = msg.KV.Get(k); exists && buffer != nil {
				*v = strings.ToLower(buffer.(string))
			}
		}

		return data, nil

	case "TLSv1.2", "TLSv1.1", "TLSv1.0":
		if buffer, exists = msg.KV.Get(ClientRnd); !exists {
			return nil, errors.New("tls client-rnd not found")
		}
		data.ClientRandom = strings.ToLower(buffer.(string))
		if buffer, exists = msg.KV.Get(SessionKey); !exists {
			return nil, errors.New("tls session-key not found")
		}
		data.SessionKey = strings.ToLower(buffer.(string))
		return data, nil
	default:
		return nil, errors.New("tls version not supported")
	}
}
