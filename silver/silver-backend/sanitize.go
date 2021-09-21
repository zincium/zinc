package main

import (
	"os"

	"github.com/zincium/zinc/modules/base"
	silver "github.com/zincium/zinc/silver/protocol"
)

type chunkFirst interface {
	GetRepo() *silver.Repository
	GetStdin() []byte
}

func (s *Server) joinSanitizePath(repo *silver.Repository) (string, error) {
	repoPath, err := base.JoinSanitizePath(s.opt.Root, repo.Location)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(repoPath); err != nil && os.IsNotExist(err) {
		return "", err
	}
	return repoPath, nil
}

func (s *Server) repoSanitize(ck chunkFirst) (string, error) {
	if ck.GetStdin() != nil {
		return "", ErrStdinDataNotEmpty
	}
	return s.joinSanitizePath(ck.GetRepo())
}
