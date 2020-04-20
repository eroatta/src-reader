package entity

import (
	"fmt"
	"go/token"
	"strings"
)

func NewDeclarationIDBuilder() *DeclarationIDBuilder {
	return &DeclarationIDBuilder{}
}

type DeclarationIDBuilder struct {
	filename string
	pkg      string
	name     string
	receiver string
	declType token.Token
}

func (b *DeclarationIDBuilder) WithFilename(filename string) *DeclarationIDBuilder {
	b.filename = filename
	return b
}

func (b *DeclarationIDBuilder) WithPackage(pkg string) *DeclarationIDBuilder {
	b.pkg = pkg
	return b
}

func (b *DeclarationIDBuilder) WithName(name string) *DeclarationIDBuilder {
	b.name = name
	return b
}

func (b *DeclarationIDBuilder) WithReceiver(recv string) *DeclarationIDBuilder {
	b.receiver = recv
	return b
}

func (b *DeclarationIDBuilder) WithType(declType token.Token) *DeclarationIDBuilder {
	b.declType = declType
	return b
}

func (b *DeclarationIDBuilder) Build() string {
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
