package domain

import (
	"errors"
	"strings"
)

func NormalizeID(value string) (string, error) {
	cleaned := strings.TrimSpace(strings.ToLower(value))
	if cleaned == "" {
		return "", errors.New("identifier cannot be empty")
	}
	for _, character := range cleaned {
		if (character < 'a' || character > 'z') && (character < '0' || character > '9') && character != '-' && character != '_' {
			return "", errors.New("identifier contains unsupported characters")
		}
	}
	return cleaned, nil
}

func NormalizeDisplayName(value string) (string, error) {
	cleaned := strings.Join(strings.Fields(value), " ")
	if len(cleaned) < 2 {
		return "", errors.New("display name must contain at least two characters")
	}
	if len(cleaned) > 80 {
		return "", errors.New("display name is too long")
	}
	return cleaned, nil
}

func EnsureDistinctPlayers(black, white string) error {
	if strings.EqualFold(strings.TrimSpace(black), strings.TrimSpace(white)) {
		return errors.New("black and white players must differ")
	}
	return nil
}
