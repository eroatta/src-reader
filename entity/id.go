package entity

import (
	"fmt"
	"go/token"
	"strings"
)

// NewIDBuilder creates a new ID builder.
func NewIDBuilder() *IDBuilder {
	return &IDBuilder{}
}

// IDBuilder builds an identifier's ID from based on several properties like:
// * the file the identifier belongs to;
// * the package name it's included in;
// * the name of the identifier;
// * the receiver, in case the identifier is a method of a struct/interface;
// * the type (function/struct/variable/constant)
type IDBuilder struct {
	filename string
	pkg      string
	name     string
	receiver string
	declType token.Token
}

// WithFilename specifies the filename where the identifier is located.
func (b *IDBuilder) WithFilename(filename string) *IDBuilder {
	b.filename = filename
	return b
}

// WithPackage specifies the package the identifier belongs to.
func (b *IDBuilder) WithPackage(pkg string) *IDBuilder {
	b.pkg = pkg
	return b
}

// WithName specifies the identifier's name.
func (b *IDBuilder) WithName(name string) *IDBuilder {
	b.name = name
	return b
}

// WithReceiver specifies the name of the interfac/struct the identifier is related to.
func (b *IDBuilder) WithReceiver(recv string) *IDBuilder {
	b.receiver = recv
	return b
}

// WithType specifies the token type for the identifier.
func (b *IDBuilder) WithType(declType token.Token) *IDBuilder {
	b.declType = declType
	return b
}

// Build creates the string ID from the provided input.
func (b *IDBuilder) Build() string {
	idBuilder := strings.Builder{}
	separator := "+++"

	idBuilder.WriteString(fmt.Sprintf("filename:%s", b.filename))
	idBuilder.WriteString(separator)

	idBuilder.WriteString(fmt.Sprintf("pkg:%s", b.pkg))
	idBuilder.WriteString(separator)

	idBuilder.WriteString(fmt.Sprintf("declType:%s", b.declType))
	idBuilder.WriteString(separator)

	if b.receiver == "" {
		idBuilder.WriteString(fmt.Sprintf("name:%s", b.name))
	} else {
		idBuilder.WriteString(fmt.Sprintf("name:%s.%s", b.receiver, b.name))
	}

	return idBuilder.String()
}
