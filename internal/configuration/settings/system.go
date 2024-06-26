package settings

import (
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

// System contains settings to configure system related elements.
type System struct {
	PUID     *uint32
	PGID     *uint32
	Timezone string
}

// Validate validates System settings.
func (s System) validate() (err error) {
	return nil
}

func (s *System) copy() (copied System) {
	return System{
		PUID:     gosettings.CopyPointer(s.PUID),
		PGID:     gosettings.CopyPointer(s.PGID),
		Timezone: s.Timezone,
	}
}

func (s *System) mergeWith(other System) {
	s.PUID = gosettings.MergeWithPointer(s.PUID, other.PUID)
	s.PGID = gosettings.MergeWithPointer(s.PGID, other.PGID)
	s.Timezone = gosettings.MergeWithString(s.Timezone, other.Timezone)
}

func (s *System) overrideWith(other System) {
	s.PUID = gosettings.OverrideWithPointer(s.PUID, other.PUID)
	s.PGID = gosettings.OverrideWithPointer(s.PGID, other.PGID)
	s.Timezone = gosettings.OverrideWithString(s.Timezone, other.Timezone)
}

func (s *System) setDefaults() {
	const defaultID = 1000
	s.PUID = gosettings.DefaultPointer(s.PUID, defaultID)
	s.PGID = gosettings.DefaultPointer(s.PGID, defaultID)
}

func (s System) String() string {
	return s.toLinesNode().String()
}

func (s System) toLinesNode() (node *gotree.Node) {
	node = gotree.New("OS Alpine settings:")

	node.Appendf("Process UID: %d", *s.PUID)
	node.Appendf("Process GID: %d", *s.PGID)

	if s.Timezone != "" {
		node.Appendf("Timezone: %s", s.Timezone)
	}

	return node
}
