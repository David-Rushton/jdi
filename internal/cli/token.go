package cli

import "fmt"

type argToken interface {
	first() argToken
	firstAvailable() argToken
	next() argToken
	nextAvailable() argToken
	consume() error
	remaining() []argToken
	getValue() string
	getType() tokenType
	isEnd() bool
}

type tokenType int

const (
	tokenTypeMaybeOption tokenType = iota
	tokenTypeMaybeOptionArgument
	tokenTypeMaybePositionalArguement
	tokenTypePositionalArguement
	tokenTypeTerminator
	tokenTypeEnd
)

type token struct {
	value                 string
	consumed              bool
	isTerminator          bool
	isPositionalArguement bool
	isTail                bool
	maybeOption           bool
	headToken             *token
	previousToken         *token
	nextToken             *token
	tailToken             *token
}

func tokeniseArgs(args []string) argToken {
	var head *token
	var last *token
	var hasTerminator bool

	tail := &token{isTail: true}

	for i, arg := range args {
		isTerminator := false
		if !hasTerminator && arg == "--" {
			isTerminator = true
			hasTerminator = true
		}

		current := &token{
			isPositionalArguement: !isTerminator && hasTerminator,
			maybeOption:           shortOptionRegex.MatchString(arg) || longOptionRegex.MatchString(arg),
			tailToken:             tail,
			value:                 arg,
		}

		switch i {
		case 0:
			head = current
		default:
			last.nextToken = current
			current.previousToken = last
		}

		current.headToken = head

		last = current
	}

	if last != nil {
		last.nextToken = tail
	}
	tail.previousToken = last

	if head == nil {
		tail.headToken = tail
		head = tail
	}

	return head
}

func (t *token) first() argToken {
	return t.headToken
}

func (t *token) firstAvailable() argToken {
	current := t.headToken

	for !current.isTail {
		if !current.consumed {
			return current
		}

		current = current.nextToken
	}

	return current.tailToken
}

func (t *token) next() argToken {
	if t.isTail {
		return t
	}

	return t.nextToken
}

func (t *token) nextAvailable() argToken {
	current := t.nextToken

	for !current.isTail {
		if !current.consumed {
			return current
		}

		current = current.nextToken
	}

	return t.tailToken
}

func (t *token) consume() error {
	if t.isTail {
		return fmt.Errorf("cannot end token")
	}

	if t.consumed {
		return fmt.Errorf("cannot consume a token that has already been consumed")
	}

	return nil
}

func (t *token) remaining() []argToken {
	result := []argToken{}

	candidate := t.headToken
	for candidate != nil {
		if !candidate.consumed && !candidate.isTail {
			result = append(result, candidate)
		}

		candidate = candidate.nextToken
	}

	return result
}

func (t *token) getValue() string {
	return t.value
}

func (t *token) getType() tokenType {
	if t.isTail {
		return tokenTypeEnd
	}

	if t.isTerminator {
		return tokenTypeTerminator
	}

	if t.isPositionalArguement {
		return tokenTypePositionalArguement
	}

	if t.previousToken.maybeOption && !t.maybeOption {
		return tokenTypeMaybeOptionArgument
	}

	if t.maybeOption {
		return tokenTypeMaybeOption
	}

	return tokenTypeMaybePositionalArguement
}

func (t *token) isEnd() bool {
	return t.isTail
}
