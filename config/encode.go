package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/6uf/apiGO"
)

func (s *Config) UploadAccount(savename, bearer, uuid, name string) {
	s.SavedConfigs = append(s.SavedConfigs, Logs{
		Name:   savename,
		Bearer: s.Encode(bearer),
		Info: apiGO.UserINFO{
			Name: s.Encode(name),
			ID:   s.Encode(uuid),
		},
		LastAuthed: time.Now().Unix(),
	})
	s.SaveConfig()
	s.LoadState()
}

func (s *Config) Encode(value string) string {
	var D Encrypt = Encrypt{
		Key:         []byte(s.Key),
		PassThrough: value,
	}
	D.ParseValue()
	return strings.Trim(D.Encrypted, "\n")
}

func (s *Config) Decode(value string) string {
	var D Encrypt = Encrypt{
		Key:         []byte(s.Key),
		PassThrough: value,
		Decrypt:     true,
	}
	D.ParseValue()
	return strings.Trim(D.Decrypted, "\n")
}

func (s *Config) GetValueFromConfig(name string) Decrypted {
	for _, data := range s.SavedConfigs {
		if data.Name == name {
			return Decrypted{
				Bearer: s.Decode(data.Bearer),
				Info: apiGO.UserINFO{
					ID:   s.Decode(data.Info.ID),
					Name: s.Decode(data.Info.Name),
				},
			}
		}
	}
	return Decrypted{}
}

func (Data *Encrypt) ParseValue() error {
	if !Data.Decrypt {
		block, err := aes.NewCipher(Data.Key)
		if err != nil {
			return err
		}
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}
		Data.Encrypted = fmt.Sprintf("%x\n", aesgcm.Seal(nil, make([]byte, 12), []byte(Data.PassThrough), nil))
	} else {
		block, err := aes.NewCipher(Data.Key)
		if err != nil {
			return err
		}
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}
		Nonce, err := hex.DecodeString(Data.PassThrough)
		if err != nil {
			return err
		}
		plaintext, err := aesgcm.Open(nil, make([]byte, 12), Nonce, nil)
		if err != nil {
			return err
		}
		Data.Decrypted = fmt.Sprintf("%s", string(plaintext))
	}
	return nil
}
