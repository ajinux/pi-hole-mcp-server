package domain

import (
	"fmt"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func Whois(domain string) (*whoisparser.WhoisInfo, error) {
	res, err := whois.Whois(domain)
	if err != nil {
		return nil, fmt.Errorf("whois.Whois() error = %v", err)
	}
	result, err := whoisparser.Parse(res)
	if err != nil {
		return nil, fmt.Errorf("whoisparser.Parse() error = %v", err)
	}
	return &result, nil
}
