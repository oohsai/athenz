//
// Copyright The Athenz Authors
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
//

package sia

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	gcpa "github.com/AthenZ/athenz/libs/go/sia/gcp/attestation"
	"github.com/AthenZ/athenz/libs/go/sia/gcp/meta"
	"github.com/AthenZ/athenz/libs/go/sia/host/ip"
	"github.com/AthenZ/athenz/libs/go/sia/host/signature"
	"github.com/AthenZ/athenz/libs/go/sia/host/utils"
	"net"
	"net/url"
)

type GCEProvider struct {
	Name string
}

// GetName returns the name of the current provider
func (gke GCEProvider) GetName() string {
	return gke.Name
}

// GetHostname returns the hostname as per the provider
func (gke GCEProvider) GetHostname(fqdn bool) string {
	return utils.GetHostname(fqdn)
}

func (gke GCEProvider) AttestationData(svc string, key crypto.PrivateKey, sigInfo *signature.SignatureInfo) (string, error) {
	result, err := meta.GetData("http://169.254.169.254", "/computeMetadata/v1/instance/service-accounts/default/identity?audience=https://zts.athenz.io&format=full")
	if err == nil {
		return string(result), nil
	}
	return "", fmt.Errorf("error while retriveing attestation data")
}

func (gke GCEProvider) PrepareKey(file string) (crypto.PrivateKey, error) {
	return "", fmt.Errorf("not implemented")
}

func (gke GCEProvider) GetCsrDn() pkix.Name {
	return pkix.Name{}
}

func (gke GCEProvider) GetSanDns(service string, includeHost bool, wildcard bool, cnames []string) []string {
	return nil
}

func (gke GCEProvider) GetSanUri(svc string, opts ip.Opts, spiffeTrustDomain, spiffeNamespace string) []*url.URL {
	return nil
}

func (gke GCEProvider) GetEmail(service string) []string {
	return nil
}

func (gke GCEProvider) GetRoleDnsNames(cert *x509.Certificate, service string) []string {
	return nil
}

func (gke GCEProvider) GetSanIp(docIp map[string]bool, ips []net.IP, opts ip.Opts) []net.IP {
	return nil
}

func (gke GCEProvider) GetSuffix() string {
	return ""
}

func (gke GCEProvider) CloudAttestationData(base, svc, ztsServerName string) (string, error) {
	return gcpa.New(base, svc, ztsServerName)
}

func (gke GCEProvider) GetAccountDomainServiceFromMeta(base string) (string, string, string, error) {
	account, err := meta.GetProject(base)
	if err != nil {
		return "", "", "", err
	}
	domain, err := meta.GetDomain(base)
	if err != nil {
		return account, "", "", err
	}
	service, err := meta.GetService(base)
	if err != nil {
		return account, domain, "", err
	}
	return account, domain, service, nil
}

func (tp GCEProvider) GetAccessManagementProfileFromMeta(base string) (string, error) {
	profile, err := meta.GetProfile(base)
	if err != nil {
		return "", err
	}
	return profile, nil
}

func (tp GCEProvider) GetAdditionalSshHostPrincipals(base string) (string, error) {
	instanceName, err := meta.GetInstanceName(base)
	if err != nil {
		return "", err
	}
	project, err := meta.GetProject(base)
	if err != nil {
		return fmt.Sprintf("%s", instanceName), nil
	}
	instanceId, err := meta.GetInstanceId(base)
	if err != nil {
		return fmt.Sprintf("%s,%s.c.%s.internal", instanceName, instanceName, project), nil
	}
	return fmt.Sprintf("%s,compute.%s,%s.c.%s.internal", instanceName, instanceId, instanceName, project), nil
}
