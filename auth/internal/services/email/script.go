package email

import (
	"bytes"

	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

func (es *emailService) SendPasswordResetToken(email, token string) error {
	tracer := EmailErrorTracer + ": SendPasswordResetToken()"

	var body bytes.Buffer
	data := struct{ URL string }{URL: utils.GenerateURLWithTokenQuery("/auth/reset-password", token)}
	if err := es.template.ExecuteTemplate(&body, "password_reset.html", data); err != nil {
		return ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
	}

	message := es.BuildMessage([]string{email}, "Password Reset Request", body.String(), "")
	return es.SendEmail(message)
}
