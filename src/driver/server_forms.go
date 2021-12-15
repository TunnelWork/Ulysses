package driver

// ProvisioningAccount
type (
	FormUpdateProvisioningAccount struct {
		SerialNumber         uint64
		AccountConfiguration interface{}
	}

	FormDeleteProvisioningAccount struct {
		SerialNumber uint64
	}

	FormSuspendProvisioningAccount struct {
		SerialNumber uint64
	}

	FormRefreshProvisioningAccount struct {
		SerialNumber uint64
	}
)
