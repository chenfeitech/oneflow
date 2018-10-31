package config

import (
	"os"
	"strings"
	"time"
)

var provision = _provision{strings.HasPrefix(strings.ToLower(os.Getenv("FLOW_ENV")), "prod")}

type _provision struct {
	is_online bool
}

func (p *_provision) String(prod, dev string) string {
	if p.is_online {
		return prod
	} else {
		return dev
	}
}

func (p *_provision) Uint(prod, dev uint) uint {
	if p.is_online {
		return prod
	} else {
		return dev
	}
}

func (p *_provision) Uint64(prod, dev uint64) uint64 {
	if p.is_online {
		return prod
	} else {
		return dev
	}
}

func (p *_provision) Bool(prod, dev bool) bool {
	if p.is_online {
		return prod
	} else {
		return dev
	}
}

func (p *_provision) Duration(prod, dev time.Duration) time.Duration {
	if p.is_online {
		return prod
	} else {
		return dev
	}
}

func (p *_provision) Float64(prod, dev float64) float64 {
	if p.is_online {
		return prod
	} else {
		return dev
	}
}

func String(prod, dev string) string {
	return provision.String(prod, dev)
}

func Uint(prod, dev uint) uint {
	return provision.Uint(prod, dev)
}

func Uint64(prod, dev uint64) uint64 {
	return provision.Uint64(prod, dev)
}

func Bool(prod, dev bool) bool {
	return provision.Bool(prod, dev)
}

func Duration(prod, dev time.Duration) time.Duration {
	return provision.Duration(prod, dev)
}

func Float64(prod, dev float64) float64 {
	return provision.Float64(prod, dev)
}
