package goclient

import (
	"fmt"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	eth2apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
)

// GetBeaconBlock returns beacon block by the given slot and committee index
func (gc *goClient) GetBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (ssz.Marshaler, spec.DataVersion, error) {
	// TODO need to support blinded?
	// TODO what with fee recipient?
	sig := phase0.BLSSignature{}
	copy(sig[:], randao[:])

	reqStart := time.Now()
	beaconBlock, err := gc.client.BeaconBlockProposal(gc.ctx, slot, sig, graffiti)
	if err != nil {
		return nil, DataVersionNil, err
	}
	metricsProposerDataRequest.Observe(time.Since(reqStart).Seconds())

	switch beaconBlock.Version {
	case spec.DataVersionPhase0:
		return beaconBlock.Phase0, beaconBlock.Version, nil
	case spec.DataVersionAltair:
		return beaconBlock.Altair, beaconBlock.Version, nil
	case spec.DataVersionBellatrix:
		return beaconBlock.Bellatrix, beaconBlock.Version, nil
	case spec.DataVersionCapella:
		return beaconBlock.Capella, beaconBlock.Version, nil
	default:
		return nil, DataVersionNil, errors.New(fmt.Sprintf("beacon block version %s not supported", beaconBlock.Version))
	}
}

func (gc *goClient) GetBlindedBeaconBlock(slot phase0.Slot, committeeIndex phase0.CommitteeIndex, graffiti, randao []byte) (ssz.Marshaler, spec.DataVersion, error) {
	return nil, DataVersionNil, nil
}

func (gc *goClient) SubmitBlindedBeaconBlock(block *api.VersionedBlindedBeaconBlock, sig phase0.BLSSignature) error {
	signedBlock := &api.VersionedSignedBlindedBeaconBlock{
		Version: block.Version,
	}
	switch block.Version {
	case spec.DataVersionBellatrix:
		if block.Bellatrix == nil {
			return errors.New("bellatrix blinded block is nil")
		}
		signedBlock.Bellatrix = &apiv1bellatrix.SignedBlindedBeaconBlock{
			Message: block.Bellatrix,
		}
		copy(signedBlock.Bellatrix.Signature[:], sig[:])
	case spec.DataVersionCapella:
		if block.Capella == nil {
			return errors.New("capella blinded block is nil")
		}
		signedBlock.Capella = &apiv1capella.SignedBlindedBeaconBlock{
			Message: block.Capella,
		}
		copy(signedBlock.Capella.Signature[:], sig[:])
	default:
		return errors.New("unknown block version")
	}

	return gc.client.SubmitBlindedBeaconBlock(gc.ctx, signedBlock)
}

// SubmitBeaconBlock submit the block to the node
func (gc *goClient) SubmitBeaconBlock(block *spec.VersionedBeaconBlock, sig phase0.BLSSignature) error {
	signedBlock := &spec.VersionedSignedBeaconBlock{
		Version: block.Version,
	}
	switch block.Version {
	case spec.DataVersionPhase0:
		if block.Phase0 == nil {
			return errors.New("phase0 block is nil")
		}
		signedBlock.Phase0 = &phase0.SignedBeaconBlock{
			Message: block.Phase0,
		}
		copy(signedBlock.Phase0.Signature[:], sig[:])
	case spec.DataVersionAltair:
		if block.Altair == nil {
			return errors.New("altair block is nil")
		}
		signedBlock.Altair = &altair.SignedBeaconBlock{
			Message: block.Altair,
		}
		copy(signedBlock.Altair.Signature[:], sig[:])
	case spec.DataVersionBellatrix:
		if block.Bellatrix == nil {
			return errors.New("bellatrix block is nil")
		}
		signedBlock.Bellatrix = &bellatrix.SignedBeaconBlock{
			Message: block.Bellatrix,
		}
		copy(signedBlock.Bellatrix.Signature[:], sig[:])
	case spec.DataVersionCapella:
		if block.Capella == nil {
			return errors.New("capella block is nil")
		}
		signedBlock.Capella = &capella.SignedBeaconBlock{
			Message: block.Capella,
		}
		copy(signedBlock.Capella.Signature[:], sig[:])
	default:
		return errors.New("unknown block version")
	}

	return gc.client.SubmitBeaconBlock(gc.ctx, signedBlock)
}

func (gc *goClient) SubmitProposalPreparation(feeRecipients map[phase0.ValidatorIndex]bellatrix.ExecutionAddress) error {
	var preparations []*eth2apiv1.ProposalPreparation
	for index, recipient := range feeRecipients {
		preparations = append(preparations, &eth2apiv1.ProposalPreparation{
			ValidatorIndex: index,
			FeeRecipient:   recipient,
		})
	}
	return gc.client.SubmitProposalPreparations(gc.ctx, preparations)
}
