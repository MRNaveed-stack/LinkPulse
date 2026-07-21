package config

import "golang.org/x/time/rate"

const (
	RegisterRate        = rate.Limit(0.1) // 6 per minute
	RegisterBurst       = 3
	LoginRate           = rate.Limit(0.2) // 12 per minute
	LoginBurst          = 5
	ForgotPasswordRate  = rate.Limit(0.05) // 3 per minute
	ForgotPasswordBurst = 3
	RedirectRate        = rate.Limit(100.0)
	RedirectBurst       = 200
)
