package utils

import "errors"

var (
	ErrTLVHeader                   = errors.New("error tlv header")
	ErrTLVReadReachMaxLength       = errors.New("error read tlv reach max length")
	ErrTLVReadPayloadInvalidLength = errors.New("error read tlv invalid length")
)
