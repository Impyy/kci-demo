package profile

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Impyy/kci-demo/cmd/kci-demo/crypto"
)

const (
	profilePath = "profile.json"
)

type profileFormat struct {
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
}

type Profile struct {
	PublicKey *[crypto.KeySize]byte
	SecretKey *[crypto.KeySize]byte
}

func Load() (*Profile, error) {
	bytes, err := ioutil.ReadFile(profilePath)
	if err != nil {
		return nil, err
	}

	temp := profileFormat{}
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return nil, err
	}

	publicKey, err := hex.DecodeString(temp.PublicKey)
	if err != nil {
		return nil, err
	}
	secretKey, err := hex.DecodeString(temp.SecretKey)
	if err != nil {
		return nil, err
	}

	p := &Profile{
		PublicKey: new([crypto.KeySize]byte),
		SecretKey: new([crypto.KeySize]byte),
	}
	copy(p.PublicKey[:], publicKey)
	copy(p.SecretKey[:], secretKey)
	return p, nil
}

func New() (*Profile, error) {
	publicKey, secretKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	return &Profile{PublicKey: publicKey, SecretKey: secretKey}, nil
}

func Exists() bool {
	_, err := os.Stat(profilePath)
	return err == nil
}

func Save(p *Profile) error {
	temp := profileFormat{
		PublicKey: hex.EncodeToString(p.PublicKey[:]),
		SecretKey: hex.EncodeToString(p.SecretKey[:]),
	}

	bytes, err := json.MarshalIndent(&temp, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(profilePath, bytes, 0666)
}
