package dnsclient

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// ExtractTopLevelDomain strips subdomain prefixes and returns the top-level domain
// Example: api.example.com -> example.com, www.subdomain.example.com -> example.com
func ExtractTopLevelDomain(domain string) string {
	// Remove trailing dot if present
	domain = strings.TrimSuffix(domain, ".")

	// Split by dots
	parts := strings.Split(domain, ".")

	// If already a TLD or invalid, return as-is
	if len(parts) <= 2 {
		return domain
	}

	// Return last two parts (domain.tld)
	return strings.Join(parts[len(parts)-2:], ".")
}

func query(domain string, qtype uint16) ([]dns.RR, error) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), qtype)

	// Use a public resolver. Can use 1.1.1.1, 8.8.8.8, or your own.
	resp, _, err := c.Exchange(&m, "1.1.1.1:53")
	if err != nil {
		return nil, err
	}
	if resp.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("bad rcode: %d", resp.Rcode)
	}
	return resp.Answer, nil
}

func fetchA(domain string) []string {
	records, _ := query(domain, dns.TypeA)
	var ips []string
	for _, rr := range records {
		if a, ok := rr.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips
}

func fetchAAAA(domain string) []string {
	records, _ := query(domain, dns.TypeAAAA)
	var ips []string
	for _, rr := range records {
		if a, ok := rr.(*dns.AAAA); ok {
			ips = append(ips, a.AAAA.String())
		}
	}
	return ips
}

func fetchNS(domain string) []string {
	records, _ := query(domain, dns.TypeNS)
	var nss []string
	for _, rr := range records {
		if ns, ok := rr.(*dns.NS); ok {
			nss = append(nss, ns.Ns)
		}
	}
	return nss
}

func fetchMX(domain string) []string {
	records, _ := query(domain, dns.TypeMX)
	var mxs []string
	for _, rr := range records {
		if mx, ok := rr.(*dns.MX); ok {
			mxs = append(mxs, fmt.Sprintf("%s (pref %d)", mx.Mx, mx.Preference))
		}
	}
	return mxs
}

func fetchTXT(domain string) []string {
	records, _ := query(domain, dns.TypeTXT)
	var txts []string
	for _, rr := range records {
		if txt, ok := rr.(*dns.TXT); ok {
			txts = append(txts, strings.Join(txt.Txt, " "))
		}
	}
	return txts
}

type DNSRecord struct {
	Domain string
	A      []string
	AAAA   []string
	NS     []string
	MX     []string
	TXT    []string
}

func GetAllRecords(domain string) (DNSRecord, error) {
	// Strip to top-level domain
	tld := ExtractTopLevelDomain(domain)

	// Query all record types
	record := DNSRecord{
		Domain: tld,
		A:      fetchA(tld),
		AAAA:   fetchAAAA(tld),
		NS:     fetchNS(tld),
		MX:     fetchMX(tld),
		TXT:    fetchTXT(tld),
	}

	// Check if domain exists (at least one record type should have data)
	if len(record.A) == 0 && len(record.AAAA) == 0 && len(record.NS) == 0 && len(record.MX) == 0 && len(record.TXT) == 0 {
		return record, fmt.Errorf("no such domain exists: %s", tld)
	}

	return record, nil
}
