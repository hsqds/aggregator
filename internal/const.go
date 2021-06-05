package internal

import (
	"log"
)

const (
	loggerFlags         = log.Ltime | log.Lmsgprefix
	pgMutextLogPrefix   = "\033[31mPGMutex:\t\033[0m"
	loaderLogPrefix     = "\033[36mLoader:\t\033[0m"
	repositoryLogPrefix = "\033[34mRepository:\t\033[0m"
	parserLogPrefix     = "\033[33mParser:\t\033[0m"
	aggregatorLogPrefix = "\033[32mAggregator:\t\033[0m"
)
