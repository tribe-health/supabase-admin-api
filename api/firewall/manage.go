package firewall

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

//go:embed static/rules.conf
var res embed.FS

type Config struct {
	PrivilegedPorts []int `yaml:"privileged_ports" required:"true"`
	CustomerPorts   []int `yaml:"customer_ports" required:"true"`

	PrivilegedPortsWhitelist []string `yaml:"privileged_ports_whitelist" required:"true"`
}

func parseSubnets(networks []string) ([]net.IPNet, error) {
	out := make([]net.IPNet, 0)
	if len(networks) == 0 {
		return out, fmt.Errorf("at least one subnet is required")
	}
	for _, subnet := range networks {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			return out, errors.Wrapf(err, "failed to parse %s", subnet)
		}
		out = append(out, *ipNet)
	}
	return out, nil
}

func (c *Config) CreateFirewallManager() (*Manager, error) {
	if len(c.PrivilegedPorts) == 0 || len(c.CustomerPorts) == 0 {
		return nil, fmt.Errorf("at least one port must be specified for each category")
	}
	subnets, err := parseSubnets(c.PrivilegedPortsWhitelist)
	if err != nil {
		return nil, err
	}
	return &Manager{
		PrivilegedPorts:          c.PrivilegedPorts,
		CustomerPorts:            c.CustomerPorts,
		PrivilegedPortsWhitelist: subnets,
	}, nil
}

type Manager struct {
	PrivilegedPorts []int
	CustomerPorts   []int

	PrivilegedPortsWhitelist []net.IPNet
}

type FirewallRequest struct {
	UserWhitelist []string `json:"user_whitelist" required:"true"`
}

func (f *Manager) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	var req FirewallRequest

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&req); err != nil {
		return errors.Wrap(err, "invalid request")
	}
	subnets, err := parseSubnets(req.UserWhitelist)
	if err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	processedRules, err := f.Process(subnets)
	if err != nil {
		return errors.Wrap(err, "failed to process firewall request")
	}
	err = f.Apply(processedRules)
	if err != nil {
		return errors.Wrap(err, "failed to apply requested rules")
	}
	return nil
}

func (f *Manager) Process(userNetworks []net.IPNet) (string, error) {
	tpl, err := template.ParseFS(res, "static/rules.conf")
	if err != nil {
		return "", errors.Wrap(err, "failed to read template")
	}
	var output bytes.Buffer
	data := struct {
		UserWhitelist            string
		PrivilegedPorts          string
		CustomerPorts            string
		PrivilegedPortsWhitelist string
	}{
		UserWhitelist:            f.getWhitelist(userNetworks),
		PrivilegedPorts:          f.getPorts(f.PrivilegedPorts),
		CustomerPorts:            f.getPorts(f.CustomerPorts),
		PrivilegedPortsWhitelist: f.getWhitelist(f.PrivilegedPortsWhitelist),
	}
	err = tpl.Execute(&output, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to render rules")
	}
	return output.String(), nil
}

func (f *Manager) Apply(rules string) error {
	file, err := ioutil.TempFile("", "nft_rules")
	if err != nil {
		return errors.Wrap(err, "failed to create temp file for fw rules")
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(rules)
	if err != nil {
		return errors.Wrap(err, "failed to write new rules to temp file")
	}
	cmd := exec.Command("sudo", "nft", "-f", file.Name())
	stdout, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "failed to apply new rules")
	}
	logrus.WithField("nftOutput", stdout).Info("applied new firewall rules")
	return nil
}

func (f *Manager) getPorts(ports []int) string {
	strPorts := make([]string, 0)
	for _, port := range ports {
		strPorts = append(strPorts, strconv.Itoa(port))
	}
	return strings.Join(strPorts, ", ")
}

func (f *Manager) getWhitelist(nets []net.IPNet) string {
	allowedSubnets := make([]string, 0)
	for _, net := range nets {
		allowedSubnets = append(allowedSubnets, net.String())
	}
	return strings.Join(allowedSubnets, ", ")
}
