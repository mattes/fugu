package run

func Parse() {

	cmd.Bool([]string{"#rm", "-rm"}, false, "Automatically remove the container when it exits (incompatible with -d)")
	cmd.Bool([]string{"d", "-detach"}, false, "Detached mode: run container in the background and print new container ID")
	cmd.Bool([]string{"#n", "#-networking"}, true, "Enable networking for this container")
	cmd.Bool([]string{"#privileged", "-privileged"}, false, "Give extended privileges to this container")
	cmd.Bool([]string{"P", "-publish-all"}, false, "Publish all exposed ports to the host interfaces")
	cmd.Bool([]string{"i", "-interactive"}, false, "Keep STDIN open even if not attached")
	cmd.Bool([]string{"t", "-tty"}, false, "Allocate a pseudo-TTY")
	cmd.String([]string{"#cidfile", "-cidfile"}, "", "Write the container ID to the file")
	cmd.String([]string{"#entrypoint", "-entrypoint"}, "", "Overwrite the default ENTRYPOINT of the image")
	cmd.String([]string{"h", "-hostname"}, "", "Container host name")
	cmd.String([]string{"m", "-memory"}, "", "Memory limit (format: <number><optional unit>, where unit = b, k, m or g)")
	cmd.String([]string{"u", "-user"}, "", "Username or UID")
	cmd.String([]string{"w", "-workdir"}, "", "Working directory inside the container")
	cmd.Int64([]string{"c", "-cpu-shares"}, 0, "CPU shares (relative weight)")
	cmd.String([]string{"-cpuset"}, "", "CPUs in which to allow execution (0-3, 0,1)")
	cmd.String([]string{"-net"}, "bridge", "Set the Network mode for the container\n'bridge': creates a new network stack for the container on the docker bridge\n'none': no networking for this container\n'container:<name|id>': reuses another container network stack\n'host': use the host network stack inside the container.  Note: the host mode gives the container full access to local system services such as D-bus and is therefore considered insecure.")
	cmd.String([]string{"-restart"}, "", "Restart policy to apply when a container exits (no, on-failure[:max-retry], always)")
	cmd.Bool([]string{"#sig-proxy", "-sig-proxy"}, true, "Proxy received signals to the process (even in non-TTY mode). SIGCHLD, SIGSTOP, and SIGKILL are not proxied.")
	cmd.String([]string{"#name", "-name"}, "", "Assign a name to the container")

	cmd.Var(&flAttach, []string{"a", "-attach"}, "Attach to STDIN, STDOUT or STDERR.")
	cmd.Var(&flVolumes, []string{"v", "-volume"}, "Bind mount a volume (e.g., from the host: -v /host:/container, from Docker: -v /container)")
	cmd.Var(&flLinks, []string{"#link", "-link"}, "Add link to another container in the form of name:alias")
	cmd.Var(&flDevices, []string{"-device"}, "Add a host device to the container (e.g. --device=/dev/sdc:/dev/xvdc)")
	cmd.Var(&flEnv, []string{"e", "-env"}, "Set environment variables")
	cmd.Var(&flEnvFile, []string{"-env-file"}, "Read in a line delimited file of environment variables")

	cmd.Var(&flPublish, []string{"p", "-publish"}, fmt.Sprintf("Publish a container's port to the host\nformat: %s\n(use 'docker port' to see the actual mapping)", nat.PortSpecTemplateFormat))
	cmd.Var(&flExpose, []string{"#expose", "-expose"}, "Expose a port from the container without publishing it to your host")
	cmd.Var(&flDns, []string{"#dns", "-dns"}, "Set custom DNS servers")
	cmd.Var(&flDnsSearch, []string{"-dns-search"}, "Set custom DNS search domains")
	cmd.Var(&flVolumesFrom, []string{"#volumes-from", "-volumes-from"}, "Mount volumes from the specified container(s)")
	cmd.Var(&flLxcOpts, []string{"#lxc-conf", "-lxc-conf"}, "(lxc exec-driver only) Add custom lxc options --lxc-conf=\"lxc.cgroup.cpuset.cpus = 0,1\"")

	cmd.Var(&flCapAdd, []string{"-cap-add"}, "Add Linux capabilities")
	cmd.Var(&flCapDrop, []string{"-cap-drop"}, "Drop Linux capabilities")
}
