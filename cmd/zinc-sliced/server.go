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
	Root string
}

type Server struct {
	opt *ServerOptions
	srv *grpc.Server
}

func NewServer(opt *ServerOptions) *Server {
	return &Server{opt: opt, srv: grpc.NewServer()}
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.srv == nil {
		s.srv.GracefulStop()
	}

}

func (s *Server) ListenAndServe(listen string) error {
	sliced.RegisterSlicerServer(s.srv, s)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	return s.srv.Serve(ln)
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
func (s *Server) AdvertiseRefs(stream sliced.Slicer_AdvertiseRefsServer) error {
	return nil
}
func (s *Server) PostUploadPack(stream sliced.Slicer_PostUploadPackServer) error {

	return nil
}
func (s *Server) PostReceivePack(stream sliced.Slicer_PostReceivePackServer) error {
	return nil
}
