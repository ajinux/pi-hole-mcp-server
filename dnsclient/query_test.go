package dnsclient

import (
	"fmt"
	"testing"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func Test_query(t *testing.T) {
	res, err := whois.Whois("likexian.com")
	if err != nil {
		t.Errorf("whois.Whois() error = %v", err)
		return
	}
	result, err := whoisparser.Parse(res)
	if err != nil {
		t.Errorf("whoisparser.Parse() error = %v", err)
		return
	}
	fmt.Printf("res = %+v", result)
}
