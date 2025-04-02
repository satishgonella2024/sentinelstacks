package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// SignatureInfo contains metadata about a signature
type SignatureInfo struct {
	Signer    string    `json:"signer"`
	Timestamp time.Time `json:"timestamp"`
	Algorithm string    `json:"algorithm"`
	Version   string    `json:"version"`
	KeyID     string    `json:"keyId"`
}

// Signature represents a cryptographic signature with metadata
type Signature struct {
	Data    string       `json:"data"`
	Info    SignatureInfo `json:"info"`
	Content []byte       `json:"-"` // Original content (not serialized)
}

// KeyManager handles cryptographic keys for signing and verification
type KeyManager struct {
	KeysDir string
}

// NewKeyManager creates a new key manager
func NewKeyManager(keysDir string) (*KeyManager, error) {
	// Ensure keys directory exists
	if err := os.MkdirAll(keysDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	return &KeyManager{
		KeysDir: keysDir,
	}, nil
}

// GenerateKeyPair generates a new RSA key pair
func (km *KeyManager) GenerateKeyPair(keyID string, keySize int) error {
	if keySize == 0 {
		keySize = 2048 // Default key size
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Encode private key to PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key to PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Write keys to files
	privateKeyPath := filepath.Join(km.KeysDir, keyID+".key")
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	publicKeyPath := filepath.Join(km.KeysDir, keyID+".pub")
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

// Sign creates a signature for the given content
func (km *KeyManager) Sign(content []byte, keyID string, signer string) (*Signature, error) {
	// Load private key
	privateKey, err := km.loadPrivateKey(keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Create hash of content
	hash := sha256.Sum256(content)

	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign content: %w", err)
	}

	// Encode signature to base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Create signature info
	info := SignatureInfo{
		Signer:    signer,
		Timestamp: time.Now().UTC(),
		Algorithm: "RSA-PKCS1-SHA256",
		Version:   "1.0",
		KeyID:     keyID,
	}

	return &Signature{
		Data:    signatureBase64,
		Info:    info,
		Content: content,
	}, nil
}

// SignFile creates a signature for a file
func (km *KeyManager) SignFile(filePath string, keyID string, signer string) (*Signature, error) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Sign content
	return km.Sign(content, keyID, signer)
}

// Verify checks if a signature is valid
func (km *KeyManager) Verify(signature *Signature, content []byte) error {
	// Load public key
	publicKey, err := km.loadPublicKey(signature.Info.KeyID)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}

	// Decode signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(signature.Data)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Create hash of content
	hash := sha256.Sum256(content)

	// Verify signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	return nil
}

// VerifyFile verifies a signature against a file
func (km *KeyManager) VerifyFile(signature *Signature, filePath string) error {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Verify content
	return km.Verify(signature, content)
}

// ImportPublicKey imports a public key from a file
func (km *KeyManager) ImportPublicKey(keyID string, keyPath string) error {
	// Read key file
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	// Verify it's a valid public key
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return fmt.Errorf("invalid public key format")
	}

	// Save to keys directory
	publicKeyPath := filepath.Join(km.KeysDir, keyID+".pub")
	if err := os.WriteFile(publicKeyPath, keyData, 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

// loadPrivateKey loads a private key from the keys directory
func (km *KeyManager) loadPrivateKey(keyID string) (*rsa.PrivateKey, error) {
	// Read key file
	keyPath := filepath.Join(km.KeysDir, keyID+".key")
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key format")
	}

	// Parse private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// loadPublicKey loads a public key from the keys directory
func (km *KeyManager) loadPublicKey(keyID string) (*rsa.PublicKey, error) {
	// Read key file
	keyPath := filepath.Join(km.KeysDir, keyID+".pub")
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key format")
	}

	// Parse public key
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Cast to RSA public key
	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}

// SignStream signs a data stream
func (km *KeyManager) SignStream(reader io.Reader, keyID string, signer string) (*Signature, error) {
	// Hash the stream
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return nil, fmt.Errorf("failed to hash stream: %w", err)
	}
	hashSum := hash.Sum(nil)

	// Load private key
	privateKey, err := km.loadPrivateKey(keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashSum)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %w", err)
	}

	// Encode signature to base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Create signature info
	info := SignatureInfo{
		Signer:    signer,
		Timestamp: time.Now().UTC(),
		Algorithm: "RSA-PKCS1-SHA256",
		Version:   "1.0",
		KeyID:     keyID,
	}

	return &Signature{
		Data:    signatureBase64,
		Info:    info,
		Content: hashSum,
	}, nil
}

// VerifyStream verifies a signature against a data stream
func (km *KeyManager) VerifyStream(reader io.Reader, signature *Signature) error {
	// Hash the stream
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return fmt.Errorf("failed to hash stream: %w", err)
	}
	hashSum := hash.Sum(nil)

	// Load public key
	publicKey, err := km.loadPublicKey(signature.Info.KeyID)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}

	// Decode signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(signature.Data)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Verify signature
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashSum, signatureBytes)
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	return nil
}
