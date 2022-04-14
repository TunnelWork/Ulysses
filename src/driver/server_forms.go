package driver

// ProvisioningAccount
type (
	FormUpdateProvisioningAccount struct {
		SerialNumber         uint64 `json:"serial_number" binding:"required"`
		AccountConfiguration interface{}
	}

	FormDeleteProvisioningAccount struct {
		SerialNumber uint64 `json:"serial_number" binding:"required"`
	}

	FormSuspendProvisioningAccount struct {
		SerialNumber uint64 `json:"serial_number" binding:"required"`
	}

	FormRefreshProvisioningAccount struct {
		SerialNumber uint64 `json:"serial_number" binding:"required"`
	}
)
