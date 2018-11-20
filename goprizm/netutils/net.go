package netutils

import (
	"fmt"
	"goprizm/log"
	"net"
	"strings"
	"time"
)

// NormalizeMac transforms a MAC to all lower case with delimitters removed.
func NormalizeMAC(str string) (string, error) {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, "-", "", -1)
	str = strings.Replace(str, ".", "", -1)
	str = strings.Replace(str, ":", "", -1)
	str = strings.ToLower(str)

	for _, c := range str {
		if !((c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9')) {
			return str, fmt.Errorf("normalize - err: invalid characters(%s)", str)
		}
	}

	return str, nil
}

//  IsIPAddr return true if given string is a valid IPv4/IPv6 IP address.
func IsIPAddr(s string) bool {
	return net.ParseIP(s) != nil
}

// Return true if fqdn is actually a FQDN of h.
//  isFQDN("abc.domain.name", "abc") => true
//  isFQDN("abc.domain.name", "") => true
//  isFQDN("abc.domain.name", "xyz") => false
//  isFQDN("abc.domain.name", "abc.xyz.com") => false
func IsFQDN(fqdn, hostname string) bool {
	return len(fqdn) > len(hostname) && strings.HasPrefix(fqdn, hostname)
}

// RetryOp retries given operation atmost n times if it returns error. Wait for after
// duration before every retry. If n < 0, retry is performed till op succeeds.
func RetryOp(n int, after time.Duration, opName string, op func() error) (err error) {
	i := 0
	for {
		// If n is +ve and num of retries reached n.
		if n > 0 && i == n {
			return err
		}

		// Execute op and on error sleep for given duration.
		if err = op(); err != nil {
			i += 1
			log.Errorf("%s(retry:%d) failed err:%v", opName, i, err)
			time.Sleep(after)
			continue
		}
		return nil
	}
}
