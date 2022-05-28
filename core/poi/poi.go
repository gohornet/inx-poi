package poi

import (
	"bytes"
	"crypto"
	"encoding"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/gohornet/hornet/pkg/model/milestone"
	restapipkg "github.com/gohornet/hornet/pkg/restapi"
	"github.com/gohornet/hornet/pkg/whiteflag"

	// import implementation
	_ "golang.org/x/crypto/blake2b"
)

func createProof(c echo.Context) (*ProofRequestAndResponse, error) {

	blockID, err := restapipkg.ParseBlockIDParam(c)
	if err != nil {
		return nil, err
	}

	metadata, err := deps.NodeBridge.BlockMetadata(blockID)
	if err != nil {
		return nil, err
	}

	msIndex := metadata.GetReferencedByMilestoneIndex()
	if msIndex == 0 {
		return nil, errors.WithMessagef(restapipkg.ErrInvalidParameter, "block %s is not referenced by a milestone", blockID.ToHex())
	}

	ms, err := deps.NodeBridge.Milestone(msIndex)
	if err != nil {
		return nil, err
	}

	block, err := deps.NodeBridge.Block(blockID)
	if err != nil {
		return nil, err
	}

	blockIDs, err := FetchMilestoneCone(msIndex)
	if err != nil {
		return nil, err
	}

	var blockIDIndex int
	includedBlocks := []encoding.BinaryMarshaler{}
	for i, b := range blockIDs {
		if b == blockID {
			blockIDIndex = i
		}
		includedBlocks = append(includedBlocks, b)
	}

	hasher := whiteflag.NewHasher(crypto.BLAKE2b_256)

	proof, err := hasher.ComputeInclusionProof(includedBlocks, blockIDIndex)
	if err != nil {
		return nil, err
	}

	hash := proof.Hash(hasher)

	if !bytes.Equal(hash, ms.Milestone.InclusionMerkleRoot[:]) {
		return nil, errors.WithMessage(echo.ErrInternalServerError, "valid proof cannot be created")
	}

	return &ProofRequestAndResponse{
		Milestone: ms.Milestone,
		Block:     block,
		Proof:     proof,
	}, nil
}

func validateProof(c echo.Context) (*ValidateProofResponse, error) {

	req := &ProofRequestAndResponse{}
	if err := c.Bind(req); err != nil {
		return nil, errors.WithMessagef(restapipkg.ErrInvalidParameter, "invalid request, error: %s", err)
	}

	if req.Proof == nil || req.Milestone == nil || req.Block == nil {
		return nil, errors.WithMessage(restapipkg.ErrInvalidParameter, "invalid request")
	}

	// Hash the contained block to get the ID
	blockID, err := req.Block.ID()
	if err != nil {
		return nil, err
	}

	// Check if the contained proof contains the blockID
	containsValue, err := req.Proof.ContainsValue(blockID)
	if err != nil {
		return nil, err
	}
	if !containsValue {
		return &ValidateProofResponse{Valid: false}, nil
	}

	// Verify the contained Milestone signatures
	keySet := deps.KeyManager.PublicKeysSetForMilestoneIndex(milestone.Index(req.Milestone.Index))
	if err := req.Milestone.VerifySignatures(deps.MilestonePublicKeyCount, keySet); err != nil {
		return &ValidateProofResponse{Valid: false}, nil
	}

	hash := req.Proof.Hash(whiteflag.NewHasher(crypto.BLAKE2b_256))

	return &ValidateProofResponse{Valid: bytes.Equal(hash, req.Milestone.InclusionMerkleRoot[:])}, nil
}