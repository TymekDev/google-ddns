package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var (
		interval           time.Duration
		domains            []string
		username, password string
	)
	pflag.StringVarP(&username, "username", "u", "", "")
	pflag.StringVarP(&password, "password", "p", "", "")
	pflag.StringSliceVarP(&domains, "domain", "d", nil, "")
	pflag.DurationVarP(&interval, "interval", "n", 5*time.Minute, "")
	pflag.Parse()

	if len(username) == 0 || len(password) == 0 || len(domains) == 0 {
		log.Fatalln("username, password, or domains are missing")
	}

	for range time.Tick(interval) {
		public, err := publicIP()
		if err != nil {
			log.Println("ERROR", err)
			continue
		}

		for _, domain := range domains {
			ip, err := lookup(domain)
			if err != nil {
				log.Println("ERROR", err)
				continue
			}

			if ip.String() == public.String() {
				continue
			}

			if err := update(username, password, domain); err != nil {
				log.Println("ERROR", err)
				continue
			}
		}
	}
}

func publicIP() (net.IP, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(b)), nil
}

func lookup(domain string) (net.IP, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 1000 * time.Millisecond}
			return d.Dial(network, address)
		},
	}

	addrs, err := r.LookupHost(context.Background(), domain)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(addrs[0]), nil
}

func update(username, password, domain string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s:%s@domains.google.com/nic/update", username, password),
		strings.NewReader(fmt.Sprint("hostname=", domain)),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Golang tymek.makowski@gmail.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return (err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println("INFO", domain, string(b))

	return nil
}
