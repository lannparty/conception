# Conception
Run an ubuntu:latest Docker container within a Docker container.
# Why?
Gives access to debugging tools the parent container otherwise wouldn't have. 
# How?
Drops a chroot-ready tarball of ubuntu:latest in the root of your target container. Since chroot is a syscall, even the most lightweight containers should come equipped.

# Usage:
```
$ conception <TARGET CONTAINER ID>
```  
# Example with busybox target:
```
➜  conception git:(master) ✗ docker ps  
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
3e91f8e398b4        busybox             "sh"                13 seconds ago      Up 13 seconds                           hungry_moore
➜  conception git:(master) ✗ ./conception 3e91f8e398b4  
{"status":"Pulling from library/ubuntu","id":"latest"}
{"status":"Digest: sha256:b88f8848e9a1a4e4558ba7cfc4acc5879e1d0e7ac06401409062ad2627e6fb58"}
{"status":"Status: Image is up to date for ubuntu:latest"}
➜  conception git:(master) ✗ docker exec -it 3e91f8e398b4 sh  
/ # ls
backup.tar  bin         dev         etc         home        proc        root        sys         tmp         usr         var
/ # mkdir conception  
/ # tar xf backup.tar -C conception  
/ # chroot conception  
root@3e91f8e398b4:/#  
```

You now have access to apt to download any packages you'd need.
Note: If your target container is on Kubernetes or any other container management ecosystem you will need to point /etc/resolve to forward to the correct DNS servers.
