package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fediverse/accountaddress"
	"fediverse/acct"
	activityclient "fediverse/application/activity/client"
	"fediverse/application/activity/routes"
	"fediverse/application/common"
	"fediverse/application/config"
	"fediverse/application/following"
	"fediverse/application/keymanager"
	"fediverse/application/posts"
	"fediverse/application/schema"
	"fediverse/security/rsahelpers"
	"fediverse/security/rsassapkcsv115sha256"
	"fediverse/webfinger"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piprate/json-gold/ld"
)

func readFromStdin() []byte {
	var content []byte
	go func() {
		<-time.After(100 * time.Millisecond)
		if len(content) == 0 {
			fmt.Fprintf(os.Stderr, "Please provide some content\n")
		}
	}()
	for {
		buffer := make([]byte, 1024)
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			panic(err)
		}
		content = append(content, buffer[:n]...)
		if n < 1024 {
			break
		}
	}
	return content
}

func parsePrivateKey(payload []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(payload))
	if block == nil || block.Type != rsahelpers.RSAPrivateKeyLabel {
		return nil, fmt.Errorf("invalid private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}
	return privateKey, nil
}

func parsePublicKey(payload []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(payload))
	if block == nil || block.Type != rsahelpers.PublicKeyLabel {
		return nil, fmt.Errorf("invalid public key")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}
	rsaPubKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to RSA public key")
	}
	return rsaPubKey, nil
}

func main() {
	err := schema.Initialize()

	if err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Fprintf(os.Stderr, "Please provide a command\n")
		os.Exit(1)
		return
	}

	command := args[0]

	switch command {
	case "create":
		// This command creates a new post. This command expects a content as the
		// first argument.

		// TODO: accept content from stdin, and check if the content is UTF-8.

		if len(args) <= 1 {
			fmt.Fprintf(os.Stderr, "Please provide some content\n")
			os.Exit(1)
			return
		}
		err := posts.CreatePost(args[1])
		if err != nil {
			panic(err)
		}
		fmt.Println("Post successfully created!")
	case "follow":
		if len(args) <= 1 {
			fmt.Fprintf(os.Stderr, "Please provide a account address to follow\n")
			os.Exit(1)
			return
		}

		address, err := accountaddress.ParseAccountAddress(args[1])
		if err != nil && errors.Is(err, accountaddress.ErrInvalidAccountAddress()) {
			fmt.Fprintf(os.Stderr, "Invalid account address\n")
			os.Exit(1)
			return
		}

		// Perform a WebFinger lookup.
		fmt.Printf("Performing WebFinger lookup for %s...\n", acct.Acct(address).String())
		j, err := webfinger.Lookup(address.Host, acct.Acct(address).String(), []string{"self"})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error looking up account: %s\n", err.Error())
			os.Exit(2)
			return
		}

		links, ok := j.Links.Value()
		if !ok {
			fmt.Fprint(os.Stderr, "No properties found\n")
			os.Exit(1)
			return
		}

		var selfLink string
		ok = false
		for _, link := range links {
			if link.Rel == "self" {
				selfLink = link.Href
				ok = true
				break
			}
		}
		if !ok {
			fmt.Fprint(os.Stderr, "self link not found in WebFinger lookup\n")
			os.Exit(1)
		}
		if selfLink == "" {
			fmt.Fprint(os.Stderr, "self link is empty\n")
			os.Exit(1)
		}

		fmt.Println("Got self link:", selfLink)

		req, err := http.NewRequest("GET", selfLink, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating request: %s\n", err.Error())
			os.Exit(1)
		}
		req.Header.Set("Accept", "application/activity+json")

		fmt.Println("Looking up the ActivityPub actor using the self link")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Response failed %s\n", err.Error())
			os.Exit(1)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read response body: %s\n", err.Error())
			os.Exit(1)
		}
		var parsed any
		err = json.Unmarshal(body, &parsed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse response body: %s\n", err.Error())
			os.Exit(1)
		}

		proc := ld.NewJsonLdProcessor()
		options := ld.NewJsonLdOptions("")

		expanded, err := proc.Expand(parsed, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to expand JSON-LD document: %s\n", err.Error())
		}
		var inboxID string
		gotInbox := false
		for _, node := range expanded {
			predicateObjectMap, ok := node.(map[string]any)
			if !ok {
				continue
			}
			inbox, ok := predicateObjectMap["http://www.w3.org/ns/ldp#inbox"].([]any)
			if !ok {
				continue
			}
			value, ok := inbox[0].(map[string]any)
			if !ok {
				continue
			}
			inboxID, ok = value["@id"].(string)
			if !ok {
				continue
			}
			gotInbox = true
		}
		if !gotInbox {
			fmt.Fprintf(os.Stderr, "No inbox found\n")
			os.Exit(1)
		}

		fmt.Println("The URL to the inbox:", inboxID)

		id, err := following.AddFollowing(address)
		if err != nil {
			// TODO: this also fails if the user is already following the account.
			//   just silently ignore the error, and return
			panic(err)
		}

		privateKey := keymanager.GetPrivateKey()
		signingKeyIRI := common.Origin() + routes.Activity{}.Actors().Actor().Route().FullRoute(config.Username()) + "#main-key"
		followActivityIRI := common.Origin() + routes.Activity{}.Actors().Actor().Following().FullRoute(config.Username()).Route().FullRoute() + strconv.FormatInt(id, 10)
		senderIRI := common.Origin() + routes.Activity{}.Actors().Actor().Route().FullRoute(config.Username())
		recipientID := selfLink
		inboxURL := inboxID

		fmt.Println("Sending follow activity...")
		err = activityclient.Follow(
			privateKey,
			activityclient.SigningKeyIRI(signingKeyIRI),
			activityclient.FollowActivityIRI(followActivityIRI),
			activityclient.SenderIRI(senderIRI),
			activityclient.ObjectIRI(recipientID),
			activityclient.InboxURL(inboxURL),
		)

		if err != nil {
			panic(err)
		}
	case "genrsa":
		// This command generates a new RSA key pair. It accepts a `--public` flag
		// to show the public key.

		// But the typical use case for this would be to generate only the private
		// key, and then use the `getrsapublic` command to get the public key.

		fs := flag.NewFlagSet("genrsa", flag.ExitOnError)
		var showPublic bool
		fs.BoolVar(&showPublic, "public", false, "Show public key")
		fs.Parse(os.Args[2:])

		privateKey, err := rsahelpers.GenerateRSPrivateKey(2048)
		if err != nil {
			panic(err)
		}

		fmt.Print(rsahelpers.PrivateKeyToPKCS1PEMString(privateKey))
		if showPublic {
			publicKeyPEMString, err := rsahelpers.PublicKeyToPKIXString(
				&privateKey.PublicKey,
			)
			if err != nil {
				panic(err)
			}
			fmt.Print(publicKeyPEMString)
		}
	case "deriversapublic":
		// This command gets the public key from a private key. It expects a private
		// key as the first argument, and will output the public key, in the
		// standard out (console, on your terminal).

		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Please provide a private key\n")
			os.Exit(1)
			return
		}
		content := os.Args[2]
		privateKey, err := parsePrivateKey([]byte(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing private key: %s\n", err.Error())
			os.Exit(1)
			return
		}
		publicKeyPEM, err := rsahelpers.PublicKeyToPKIXString(
			&privateKey.PublicKey,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling public key: %s\n", err)
			return
		}
		fmt.Print(string(publicKeyPEM))
	case "sign":
		// This command signs a payload. It expects a private key as the first
		// argumemt, and the payload coming in via Standard Input (stdin), that you
		// can either type into the console, or pipe in from another application.

		fs := flag.NewFlagSet("sign", flag.ExitOnError)

		var content string
		fs.StringVar(&content, "content", "", "Content to sign")
		fs.Parse(os.Args[4:])

		var payload []byte
		if content == "" {
			payload = []byte(readFromStdin())
		} else {
			payload = []byte(content)
		}

		pemPrivateKey := os.Args[2]
		if pemPrivateKey == "--" {
			pemPrivateKey = os.Args[3]
		}

		privateKey, err := parsePrivateKey([]byte(pemPrivateKey))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid private key\n")
			os.Exit(1)
			return
		}

		sig, err := rsassapkcsv115sha256.Base64Signer(privateKey).Sign(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error signing payload: %s\n", err)
			os.Exit(1)
		}

		fmt.Print(sig)
	case "verify":
		fs := flag.NewFlagSet("verify", flag.ExitOnError)

		var publicKey string
		var signatureBase64 string
		var content string
		fs.StringVar(&publicKey, "publicKey", "", "A PEM-encoded public key")
		fs.StringVar(&signatureBase64, "signature", "", "A base64-encoded signature")
		fs.StringVar(&content, "content", "", "Optional content to verify. If not set, then the content will be read from stdin.")
		fs.Parse(os.Args[2:])

		if publicKey == "" {
			fmt.Fprintf(os.Stderr, "Please provide a public key\n")
			fs.Usage()
			os.Exit(1)
			return
		}
		if signatureBase64 == "" {
			fmt.Fprintf(os.Stderr, "Please provide a signature\n")
			fs.Usage()
			os.Exit(1)
			return
		}

		var payload []byte
		if content == "" {
			payload = readFromStdin()
		} else {
			payload = []byte(content)
		}

		rsaPublicKey, err := parsePublicKey([]byte(publicKey))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing public key: %s\n", err.Error())
			os.Exit(1)
			return
		}

		if err := rsassapkcsv115sha256.Base64Verifier(rsaPublicKey).Verify(payload, signatureBase64); err == nil {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature verification failed:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s.", command)
	}
}
