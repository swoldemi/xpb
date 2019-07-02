package browser

const (
	// HostEmailSelector defines the name value for the element in which to type the host's email
	HostEmailSelector = "identifier"

	// EmailSubmitSelector defines the CSS selector (by ID) used to login form's email next button
	EmailSubmitSelector = "#identifierNext"

	// HostPasswordSelector defines the name value for finding the field in which to type the login password
	HostPasswordSelector = "password"

	// LoginSubmitSelector defines the CSS selector (by ID) used to login form's final submit button
	LoginSubmitSelector = "#passwordNext"

	// OwnerRole defines the name of the Owner role
	OwnerRole = "Owner"

	// BillingRole defines the name of the Project Billing Manager role.
	BillingRole = "Project Billing Manager"

	// AddUserBtnXPATH defines the XPATH query for finding the Add user button
	AddUserBtnXPATH = "//button[contains(@aria-label,\"Add member\")]"

	// GuestEmailSelector defines the CSS selector (by ID) used to find the IAM email field
	GuestEmailSelector = "#mat-input-1"

	// HeaderSelector defines the CSS selector (by ID) used to find the arbirary header
	HeaderSelector = "#cfc-subtask-heading-0"

	// FirstRoleSelector defines the CSS selector (by ID) used to find the first role input field
	FirstRoleSelector = "#mat-form-field-label-3"

	// RolesMenuXPATH defines the XPATH query for finding the floating roles menu
	RolesMenuXPATH = "//input[@placeholder=\"undefined\"]"

	// OwnerRoleSelector defines the CSS selector used to find the OwnerRole element
	OwnerRoleSelector = "#mat-option-124"

	// AddRoleBtnXPATH defines the XPATH query for finding the button that adds additional role fields
	AddRoleBtnXPATH = "//*[@icon=\"add\"]"

	// RoleInputsXPATH defines the XPATH query for finding all fields used for role selection
	RoleInputsXPATH = "//cfc-iam-role-picker[@placeholder=\"Select a role\"]"

	// MatOptionsXPATH defines the XPATH query for finding all mat-option-X objects
	MatOptionsXPATH = "//mat-option[starts-with(@id, \"mat-option\")]"

	// RolesSubmitXPATH defines the XPATH query for finding submission button for user addition
	RolesSubmitXPATH = "//*[@cfcformsubmit=\"addRoleForm\"]"
)
