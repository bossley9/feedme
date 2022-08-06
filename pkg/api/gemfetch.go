package api

import (
	"bufio"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/tofu"
)

var (
	hosts     tofu.KnownHosts
	hostsfile *tofu.HostWriter
	scanner   *bufio.Scanner
)

func xdgDataHome() string {
	if s, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
		return s
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "share")
}

func init() {
	// Load known hosts file
	path := filepath.Join(xdgDataHome(), "gemini", "known_hosts")
	err := hosts.Load(path)
	if err != nil {
		log.Fatal(err)
	}

	hostsfile, err = tofu.OpenHostsFile(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner = bufio.NewScanner(os.Stdin)
}

const trustPrompt = `The certificate offered by %s is of unknown trust. Its fingerprint is:
%s

If you knew the fingerprint to expect in advance, verify that this matches.
Otherwise, this should be safe to trust.

[t]rust always; trust [o]nce; [a]bort
=> `

func trustCertificate(hostname string, cert *x509.Certificate) error {
	host := tofu.NewHost(hostname, cert.Raw)
	knownHost, ok := hosts.Lookup(hostname)
	if ok {
		// Check fingerprint
		if knownHost.Fingerprint != host.Fingerprint {
			return errors.New("error: fingerprint does not match!")
		}
		return nil
	}

	hosts.Add(host)
	hostsfile.WriteHost(host)
	return nil
}

func getInput(prompt string) (input string, ok bool) {
	fmt.Printf("%s ", prompt)
	scanner.Scan()
	return scanner.Text(), true
}

func do(req *gemini.Request, via []*gemini.Request) (*gemini.Response, error) {
	client := gemini.Client{
		TrustCertificate: trustCertificate,
	}
	ctx := context.Background()
	resp, err := client.Do(ctx, req)
	if err != nil {
		return resp, err
	}

	switch resp.Status.Class() {
	case gemini.StatusInput:
		input, ok := getInput(resp.Meta)
		if !ok {
			break
		}
		req.URL.ForceQuery = true
		req.URL.RawQuery = gemini.QueryEscape(input)
		return do(req, via)

	case gemini.StatusRedirect:
		via = append(via, req)
		if len(via) > 5 {
			return resp, errors.New("too many redirects")
		}

		target, err := url.Parse(resp.Meta)
		if err != nil {
			return resp, err
		}
		target = req.URL.ResolveReference(target)
		redirect := *req
		redirect.URL = target
		return do(&redirect, via)
	}

	return resp, err
}

func FetchGemini(url string) ([]byte, error) {
	req, err := gemini.NewRequest(url)
	if err != nil {
		return []byte{}, err
	}
	res, err := do(req, nil)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	if res.Status.Class() != gemini.StatusSuccess {
		message := string(res.Status) + ": " + res.Meta
		return []byte{}, errors.New(message)
	}

	return ioutil.ReadAll(res.Body)
}
