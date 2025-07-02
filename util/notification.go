package util

func GetEmailContent(mapper string) string {
	if mapper == "verifyEmail" {
		return `
Hello, %s<br><br>
To verify your Hunter account, we need to confirm your email. Please use the following One-Time Password (OTP):<br><br>
<b>OTP: %s</b><br>
ref: %s<br><br>
The OTP is expired in 15 minutes.<br><br>
Best regards,<br>
Supporter Hunter
`
	}

	return ""
}
