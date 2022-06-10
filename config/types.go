package config

import "github.com/6uf/apiGO"

type Encrypt struct {
	Key         []byte // 32 char key used to encrypt the strings
	PassThrough string // Main phrase to encrypt or decrypt.
	Encrypted   string // The encrypted phrase
	Decrypted   string // the string that is decrypted
	Decrypt     bool   // if true, it decrypts only.
}

type Decrypted struct {
	Bearer string
	Info   apiGO.UserINFO
}

type Config struct {
	SavedConfigs []Logs `json:"logs"`
	Key          string `json:"decodeKey"`
}

type Logs struct {
	Name       string `json:"savename"`
	Bearer     string `json:"bearer"`
	LastAuthed int64  `json:"lastauthed"`
	Info       apiGO.UserINFO
}
