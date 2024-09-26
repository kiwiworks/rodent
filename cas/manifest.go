package cas

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/kiwiworks/rodent/errors"
	"github.com/kiwiworks/rodent/logger"
)

type Manifest struct {
	Endpoint         string `required:"true" json:"endpoint" yaml:"endpoint"`
	ClientId         string `required:"true" split_words:"true" json:"clientId" yaml:"clientId"`
	ClientSecret     string `required:"true" split_words:"true" json:"clientSecret" yaml:"clientSecret"`
	CertificatePath  string `required:"true" split_words:"true" json:"certificatePath" yaml:"certificatePath"`
	OrganizationName string `required:"true" split_words:"true" json:"organizationName" yaml:"organizationName"`
	ApplicationName  string `required:"true" split_words:"true" json:"applicationName" yaml:"applicationName"`
}

func ManifestFromFile(filename string) (*Manifest, error) {
	var manifest Manifest

	log := logger.New()

	ext := filepath.Ext(filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "opening manifest file '%s' failed", filename)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("failed to close manifest file", zap.String("filename", filename), zap.Error(err))
		}
	}(file)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.NewDecoder(file).Decode(&manifest); err != nil {
			return nil, errors.Wrapf(err, "unable to unmarshal manifest from yaml")
		}
		return &manifest, nil
	case ".json":
		if err := json.NewDecoder(file).Decode(&manifest); err != nil {
			return nil, errors.Wrapf(err, "unable to unmarshal manifest from yaml")
		}
		return &manifest, nil
	default:
		return nil, errors.Newf("unknown manifest extension '%s'", ext)
	}
}

func (m *Manifest) asCasdoorConfig(manifestPath string) (*CasdoorClientConfig, error) {
	manifestDir := filepath.Dir(manifestPath)
	certificatePath := path.Join(manifestDir, m.CertificatePath)
	certContent, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load x509 certificate content from '%s'", certificatePath)
	}

	decodedCert, _ := pem.Decode(certContent)
	_, err = x509.ParseCertificate(decodedCert.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load x509 certificate content from '%s'", certificatePath)
	}

	endpoint, err := url.Parse(m.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse endpoint '%s'", m.Endpoint)
	}

	return &CasdoorClientConfig{
		Endpoint:         *endpoint,
		ClientId:         m.ClientId,
		ClientSecret:     m.ClientSecret,
		Certificate:      string(certContent),
		OrganizationName: m.OrganizationName,
		ApplicationName:  m.ApplicationName,
	}, nil
}
