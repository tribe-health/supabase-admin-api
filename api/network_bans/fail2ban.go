package network_bans

import (
	"log"

	"github.com/Sean-Der/fail2go"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type Fail2Ban struct {
	Fail2banSocket string
}

func (n *Fail2Ban) ListBannedIps(jails []string) ([]string, error) {
	conn := fail2go.Newfail2goConn(n.Fail2banSocket)
	filteredJails, err := n.filterJails(jails, conn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to filter jails")
	}
	allBans := make([]string, 0)
	for _, jail := range filteredJails {
		log.Printf("get %s", jail)
		_, _, _, _, _, ipList, err := conn.JailStatus(jail)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get jail %s status", jail)
		}
		allBans = append(allBans, ipList...)
	}
	return allBans, nil
}

func (n *Fail2Ban) filterJails(jails []string, conn *fail2go.Conn) ([]string, error) {
	filteredJails := make([]string, 0)
	globalStatus, err := conn.GlobalStatus()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get global status")
	}
	for _, jail := range jails {
		if slices.Index(globalStatus, jail) == -1 {
			log.Printf("ignoring jail %s\n", jail)
			continue
		}
		filteredJails = append(filteredJails, jail)
	}
	return filteredJails, nil
}

func (n *Fail2Ban) UnbanIp(jails []string, ipAddr string) error {
	conn := fail2go.Newfail2goConn(n.Fail2banSocket)
	filteredJails, err := n.filterJails(jails, conn)
	if err != nil {
		return errors.Wrap(err, "failed to filter jails")
	}
	for _, jail := range filteredJails {
		log.Printf("unban %s from %s\n", ipAddr, jail)
		_, err := conn.JailUnbanIP(jail, ipAddr)
		if err != nil {
			return errors.Wrapf(err, "failed to unban for %s", jail)
		}
	}
	return nil
}
