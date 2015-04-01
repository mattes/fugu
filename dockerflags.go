package fugu

import (
	"fmt"
	"github.com/docker/docker/nat"
	"github.com/mattes/go-collect/flags"
)

var DockerFlags = make(map[string]*flags.Flags)

func init() {

	// Copy these values from
	// https://github.com/docker/docker/tree/master/api/client

	// Define DockerFlags["build"]
	DockerFlags["build"] = flags.New("docker")
	DockerFlags["build"].String([]string{"t", "-tag"}, "", "Tag for the image (!)")
	DockerFlags["build"].Bool([]string{"q", "-quiet"}, false, "Suppress the verbose output generated by the containers")
	DockerFlags["build"].Bool([]string{"#no-cache", "-no-cache"}, false, "Do not use cache when building the image")
	DockerFlags["build"].Bool([]string{"#rm", "-rm"}, true, "Remove intermediate containers after a successful build")
	DockerFlags["build"].Bool([]string{"-force-rm"}, false, "Always remove intermediate containers")
	DockerFlags["build"].Bool([]string{"-pull"}, false, "Always attempt to pull a newer version of the image")
	DockerFlags["build"].String([]string{"f", "-file"}, "", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")

	// Define DockerFlags["run"]
	DockerFlags["run"] = flags.New("docker")
	DockerFlags["run"].Bool([]string{"#rm", "-rm"}, false, "Automatically remove the container when it exits (incompatible with -d)")
	DockerFlags["run"].Bool([]string{"d", "-detach"}, false, "Detached mode: run the container in the background and print the new container ID")
	DockerFlags["run"].Bool([]string{"#sig-proxy", "-sig-proxy"}, true, "Proxy received signals to the process (non-TTY mode only). SIGCHLD, SIGSTOP, and SIGKILL are not proxied.")
	DockerFlags["run"].String([]string{"#name", "-name"}, "", "Assign a name to the container")

	DockerFlags["run"].Var([]string{"a", "-attach"}, "Attach to STDIN, STDOUT or STDERR.")
	DockerFlags["run"].Var([]string{"v", "-volume"}, "Bind mount a volume (e.g., from the host: -v /host:/container, from Docker: -v /container)")
	DockerFlags["run"].Var([]string{"#link", "-link"}, "Add link to another container in the form of <name|id>:alias")
	DockerFlags["run"].Var([]string{"-device"}, "Add a host device to the container (e.g. --device=/dev/sdc:/dev/xvdc:rwm)")

	DockerFlags["run"].Var([]string{"e", "-env"}, "Set environment variables")
	DockerFlags["run"].Var([]string{"-env-file"}, "Read in a line delimited file of environment variables")

	DockerFlags["run"].Var([]string{"p", "-publish"}, fmt.Sprintf("Publish a container's port to the host\nformat: %s\n(use 'docker port' to see the actual mapping)", nat.PortSpecTemplateFormat))
	DockerFlags["run"].Var([]string{"#expose", "-expose"}, "Expose a port or a range of ports (e.g. --expose=3300-3310) from the container without publishing it to your host")
	DockerFlags["run"].Var([]string{"#dns", "-dns"}, "Set custom DNS servers")
	DockerFlags["run"].Var([]string{"-dns-search"}, "Set custom DNS search domains (Use --dns-search=. if you don't wish to set the search domain)")
	DockerFlags["run"].Var([]string{"-add-host"}, "Add a custom host-to-IP mapping (host:ip)")
	DockerFlags["run"].Var([]string{"#volumes-from", "-volumes-from"}, "Mount volumes from the specified container(s)")
	DockerFlags["run"].Var([]string{"#lxc-conf", "-lxc-conf"}, "(lxc exec-driver only) Add custom lxc options --lxc-conf=\"lxc.cgroup.cpuset.cpus = 0,1\"")

	DockerFlags["run"].Var([]string{"-cap-add"}, "Add Linux capabilities")
	DockerFlags["run"].Var([]string{"-cap-drop"}, "Drop Linux capabilities")
	DockerFlags["run"].Var([]string{"-security-opt"}, "Security Options")

	DockerFlags["run"].Bool([]string{"#n", "#-networking"}, true, "Enable networking for this container")
	DockerFlags["run"].Bool([]string{"#privileged", "-privileged"}, false, "Give extended privileges to this container")
	DockerFlags["run"].String([]string{"-pid"}, "", "Default is to create a private PID namespace for the container\n'host': use the host PID namespace inside the container.  Note: the host mode gives the container full access to processes on the system and is therefore considered insecure.")
	DockerFlags["run"].Bool([]string{"P", "-publish-all"}, false, "Publish all exposed ports to random ports on the host interfaces")
	DockerFlags["run"].Bool([]string{"i", "-interactive"}, false, "Keep STDIN open even if not attached")
	DockerFlags["run"].Bool([]string{"t", "-tty"}, false, "Allocate a pseudo-TTY")
	DockerFlags["run"].String([]string{"#cidfile", "-cidfile"}, "", "Write the container ID to the file")
	DockerFlags["run"].String([]string{"#entrypoint", "-entrypoint"}, "", "Overwrite the default ENTRYPOINT of the image")
	DockerFlags["run"].String([]string{"h", "-hostname"}, "", "Container host name")
	DockerFlags["run"].String([]string{"m", "-memory"}, "", "Memory limit (format: <number><optional unit>, where unit = b, k, m or g)")
	DockerFlags["run"].String([]string{"-memory-swap"}, "", "Total memory usage (memory + swap), set '-1' to disable swap (format: <number><optional unit>, where unit = b, k, m or g)")
	DockerFlags["run"].String([]string{"u", "-user"}, "", "Username or UID")
	DockerFlags["run"].String([]string{"w", "-workdir"}, "", "Working directory inside the container")
	DockerFlags["run"].Int64([]string{"c", "-cpu-shares"}, 0, "CPU shares (relative weight)")
	DockerFlags["run"].String([]string{"-cpuset"}, "", "CPUs in which to allow execution (0-3, 0,1)")
	DockerFlags["run"].String([]string{"-net"}, "bridge", "Set the Network mode for the container\n'bridge': creates a new network stack for the container on the docker bridge\n'none': no networking for this container\n'container:<name|id>': reuses another container network stack\n'host': use the host network stack inside the container.  Note: the host mode gives the container full access to local system services such as D-bus and is therefore considered insecure.")
	DockerFlags["run"].String([]string{"-mac-address"}, "", "Container MAC address (e.g. 92:d0:c6:0a:29:33)")
	DockerFlags["run"].String([]string{"-ipc"}, "", "Default is to create a private IPC namespace (POSIX SysV IPC) for the container\n'container:<name|id>': reuses another container shared memory, semaphores and message queues\n'host': use the host shared memory,semaphores and message queues inside the container.  Note: the host mode gives the container full access to local shared memory and is therefore considered insecure.")
	DockerFlags["run"].String([]string{"-restart"}, "", "Restart policy to apply when a container exits (no, on-failure[:max-retry], always)")
	DockerFlags["run"].Bool([]string{"-read-only"}, false, "Mount the container's root filesystem as read only")

	// Define DockerFlags["exec"]
	DockerFlags["exec"] = flags.New("docker")
	DockerFlags["exec"].Bool([]string{"i", "-interactive"}, false, "Keep STDIN open even if not attached")
	DockerFlags["exec"].Bool([]string{"t", "-tty"}, false, "Allocate a pseudo-TTY")
	DockerFlags["exec"].Bool([]string{"d", "-detach"}, false, "Detached mode: run command in the background")

	// Define DockerFlags["destroy"]
	DockerFlags["destroy"] = flags.New("docker")

	// Define DockerFlags["pull"]
	DockerFlags["pull"] = flags.New("docker")
	DockerFlags["pull"].Bool([]string{"a", "-all-tags"}, false, "Download all tagged images in the repository")

	// Define DockerFlags["push"]
	DockerFlags["push"] = flags.New("docker")

}
