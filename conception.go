package main

import (
  "os"
  "io"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/pkg/stdcopy"
  "github.com/docker/docker/client"
  "golang.org/x/net/context"
)

func main() {
  ctx := context.Background()
  cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    panic(err)
  }
  
  reader, err := cli.ImagePull(ctx, "ubuntu", types.ImagePullOptions{})
  if err != nil {
    panic(err)
  }
  io.Copy(os.Stdout, reader)
  
  resp, err := cli.ContainerCreate(ctx, &container.Config{
    Image: "ubuntu",
    Cmd:   []string{"tar", "-cvpf", "backup.tar", "--exclude=/backup.tar", "--one-file-system", "/"},
    Tty:   true,
  }, nil, nil, "")

  if err != nil {
    panic(err)
  }
  
  if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    panic(err)
  }
  
  statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
  select {
  case err := <-errCh:
    if err != nil {
      panic(err)
    }
  case <-statusCh:
  }
  
  out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
  if err != nil {
    panic(err)
  }

  readerCopy, _, _ := cli.CopyFromContainer(ctx , resp.ID, "/backup.tar")
  if err != nil {
    panic(err)
  }

  cli.CopyToContainer(ctx, string(os.Args[1]), "/", readerCopy, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true, CopyUIDGID: true})
  if err != nil {
    panic(err)
  }

  stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
