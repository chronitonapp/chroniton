package utils

type Verifier interface {
	Verify() []error
}

func Verify(verifier Verifier) []error {
	return verifier.Verify()
}
