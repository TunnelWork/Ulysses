package payment

// Examples
var (
	// Here's an example of gatewayConfigTemplate
	// the frontend should display a form as
	// specified by the template
	// For possible InputType, see as defined in types.go
	ExampleGatewayConfigTemplate P = P{
		"api_user": P{
			"FriendlyName": "API Username",
			"InputType":    "text",
			"Default":      "user", // Default Value
			"Description":  "Specify your API Username. Not your email address",
			"Optional":     false,
		},
		"api_token": P{
			"FriendlyName": "API Token",
			"InputType":    "password",
			//"Optional": "false", // Default: false
		},
		"api_version": P{
			"FriendlyName": "API Version",
			"InputType":    "number",
			"Default":      3,
		},
		"api_certificate": P{
			"FriendlyName": "API Certificate",
			"InputType":    "textarea",
			"Optional":     true,
		},
		"auth_mode": P{
			"FriendlyName": "Auth Mode",
			"InputType":    "radiogroup",
			// "Optional":     false,
			"Items": P{
				"auth_once": P{
					"FriendlyName": "Auth Once",
				},
				"auth_always": P{
					"FriendlyName": "Auth Always",
				},
			},
		},
		"vendor_select": P{
			"FriendlyName": "Vendor Select",
			"InputType":    "dropdown",
			"Items": P{
				"tunnelwork": P{
					"FriendlyName": "Tunnel.Work (Default)",
				},
				"gaukaswang": P{
					"FriendlyName": "Gaukas.Wang (50% Off!)",
				},
			},
		},
	}

	// This is an example for gatewayConfig corresponding to the template
	// provided above.
	// Note that all types are string.
	ExampleGatewayConfig P = P{
		"api_user":    "user",
		"api_token":   "THISISAFAKETOKENFORDEMONSTRATION",
		"apr_version": "3",
		"api_certificate": `mDMEYTlFghYJKwYBBAHaRw8BAQdA/uS2O1VY4krn4ocmQNcslLHCYPhk3/MaKoUh
		3/QCMv20EkdhdWthcyA8aUBnYXVrLmFzPoiQBBMWCAA4FiEEBduM/AI5+aeDTX3t
		ni+Jhtdvi10FAmE5RYICGyMFCwkIBwIGFQoJCAsCBBYCAwECHgECF4AACgkQni+J
		htdvi11TEAD/WuVpN/MwPZHrhdMfjy0vftvGqCeMxnMYOMqO7dqWu/EA/jgDsJO6
		9tmLgWiGJFvp5q6C6/h2Z/h+dLEliBFvhyIJtBtHYXVrYXMgV2FuZyA8aUBnYXVr
		YXMud2FuZz6IkAQTFggAOBYhBAXbjPwCOfmng0197Z4viYbXb4tdBQJhOUWwAhsj
		BQsJCAcCBhUKCQgLAgQWAgMBAh4BAheAAAoJEJ4viYbXb4tdtHcA/jinVl583X5H
		/uqWntniOVP/H/Y8BIGKA7VKixvpRoYrAQDHXgNudx55zBvxhs8uwbx50pFyKSJl
		pURMd+1CKNipD7g4BGE5RYISCisGAQQBl1UBBQEBB0Dwmfyi3YWai/M9HnGN42LX
		R+mvWH3695DHZQwzm87FZwMBCAeIeAQYFggAIBYhBAXbjPwCOfmng0197Z4viYbX
		b4tdBQJhOUWCAhsMAAoJEJ4viYbXb4tddukA/AgGRVfY8bnJJh/xfS6CzJHkvU20
		GEO3wpxOrQHqIk7vAP9SQ4BDLnDjFrTyxNOWpWuHFcvlAbdGwrKmUjq2U74WAQ==
		=8i6d`, // Gaukas: this is my PGP PUBLIC KEY (up to Oct 24, 2021), lol
		"auth_mode":     "auth_once",
		"vendor_select": "tunnelwork",
	}

	ExampleOrderCreationParams P = P{
		"ReferenceID": "TunnelWork-#109",
		"Amount": P{
			"Value":    "4.20",
			"Currency": "USD",
		},
	}

	// Will not be parsed, used for debugging purposes.
	ExampleOrderDetails P = P{
		"OrderID":     "0xDEADC0DE",
		"ReferenceID": "TunnelWork-#109",
		"Amount": P{
			"Value":    "4.20",
			"Currency": "USD",
		},
		"GatewaySpecificField1": "Some Data",
		"GatewaySpecificField2": "Some More Data",
		"GatewaySpecificField3": "Even More Data",
	}

	// Will be parsed and also recorded
	ExampleOrderStatus P = P{
		"OrderID":         "0xDEADC0DE",
		"ReferenceID":     "TunnelWork-#109",
		"Status":          "Unpaid",        // "Unpaid", "Paid", "Refunded", "Closed"
		"PayerIdentifier": "i@gaukas.wang", // Not necessarily correlatable to the user
	}

	ExampleOrderFormTemplate P = P{
		"Type": "OnSite", // "OnSite", "Button", "CreditCard"
		"OnSiteParams": P{ // Mocking a credit card interface. For real credit card, use credit card type
			"card_holder": P{
				"FriendlyName": "Card Holder",
				"InputType":    "text",
				"Description":  "Name on card",
			},
			"card_number": P{
				"FriendlyName": "Card Number",
				"InputType":    "text",
			},
			"cvv": P{
				"FriendlyName": "CVV",
				"InputType":    "password",
				"Description":  "The 3-digit security number on your card",
				//"Optional": "false", // Default: false
			},
			"network_selection": P{
				"FriendlyName": "CC Network Selection",
				"InputType":    "dropdown",
				"Items": P{
					"mastercard": P{
						"FriendlyName": "MasterCard",
					},
					"visa": P{
						"FriendlyName": "VISA",
					},
					"jcb": P{
						"FriendlyName": "JCB",
					},
				},
			},
		},
		"ButtomParams": P{
			"btn_type":        "text", // "text", "image"
			"image_url":       "./assets/img/paynow.png",
			"btn_target_attr": "_blank",
			"btn_href":        "https://example.com/pay?id=0xDEADC0DE&merchant=0x12345678",
		},
	}

	ExampleOnSiteOrderForm P = P{
		"OrderID": "0xDEADC0DE",
		"OrderForm": P{
			"card_holder":       "Gaukas Wang",
			"card_number":       "4800333344445555",
			"cvv":               "123",
			"network_selection": "visa",
		},
	}

	ExampleOrderRefundParams P = P{
		"OrderID": "0xDEADC0DE",
		"Amount": P{
			"Value":    "0.69",
			"Currency": "USD",
		},
	}
)
