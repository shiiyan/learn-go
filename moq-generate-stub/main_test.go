package main

import (
	"errors"
	"testing"
)

func TestCompleteSignUp(t *testing.T) {
	t.Parallel()

	type Expected struct {
		To      string
		Subject string
		Body    string
	}

	tests := map[string]struct {
		Name     string
		Expected Expected
	}{
		"userA": {
			Name:     "user_a",
			Expected: Expected{To: "user_a@example.com", Subject: "Subject", Body: "Body"},
		},
		"userB": {
			Name:     "user_b",
			Expected: Expected{To: "user_b@example.com", Subject: "Subject", Body: "Body"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fetcherStub := &UserFetcherMock{}
			fetcherStub.GetFunc = func(name string) User {
				return User{Name: name, Email: name + "@example.com"}
			}
			senderMock := &EmailSenderMock{}

			CompleteSignUp(tt.Name, fetcherStub, senderMock)

			sendCalls := senderMock.SendCalls()
			if len(sendCalls) != 1 {
				t.Errorf("send was called %d times", len(sendCalls))
			}
			if sendCalls[0] != tt.Expected {
				t.Errorf("unexpected send %+v", sendCalls[0])
			}
		})
	}
}

func TestCompleteSignUpSharedMock(t *testing.T) {
	t.Parallel()

	var ErrFailedToSend = errors.New("failed to send")

	fetcherStub := &UserFetcherMock{}
	fetcherStub.GetFunc = func(name string) User {
		return User{Name: name, Email: name + "@example.com"}
	}
	senderMock := &EmailSenderMock{}
	senderMock.SendFunc = func(to, subject, body string) error {
		if to == "error@example.com" {
			return ErrFailedToSend
		}

		return nil
	}

	type Expected struct {
		To      string
		Subject string
		Body    string
	}

	tests := map[string]struct {
		Name   string
		Exp    Expected
		ExpErr error
	}{
		"userA": {
			Name: "user_a",
			Exp:  Expected{To: "user_a@example.com", Subject: "Subject", Body: "Body"},
		},
		"userB": {
			Name: "user_b",
			Exp:  Expected{To: "user_b@example.com", Subject: "Subject", Body: "Body"},
		},
		"error": {
			Name:   "error",
			ExpErr: ErrFailedToSend,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := CompleteSignUp(tt.Name, fetcherStub, senderMock)

			sendCalls := senderMock.SendCalls()

			if tt.ExpErr != nil {
				if !errors.Is(actual, ErrFailedToSend) {
					t.Fatal("expected error to be ErrFailedToSend")
				}
			} else {
				if len(sendCalls) != 1 {
					t.Errorf("send was called %d times", len(sendCalls))
				}
				if sendCalls[0] != tt.Exp {
					t.Errorf("unexpected send %+v", sendCalls[0])
				}
			}

			senderMock.ResetCalls()
		})
	}
}
