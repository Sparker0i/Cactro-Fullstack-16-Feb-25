package entity

// VoteIdentifier represents unique identifiers for a vote to prevent duplicates
type VoteIdentifier struct {
	IPHash          string
	FingerprintHash string
}

// Factory method for VoteIdentifier
func NewVoteIdentifier(ipHash, fingerprintHash string) VoteIdentifier {
	return VoteIdentifier{
		IPHash:          ipHash,
		FingerprintHash: fingerprintHash,
	}
}

// Validate checks if the identifier is valid
func (v VoteIdentifier) Validate() error {
	if v.IPHash == "" {
		return ErrInvalidVoteIdentifier
	}
	if v.FingerprintHash == "" {
		return ErrInvalidVoteIdentifier
	}
	return nil
}
