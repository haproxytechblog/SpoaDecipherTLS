package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
)

// Arguments is the structure used to parse the parameters passed when the program starts
type Arguments struct {
	BindAddress   string `default:"0.0.0.0" usage:"IP address to bind on" alias:"a"`
	BindPort      int    `default:"12345" usage:"TCP port to bind on" alias:"p"`
	LogLevel      uint8  `default:"5" usage:"Set the log level (0 (no log) to 9 (trace))" alias:"l"`
	Version       bool   `default:"false" usage:"Print the version" alias:"V"`
	EnableTls     bool   `default:"false" usage:"Enable the TLS layer" alias:"tls"`
	EnableMTls    bool   `default:"false" usage:"Enable the mTLS layer" alias:"mtls"`
	CaCert        string `default:"./ca.pem" usage:"CA certificate path" alias:"ca"`
	SpoaCert      string `default:"./spoa.crt" usage:"Agent certificate path" alias:"cert"`
	KeyType       string `default:"RSA" usage:"Key type to use (RSA or ECDSA)" alias:"kt"`
	KeySize       uint16 `default:"2048" usage:"Key size (only for RSA key type)" alias:"ks"`
	KeyCurve      string `default:"P384" usage:"Key curve (P224, P256, P384, or P521) (Only for ECDSA key type)" alias:"kc"`
	GenCa         bool   `default:"false" usage:"Generate a CA certificate" alias:"genca"`
	GenSpoeCert   bool   `default:"false" usage:"Generate a certificate for haproxy" alias:"genspoe"`
	GenSpoaCert   bool   `default:"false" usage:"Generate a certificate for the agent" alias:"genspoa"`
	CertOut       string `default:"./cert.pem" usage:"Where the output certificate will be saved" alias:"out"`
	Cn            string `default:"" usage:"Common Name to use when creating a certificate" alias:"cn"`
	NssKeylogFile string `default:"./nsskeylogfile" usage:"The NSS Keylog file path" alias:"f"`
	TlsMinVersion string `default:"TLSv1.2" usage:"Minimum TLS version (supported values TLSv1.0, TLSv1.1, TLSv1.2, TLSv1.3)" alias:"tlsmin"`
	TlsMaxVersion string `default:"TLSv1.3" usage:"Maximum TLS version (supported values TLSv1.0, TLSv1.1, TLSv1.2, TLSv1.3)" alias:"tlsmax"`
}

func (a *Arguments) LogOptions(logger *log.Logger) {
	logger.Println("Version is " + version)
	logger.Println("Compiled on " + compileDate)
	logger.Println("Commit " + commit)
	logger.Println("Bind IP is " + a.BindAddress)
	logger.Println("Bind port is " + strconv.Itoa(a.BindPort))
	logger.Printf("TLS is %t\n", a.EnableTls)
	if a.EnableTls {
		logger.Println("CA certificate path is " + a.CaCert)
		logger.Println("SPOA certificate path is " + a.SpoaCert)
	}
	logger.Printf("mTLS is %t\n", a.EnableMTls)
	logger.Println("NSS Keylog file paht is " + a.NssKeylogFile)
}

func (a *Arguments) GetBindAddressAndPort() string {
	return fmt.Sprintf("%s:%d", a.BindAddress, a.BindPort)
}

func (a *Arguments) GetCaCert() string    { return a.CaCert }
func (a *Arguments) GetKeyType() string   { return a.KeyType }
func (a *Arguments) GetKeySize() uint16   { return a.KeySize }
func (a *Arguments) GetKeyCurve() string  { return a.KeyCurve }
func (a *Arguments) GetCertOut() string   { return a.CertOut }
func (a *Arguments) GetCn() string        { return a.Cn }
func (a *Arguments) GetSpoaCert() string  { return a.SpoaCert }
func (a *Arguments) GetGenCa() bool       { return a.GenCa }
func (a *Arguments) GetGenSpoeCert() bool { return a.GenSpoeCert }
func (a *Arguments) GetGenSpoaCert() bool { return a.GenSpoaCert }
func (a *Arguments) GetMTls() bool        { return a.EnableMTls }
func (a *Arguments) GetTlsMinVersion() uint16 {
	switch a.TlsMinVersion {
	case "TLSv1.0":
		return tls.VersionTLS10
	case "TLSv1.1":
		return tls.VersionTLS11
	case "TLSv1.2":
		return tls.VersionTLS12
	case "TLSv1.3":
		return tls.VersionTLS13
	default:
		panic(fmt.Sprintf("%s is not a supported TLS version", a.TlsMinVersion))
	}
}
func (a *Arguments) GetTlsMaxVersion() uint16 {
	switch a.TlsMaxVersion {
	case "TLSv1.0":
		return tls.VersionTLS10
	case "TLSv1.1":
		return tls.VersionTLS11
	case "TLSv1.2":
		return tls.VersionTLS12
	case "TLSv1.3":
		return tls.VersionTLS13
	default:
		panic(fmt.Sprintf("%s is not a supported TLS version", a.TlsMinVersion))
	}
}
