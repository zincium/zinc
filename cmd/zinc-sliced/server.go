package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zincium/zinc/modules/env"
	"github.com/zincium/zinc/sliced"
	"google.golang.org/grpc"
)

type ServerOptions struct {
	GitPath string
	Root    string
}

type Server struct {
	*grpc.Server
	opt *ServerOptions
}

func NewServer(opt *ServerOptions) *Server {
	return &Server{opt: opt, Server: grpc.NewServer()}
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.Server == nil {
		s.Server.GracefulStop()
	}

}

func (s *Server) ListenAndServe(listen string) error {
	sliced.RegisterSlicerServer(s, s)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

func (s *Server) UploadPack(stream sliced.Slicer_UploadPackServer) error {
	ctx := stream.Context()
	pack, err := stream.Recv()
	if err != nil {
		return err
	}
	repoPath := filepath.Join(s.opt.Root, pack.Repo.Location)
	cmd := exec.CommandContext(ctx, "git", "upload-pack", repoPath)
	cmd.Env = append(env.Environ(), "GIT_ZINC_PROTOCOL=ssh")
	if len(pack.Protocol) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+pack.Protocol)
	}
	fmt.Fprintf(os.Stderr, "args: %s\n", strings.Join(cmd.Args, " "))
	return nil
}

func (s *Server) ReceivePack(stream sliced.Slicer_ReceivePackServer) error {
	return nil
}

func (s *Server) UploadArchive(stream sliced.Slicer_UploadArchiveServer) error {
	return nil
}

func (s *Server) InfoRefsUploadPack(req *sliced.InfoRefsRequest, stream sliced.Slicer_InfoRefsUploadPackServer) error {
	return nil
}

func (s *Server) InfoRefsReceivePack(req *sliced.InfoRefsRequest, stream sliced.Slicer_InfoRefsReceivePackServer) error {
	return nil
}

func (s *Server) PostUploadPack(stream sliced.Slicer_PostUploadPackServer) error {

	return nil
}
func (s *Server) PostReceivePack(stream sliced.Slicer_PostReceivePackServer) error {
	return nil
}
