package payment

type (
	InputType string
)

// InputType compatible for gatewayConfigTemplate
var (
	Text     InputType = "text"
	Password InputType = "password"
	Textarea InputType = "textarea"

	Number InputType = "number"

	Radiogroup InputType = "radiogroup"
	Dropdown   InputType = "dropdown"
)
