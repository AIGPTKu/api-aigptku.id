package domain

import (
	domainCt "github.com/AIGPTku/api-aigptku.id/app/controller/domain"
)

type AskContent = domainCt.AskContent

type RequestAsk struct {
	FuncCall chan domainCt.FuncCall
	Result chan string
	Finish chan bool
	AskContent []AskContent
	UseDefaultSystem bool
	UseFunction bool
}