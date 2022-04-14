package driver

// Affiliation
type (
	FormAffiliation struct {
	}
)

// MFA
type (
	FormMfaInitSignUp struct {
		Type string `json:"type" binding:"required"`
	}

	FormMfaCompleteSignUp struct {
		Type     string            `json:"type" binding:"required"`
		Response map[string]string `json:"response" binding:"required"`
	}

	FormMfaNewChallenge struct {
		Type string `json:"type"`
	}

	// Checkpoint Only
	FormMfaSubmitChallenge struct {
		Mfa *MfaChallengeResponse `json:"mfa" binding:"required"`
	}

	MfaChallengeResponse struct {
		Type     string            `json:"type" binding:"required"`
		Response map[string]string `json:"response" binding:"required"`
	}
)

// User
type (
	FormCreateUser struct {
		Email       string `json:"email" binding:"required"`
		PublicKey   string `json:"public_key" binding:"required"`
		Affiliation uint64 `json:"affiliation"`

		Signature string `json:"signature" binding:"required"` // Signature == Sign("REGISTER")
	}

	FormUpdateUser struct {
		ID          uint64 `json:"id" binding:"required"`
		Email       string `json:"email"`
		PublicKey   string `json:"public_key"`
		Role        uint32 `json:"role"`
		Affiliation uint64 `json:"affiliation"`
	}

	FormChangePublicKey struct {
	}
)
