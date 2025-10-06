package main

//go:generate moq -out emailsender_moq_test.go -stub -with-resets . EmailSender

type EmailSender interface {
	Send(to, subject, body string) error
}

type User struct {
	Name  string
	Email string
}

//go:generate moq -out userfetcher_moq_test.go -stub -with-resets . UserFetcher

type UserFetcher interface {
	Get(name string) User
}

func CompleteSignUp(name string, fetcher UserFetcher, sender EmailSender) error {
	user := fetcher.Get(name)
	err := sender.Send(user.Email, "Subject", "Body")
	return err
}
