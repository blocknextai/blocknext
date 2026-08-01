package sendgrid

type sendgridAddress struct {
	Email string `json:"email"`
}

type sendgridPersonalization struct {
	To  []sendgridAddress `json:"to"`
	Cc  []sendgridAddress `json:"cc,omitempty"`
	Bcc []sendgridAddress `json:"bcc,omitempty"`
}

type sendgridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type sendgridRequest struct {
	Personalizations []sendgridPersonalization `json:"personalizations"`
	From             sendgridAddress           `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendgridContent         `json:"content"`
}

type sendgridErrorItem struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Help    string `json:"help"`
}

type sendgridErrorResponse struct {
	Errors []sendgridErrorItem `json:"errors"`
}
