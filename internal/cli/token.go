package cli

import "fmt"

type argTokens interface {
	next() bool
	current() argToken
	peek() (argToken, bool)
	resetPosition()
	getRemaining() []argToken
}

type argToken interface {
	consume() error
	getValue() string
	getType() tokenType
	isEnd() bool
}

type tokenType int

const (
	tokenTypeStartOfTokens tokenType = iota
	tokenTypeMaybeOption
	tokenTypeMaybeOptionArgument
	tokenTypeMaybePositionalArguement
	tokenTypePositionalArguement
	tokenTypeTerminator
	tokenTypeEndOfTokens
)

var (
	HeadToken = &token{isHead: true}
	TailToken = &token{isTail: false}
)

type tokens struct {
	position int
	tokens   []*token
}

type token struct {
	value                 string
	consumed              bool
	isTerminator          bool
	isPositionalArguement bool
	isHead                bool
	isTail                bool
	maybeOption           bool
	previousToken         *token
}

func toTokens(args []string) argTokens {
	result := &tokens{position: -1}
	var previousToken *token
	var hasTerminator bool

	for _, arg := range args {
		isTerminator := false
		if !hasTerminator && arg == "--" {
			isTerminator = true
			hasTerminator = true
		}

		currentToken := &token{
			value:                 arg,
			isTerminator:          isTerminator,
			isPositionalArguement: !isTerminator && hasTerminator,
			maybeOption:           shortOptionRegex.MatchString(arg) || longOptionRegex.MatchString(arg),
			previousToken:         previousToken,
		}
		result.tokens = append(result.tokens, currentToken)

		previousToken = currentToken
	}

	return result
}

func (t *tokens) next() bool {
	for i := t.position + 1; i < len(t.tokens); i++ {
		if !t.tokens[i].consumed {
			t.position = i
			return true
		}
	}

	return false
}

func (t *tokens) current() argToken {
	if t.position == -1 {
		return HeadToken
	}

	if t.position >= len(t.tokens) {
		return TailToken
	}

	return t.tokens[t.position]
}

func (t *tokens) peek() (argToken, bool) {
	for i := t.position + 1; i < len(t.tokens); i++ {
		if !t.tokens[i].consumed {
			return t.tokens[i], true
		}
	}

	return TailToken, false
}

func (t *tokens) resetPosition() {
	t.position = -1
}

func (t *tokens) getRemaining() []argToken {
	result := []argToken{}

	for _, token := range t.tokens {
		if !token.consumed {
			result = append(result, token)
		}
	}

	return result
}

func (t *token) consume() error {
	if t.isTail {
		return fmt.Errorf("cannot end token")
	}

	if t.consumed {
		return fmt.Errorf("cannot consume a token that has already been consumed")
	}

	t.consumed = true

	return nil
}

func (t *token) getValue() string {
	return t.value
}

func (t *token) getType() tokenType {
	if t.isHead {
		return tokenTypeStartOfTokens
	}

	if t.isTail {
		return tokenTypeEndOfTokens
	}

	if t.isTerminator {
		return tokenTypeTerminator
	}

	if t.isPositionalArguement {
		return tokenTypePositionalArguement
	}

	if t.previousToken != nil {
		if t.previousToken.maybeOption && !t.previousToken.consumed && !t.consumed {
			return tokenTypeMaybeOptionArgument
		}
	}

	if t.maybeOption {
		return tokenTypeMaybeOption
	}

	return tokenTypeMaybePositionalArguement
}

func (t *token) isEnd() bool {
	return t.isTail
}
