package qbft

import (
	specqbft "github.com/ssvlabs/ssv-spec/qbft"
	spectypes "github.com/ssvlabs/ssv-spec/types"

	"github.com/ssvlabs/ssv/ssvsigner/ekm"

	"github.com/ssvlabs/ssv/protocol/v2/qbft/roundtimer"
)

type signing interface {
	// GetShareSigner returns a BeaconSigner instance
	GetShareSigner() ekm.BeaconSigner
	// GetSignatureDomainType returns the Domain type used for signatures
	GetSignatureDomainType() spectypes.DomainType
}

type IConfig interface {
	signing
	// GetProposerF returns func used to calculate proposer
	GetProposerF() specqbft.ProposerF
	// GetNetwork returns a p2p Network instance
	GetNetwork() specqbft.Network
	// GetTimer returns round timer
	GetTimer() roundtimer.Timer
	// GetRoundCutOff returns the round cut off
	GetCutOffRound() specqbft.Round
	// IsSMRMode returns true if using 2-phase SMR consensus instead of 4-phase QBFT
	IsSMRMode() bool
}

type Config struct {
	BeaconSigner ekm.BeaconSigner
	Domain       spectypes.DomainType
	ProposerF    specqbft.ProposerF
	Network      specqbft.Network
	Timer        roundtimer.Timer
	CutOffRound  specqbft.Round
	// SMRMode enables 2-phase SMR consensus instead of 4-phase QBFT.
	// In SMR mode, upon receiving a proposal, nodes send a commit directly
	// (skipping the prepare phase), achieving consensus in 2 rounds.
	SMRMode bool
}

// GetShareSigner returns a BeaconSigner instance
func (c *Config) GetShareSigner() ekm.BeaconSigner {
	return c.BeaconSigner
}

// GetSignatureDomainType returns the Domain type used for signatures
func (c *Config) GetSignatureDomainType() spectypes.DomainType {
	return c.Domain
}

// GetProposerF returns func used to calculate proposer
func (c *Config) GetProposerF() specqbft.ProposerF {
	return c.ProposerF
}

// GetNetwork returns a p2p Network instance
func (c *Config) GetNetwork() specqbft.Network {
	return c.Network
}

// GetTimer returns round timer
func (c *Config) GetTimer() roundtimer.Timer {
	return c.Timer
}

func (c *Config) GetCutOffRound() specqbft.Round {
	return c.CutOffRound
}

// IsSMRMode returns true if using 2-phase SMR consensus.
func (c *Config) IsSMRMode() bool {
	return c.SMRMode
}
