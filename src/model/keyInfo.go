package model

import (
	"github.com/boltdb/bolt"
)

type keyInfo struct {
	base
	Weight int
	Script string
}
