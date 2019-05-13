// Package auth is a stub, providing fake authentication.
// In real world this package must provide interface to authentication subsystem.
// For testing there are three "valid" tokens:
// correct-token-1, correct-token-2, correct-token-3
package auth

import "errors"

type User struct {
        Username string
}

var ErrNotFound = errors.New("user not found")

type TokenAuthenticator struct{
        sessions map[string]*User
}

func New() *TokenAuthenticator {
        var sessions = map[string]*User {
                "correct-token-1": &User{"user1"},
                "correct-token-2": &User{"user2"},
                "correct-token-3": &User{"user3"},
        }
        return &TokenAuthenticator{sessions}
}

func (t *TokenAuthenticator) Check(token string) (*User, error) {
        user, ok := t.sessions[token]
        if !ok {
                return nil, ErrNotFound
        }
        return user, nil
}

func (t *TokenAuthenticator) Add(token string, user *User) {
        t.sessions[token] = user
}

func (t *TokenAuthenticator) Delete(token string) {
        delete(t.sessions, token)
}
