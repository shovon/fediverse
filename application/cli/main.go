package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fediverse/application/posts"
	"fediverse/application/schema"
	"fediverse/cryptohelpers/rsahelpers"
	"flag"
	"fmt"
	"os"
	"time"
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
		fmt.Println("Please provide a command")
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
			fmt.Println("Please provide some content")
			os.Exit(1)
			return
		}
		posts.CreatePost(args[1])
		fmt.Println("Post successfully created!")
	case "genrsa":
		// This command generates a new RSA key pair. It accepts a `--public` flag
		// to show the public key.
		//
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
			fmt.Println("Please provide a private key")
			os.Exit(1)
			return
		}
		content := os.Args[2]
		privateKey, err := parsePrivateKey([]byte(content))
		if err != nil {
			fmt.Println("Error parsing private key:", err)
			os.Exit(1)
			return
		}
		publicKeyPEM, err := rsahelpers.PublicKeyToPKIXString(
			&privateKey.PublicKey,
		)
		if err != nil {
			fmt.Println("Error marshaling public key:", err)
			return
		}
		fmt.Print(string(publicKeyPEM))
	case "sign":
		// This command signs a payload. It expects a private key as the first
		// argumemt, and the payload coming in via Standard Input (stdin), that you
		// can either type into the console, or pipe in from another application.

		if len(os.Args) < 3 {
			fmt.Println("Please provide a private key")
			os.Exit(1)
			return
		}

		fs := flag.NewFlagSet("sign", flag.ExitOnError)

		var content string
		fs.StringVar(&content, "content", "", "Content to sign")
		fs.Parse(os.Args[2:])

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
			fmt.Println("Invalid private key")
			os.Exit(1)
			return
		}

		hash := sha256.Sum256(payload)
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
		if err != nil {
			fmt.Println("Error signing payload:", err)
			os.Exit(1)
		}

		fmt.Print(string(base64.StdEncoding.EncodeToString(signature)))
	case "verify":
		// I need the public key, the signature, and the payload

		fs := flag.NewFlagSet("verify", flag.ExitOnError)

		var publicKey string
		var signatureBase64 string
		var content string
		fs.StringVar(&publicKey, "publicKey", "", "A PEM-encoded public key")
		fs.StringVar(&signatureBase64, "signature", "", "A base64-encoded signature")
		fs.StringVar(&content, "content", "", "Optional content to verify. If not set, then the content will be read from stdin.")
		fs.Parse(os.Args[2:])

		if publicKey == "" {
			fmt.Println("Please provide a public key")
			fs.Usage()
			os.Exit(1)
			return
		}
		if signatureBase64 == "" {
			fmt.Println("Please provide a signature")
			fs.Usage()
			os.Exit(1)
			return
		}

		var payload []byte
		if content == "" {
			payload = []byte(readFromStdin())
		} else {
			payload = []byte(content)
		}

		rsaPublicKey, err := parsePublicKey([]byte(publicKey))
		if err != nil {
			fmt.Println("Error parsing public key:", err)
			os.Exit(1)
			return
		}

		hash := sha256.Sum256(payload)

		signature, err := base64.StdEncoding.DecodeString(signatureBase64)
		if err != nil {
			fmt.Println("Error decoding signature:", err)
			os.Exit(1)
			return
		}

		err = rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signature[:])
		if err == nil {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature verification failed:", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command %s.", command)
	}
}
