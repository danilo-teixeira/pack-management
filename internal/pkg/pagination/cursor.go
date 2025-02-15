package pagination

import (
	"encoding/base64"
	"errors"
	"reflect"
	"slices"
	"strings"

	"github.com/uptrace/bun"
)

type (
	Metadata struct {
		PageSize   int    `json:"page_size"`
		NextCursor string `json:"next_cursor,omitempty"`
		PrevCursor string `json:"prev_cursor,omitempty"`
	}

	CursorConfig struct {
		PageSize        int
		PageCursor      *string
		CursorDirection string
		OrderStrategy   string
		CursorField     string
	}
)

var (
	ErrInvalidCursor      = errors.New("invalid cursor")
	ErrInvalidCursorField = errors.New("cursor field is not a string")
)

const (
	AscDirection  = "ASC"
	DescDirection = "DESC"
)

func EncodeCursor(order string, cursor string) string {
	return base64.StdEncoding.EncodeToString([]byte(order + ":" + cursor))
}

func DecodeCursor(encodedCursor string) (string, string, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return "", "", err
	}

	cursor := string(decodedCursor)
	cursorParts := splitCursor(cursor)

	if len(cursorParts) < 2 {
		return "", "", ErrInvalidCursor
	}

	return cursorParts[0], cursorParts[1], nil
}

func InvertDirection(direction string) string {
	if direction == AscDirection {
		return DescDirection
	}

	return AscDirection
}

func getItemValue[T any](item T, cursorField string) (string, error) {
	reflectedValue := reflect.ValueOf(item)

	if reflectedValue.Kind() == reflect.Ptr {
		reflectedValue = reflectedValue.Elem()
	}

	cursorFieldInterface := reflectedValue.FieldByName(cursorField).Interface()

	value, ok := cursorFieldInterface.(string)
	if ok {
		return value, nil
	}

	return "", ErrInvalidCursorField
}

func BuildMetadata[T any](config CursorConfig, items []T) ([]T, *Metadata, error) {
	hasMoreItems := len(items) > config.PageSize
	if hasMoreItems {
		items = items[:len(items)-1]
	}

	if config.CursorDirection != config.OrderStrategy {
		slices.Reverse(items)
	}

	var nextCursor string
	if hasMoreItems || config.CursorDirection != config.OrderStrategy {
		lastItem, err := getItemValue(items[len(items)-1], config.CursorField)
		if err != nil {
			return nil, nil, err
		}

		nextCursor = EncodeCursor(config.OrderStrategy, lastItem)
	}

	var prevCursor string
	if len(items) > 0 && config.PageCursor != nil && (hasMoreItems || config.CursorDirection == config.OrderStrategy) {
		firstItem, err := getItemValue(items[0], config.CursorField)
		if err != nil {
			return nil, nil, err
		}

		prevCursor = EncodeCursor(InvertDirection(config.OrderStrategy), firstItem)
	}

	return items, &Metadata{
		PageSize:   config.PageSize,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
	}, nil
}

func decodeCursor(config CursorConfig) (string, string, error) {
	if config.PageCursor == nil {
		return config.OrderStrategy, "", nil
	}

	cursorDirection, cursorValue, err := DecodeCursor(*config.PageCursor)
	if err != nil {
		return "", "", err
	}

	return cursorDirection, cursorValue, nil
}

func BuildCursorQuery(config CursorConfig, query *bun.SelectQuery) (*bun.SelectQuery, string, error) {
	cursorDirection, cursorValue, err := decodeCursor(config)
	if err != nil {
		return nil, "", err
	}

	if cursorValue != "" {
		if cursorDirection == DescDirection {
			query.Where("pack.id < ?", cursorValue)
		} else {
			query.Where("pack.id > ?", cursorValue)
		}
	}

	query.Order("pack.id " + cursorDirection)

	return query, cursorDirection, nil
}

func splitCursor(cursor string) []string {
	return strings.Split(cursor, ":")
}
