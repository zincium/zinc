package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/zincium/zinc/modules/env"
	"github.com/zincium/zinc/modules/process"
	"github.com/zincium/zinc/modules/streamio"
	"github.com/zincium/zinc/sliced"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	ErrStdinDataNotEmpty = errors.New("stream: first pack stdin not empty")
)

type ServerOptions struct {
	GitPath   string
	Root      string
	HooksPath string
}

type Server struct {
	*grpc.Server
	opt *ServerOptions
	hs  *health.Server
}

func NewServer(opt *ServerOptions) *Server {
	return &Server{opt: opt, Server: grpc.NewServer(), hs: health.NewServer()}
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.Server == nil {
		s.Server.GracefulStop()
	}

}

func (s *Server) ListenAndServe(listen string) error {
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	sliced.RegisterSlicerServer(s, s)
	grpc_health_v1.RegisterHealthServer(s, s.hs)
	return s.Serve(ln)
}

func (s *Server) UploadPack(stream sliced.Slicer_UploadPackServer) error {
	ctx, cancelCtx := context.WithCancel(stream.Context())
	defer cancelCtx()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	repoPath, err := s.repoSanitize(req)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "unable find repo %v", err)
	}
	// gRPC doesn't allow concurrent writes to a stream, so we need to
	// synchronize writing stdout and stderrr.
	// HTTP not need call NewSyncWriter (http not care stderr)
	var mtx sync.Mutex
	stdout := streamio.NewSyncWriter(&mtx, func(p []byte) error {
		return stream.Send(&sliced.UploadPackResponse{Stdout: p})
	})
	stderr := streamio.NewSyncWriter(&mtx, func(p []byte) error {
		return stream.Send(&sliced.UploadPackResponse{Stderr: p})
	})
	stdin := streamio.NewReader(func() ([]byte, error) {
		pack, err := stream.Recv()
		return pack.Stdin, err
	})
	cmd := exec.CommandContext(ctx, "git", "upload-pack", repoPath)
	cmd.Env = append(env.Environ(), "ZINC_PROTOCOL=ssh")
	if len(req.Protocol) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Protocol)
	}
	// https://github.com/golang/go/blob/0d8a4bfc962a606584be0a76ed708f86b44164c7/src/os/exec/exec.go#L244
	// go exec package create pipe when stderr/stdout/stdin not file
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Stdin = stdin
	if err := cmd.Start(); err != nil {
		return status.Errorf(codes.Internal, "create command '%s' error %v", strings.Join(cmd.Args, " "), err)
	}
	defer process.Finalize(cmd)
	if err := cmd.Wait(); err != nil {
		if exitCode, ok := process.ExitCode(err); ok {
			_ = stream.Send(&sliced.UploadPackResponse{
				ExitCode: exitCode,
			})
		}
	}
	return nil
}

func (s *Server) ReceivePack(stream sliced.Slicer_ReceivePackServer) error {
	ctx, cancelCtx := context.WithCancel(stream.Context())
	defer cancelCtx()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	repoPath, err := s.repoSanitize(req)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "unable find repo %v", err)
	}
	// gRPC doesn't allow concurrent writes to a stream, so we need to
	// synchronize writing stdout and stderrr.
	// HTTP not need call NewSyncWriter (http not care stderr)
	var mtx sync.Mutex
	stdout := streamio.NewSyncWriter(&mtx, func(p []byte) error {
		return stream.Send(&sliced.ReceivePackResponse{Stdout: p})
	})
	stderr := streamio.NewSyncWriter(&mtx, func(p []byte) error {
		return stream.Send(&sliced.ReceivePackResponse{Stderr: p})
	})
	stdin := streamio.NewReader(func() ([]byte, error) {
		pack, err := stream.Recv()
		return pack.Stdin, err
	})
	cmd := exec.CommandContext(ctx,
		s.opt.GitPath,
		"-c", "core.hooksPath="+s.opt.HooksPath,
		"receive-pack", repoPath)
	cmd.Env = append(env.Environ(),
		"ZINC_PROTOCOL=ssh",
		"ZINC_GIT_REPO_ID="+strconv.FormatInt(req.Repo.Id, 10),
		"ZINC_GIT_USER_ID="+strconv.FormatInt(req.Uid, 10), // deploy uid is virtual id
	)
	if len(req.Protocol) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Protocol)
	}
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Stdin = stdin
	if err := cmd.Start(); err != nil {
		return status.Errorf(codes.Internal, "create command '%s' error %v", strings.Join(cmd.Args, " "), err)
	}
	if err := cmd.Wait(); err != nil {
		if exitCode, ok := process.ExitCode(err); ok {
			_ = stream.Send(&sliced.ReceivePackResponse{
				ExitCode: exitCode,
			})
		}
	}
	return nil
}

func (s *Server) UploadArchive(stream sliced.Slicer_UploadArchiveServer) error {
	ctx, cancelCtx := context.WithCancel(stream.Context())
	defer cancelCtx()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	repoPath, err := s.repoSanitize(req)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "unable find repo %v", err)
	}
	var mtx sync.Mutex
	stdout := streamio.NewSyncWriter(&mtx, func(p []byte) error {
		return stream.Send(&sliced.UploadArchiveResponse{Stdout: p})
	})
	stderr := streamio.NewSyncWriter(&mtx, func(p []byte) error {
		return stream.Send(&sliced.UploadArchiveResponse{Stderr: p})
	})
	stdin := streamio.NewReader(func() ([]byte, error) {
		pack, err := stream.Recv()
		return pack.Stdin, err
	})
	cmd := exec.CommandContext(ctx, s.opt.GitPath, "upload-archive", repoPath)
	cmd.Env = env.Environ()
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Stdin = stdin
	if err := cmd.Start(); err != nil {
		return status.Errorf(codes.Internal, "create command '%s' error %v", strings.Join(cmd.Args, " "), err)
	}
	if err := cmd.Wait(); err != nil {
		if exitCode, ok := process.ExitCode(err); ok {
			_ = stream.Send(&sliced.UploadArchiveResponse{
				ExitCode: exitCode,
			})
		}
	}
	return nil
}

func (s *Server) handleInfoRefs(ctx context.Context, service, repoPath string, req *sliced.InfoRefsRequest, w io.Writer) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()
	cmd := exec.CommandContext(ctx, s.opt.GitPath,
		service, "--stateless-rpc", "--advertise-refs", repoPath)
	cmd.Env = append(env.Environ(), "ZINC_PROTOCOL=http")
	if len(req.Protocol) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Protocol)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return status.Errorf(codes.Internal, "unable create pipe %v", err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return status.Errorf(codes.Internal, "create command '%s' error %v", strings.Join(cmd.Args, " "), err)
	}
	if req.Protocol != "version=2" {
		enc := pktline.NewEncoder(w)
		_ = enc.EncodeString("# service=git- " + service + "\n")
		_ = enc.Flush()
	}
	if _, err := io.Copy(w, stdout); err != nil {
		return status.Errorf(codes.Internal, "handleInfoRefs: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		if exitCode, ok := process.ExitCode(err); ok {
			fmt.Fprintf(os.Stderr, "ExitCode %d\n", exitCode)
		}
	}
	return nil
}

func (s *Server) InfoRefsUploadPack(req *sliced.InfoRefsRequest, stream sliced.Slicer_InfoRefsUploadPackServer) error {
	repoPath, err := s.joinSanitizePath(req.Repo)
	if err != nil {
		return err
	}
	w := streamio.NewWriter(func(p []byte) error {
		return stream.Send(&sliced.InfoRefsResponse{Stdout: p})
	})
	return s.handleInfoRefs(stream.Context(), "upload-pack", repoPath, req, w)
}

func (s *Server) InfoRefsReceivePack(req *sliced.InfoRefsRequest, stream sliced.Slicer_InfoRefsReceivePackServer) error {
	repoPath, err := s.joinSanitizePath(req.Repo)
	if err != nil {
		return err
	}
	w := streamio.NewWriter(func(p []byte) error {
		return stream.Send(&sliced.InfoRefsResponse{Stdout: p})
	})
	return s.handleInfoRefs(stream.Context(), "receive-pack", repoPath, req, w)
}

func (s *Server) PostUploadPack(stream sliced.Slicer_PostUploadPackServer) error {
	ctx, cancelCtx := context.WithCancel(stream.Context())
	defer cancelCtx()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	repoPath, err := s.repoSanitize(req)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "unable find repo %v", err)
	}
	stdout := streamio.NewWriter(func(p []byte) error {
		return stream.Send(&sliced.PostUploadPackResponse{Stdout: p})
	})
	stdin := streamio.NewReader(func() ([]byte, error) {
		pack, err := stream.Recv()
		return pack.Stdin, err
	})
	cmd := exec.CommandContext(ctx, s.opt.GitPath, "upload-pack", "--stateless-rpc", repoPath)
	cmd.Env = env.Environ()
	if len(req.Protocol) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Protocol)
	}
	cmd.Stdout = stdout
	cmd.Stdin = stdin
	if err := cmd.Start(); err != nil {
		return status.Errorf(codes.Internal, "create command '%s' error %v", strings.Join(cmd.Args, " "), err)
	}
	if err := cmd.Wait(); err != nil {
		if exitCode, ok := process.ExitCode(err); ok {
			fmt.Fprintf(os.Stderr, "exitcode %d", exitCode)
		}
	}
	return nil
}

func (s *Server) PostReceivePack(stream sliced.Slicer_PostReceivePackServer) error {
	ctx, cancelCtx := context.WithCancel(stream.Context())
	defer cancelCtx()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	repoPath, err := s.repoSanitize(req)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "unable find repo %v", err)
	}
	stdout := streamio.NewWriter(func(p []byte) error {
		return stream.Send(&sliced.PostReceivePackResponse{Stdout: p})
	})
	stdin := streamio.NewReader(func() ([]byte, error) {
		pack, err := stream.Recv()
		return pack.Stdin, err
	})
	cmd := exec.CommandContext(ctx, s.opt.GitPath,
		"-c", "core.hooksPath="+s.opt.HooksPath,
		"receive-pack", "--stateless-rpc", repoPath)
	cmd.Env = append(env.Environ(),
		"ZINC_PROTOCOL=http",
		"ZINC_GIT_REPO_ID="+strconv.FormatInt(req.Repo.Id, 10),
		"ZINC_GIT_USER_ID="+strconv.FormatInt(req.Uid, 10), // deploy uid is virtual id
	)
	if len(req.Protocol) != 0 {
		cmd.Env = append(cmd.Env, "GIT_PROTOCOL="+req.Protocol)
	}
	cmd.Stdout = stdout
	cmd.Stdin = stdin
	if err := cmd.Start(); err != nil {
		return status.Errorf(codes.Internal, "create command '%s' error %v", strings.Join(cmd.Args, " "), err)
	}
	if err := cmd.Wait(); err != nil {
		if exitCode, ok := process.ExitCode(err); ok {
			fmt.Fprintf(os.Stderr, "exitcode %d", exitCode)
		}
	}
	return nil
}
