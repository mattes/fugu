package fugu

import (
	"github.com/mattes/go-collect"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/user"
	"testing"

	fileSource "github.com/mattes/go-collect/source/file"
)

func init() {
	collect.RegisterSource(&fileSource.File{})
}

// set to testDesc to only run this test
var RunOnly = ""

type DockerCommandTest struct {
	testDesc string
	command  string
	// dataIn   *data.Data
	argsIn []string
	strOut string
	errOut error
}

func (dct *DockerCommandTest) Test(t *testing.T) {
	if RunOnly != "" && RunOnly != dct.testDesc {
		return
	}

	c := collect.New()
	data, remainingArgs, err := c.Parse(dct.argsIn, FuguFlags[dct.command], DockerFlags[dct.command])
	assert.NoError(t, err, dct.testDesc)
	if err == nil {
		str, err := DockerCommands[dct.command](c, data, remainingArgs)
		assert.Equal(t, dct.errOut, err, dct.testDesc)
		if err == nil {
			assert.Equal(t, dct.strOut, str, dct.testDesc)
		}
	}
}

type CommandTest struct {
	testDesc       string
	command        string
	argsIn         []string
	errOut         error
	stdoutContains []string
}

func (dct *CommandTest) Test(t *testing.T) {
	if RunOnly != "" && RunOnly != dct.testDesc {
		return
	}

	c := collect.New()
	data, remainingArgs, err := c.Parse(dct.argsIn, FuguFlags[dct.command], DockerFlags[dct.command])
	assert.NoError(t, err, dct.testDesc)

	if err == nil {
		// capture stdout
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err = Commands[dct.command](c, data, remainingArgs)

		w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout

		if !assert.Equal(t, dct.errOut, err, dct.testDesc) {
			t.Logf("%s", out)
		}

		if len(dct.stdoutContains) > 0 {
			assert.NotEmpty(t, string(out), dct.testDesc)
		}
		if len(dct.stdoutContains) == 0 {
			assert.Empty(t, string(out), dct.testDesc)
		}
		for _, str := range dct.stdoutContains {
			assert.Contains(t, string(out), str, dct.testDesc)
		}
	}
}

func TestCommandBuild(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "image is missing",
		command:  "build",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "build with label",
		command:  "build",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker build --tag=redis .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "build with unknown label",
		command:  "build",
		argsIn:   []string{"label2", "--source=file://examples/fugu.labels.yml"},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "tag flag",
		command:  "build",
		argsIn:   []string{"--image=foo", "--tag=bar"},
		strOut:   "docker build --tag=foo:bar .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "no tag given",
		command:  "build",
		argsIn:   []string{"--image=foo"},
		strOut:   "docker build --tag=foo .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "given path",
		command:  "build",
		argsIn:   []string{"--image=foo", "--path=bar"},
		strOut:   "docker build --tag=foo bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "given url",
		command:  "build",
		argsIn:   []string{"--image=foo", "--url=bar"},
		strOut:   "docker build --tag=foo bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "both given path and url",
		command:  "build",
		argsIn:   []string{"--image=foo", "--url=url", "--path=path"},
		strOut:   "docker build --tag=foo url",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "path or url given via args (no label)",
		command:  "build",
		argsIn:   []string{"--image=foo", "pathOrUrl"},
		strOut:   "docker build --tag=foo pathOrUrl",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "path or url given via args (with label)",
		command:  "build",
		argsIn:   []string{"label1", "--image=foo", "--source=file://examples/fugu.labels.yml", "pathOrUrl"},
		strOut:   "docker build --tag=foo pathOrUrl",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "build: invalid number of args (no label)",
		command:  "build",
		argsIn:   []string{"--image=foo", "pathOrUrl", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "build: invalid number of args (with label label)",
		command:  "build",
		argsIn:   []string{"label1", "--image=foo", "--source=file://examples/fugu.labels.yml", "pathOrUrl", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "tag-git-branch",
		command:  "build",
		argsIn:   []string{"--image=foo", "--tag-git-branch"},
		strOut:   "docker build --tag=foo:current-branch .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "tag-git-branch overwrites given tag flag",
		command:  "build",
		argsIn:   []string{"--image=foo", "--tag=bar", "--tag-git-branch"},
		strOut:   "docker build --tag=foo:current-branch .",
		errOut:   nil,
	}).Test(t)
}

func TestCommandRun(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "image is missing",
		command:  "run",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain run",
		command:  "run",
		argsIn:   []string{"--image=foo"},
		strOut:   "docker run foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain run with label",
		command:  "run",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker run --name=my-redis redis",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain run without label",
		command:  "run",
		argsIn:   []string{"--source=file://examples/fugu.simple.yml"},
		strOut:   "docker run --name=my-ubuntu ubuntu",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain run with unknown label",
		command:  "run",
		argsIn:   []string{"label-unknown", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker run --name=my-redis redis label-unknown",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from flag",
		command:  "run",
		argsIn:   []string{"--image=foo", "--command=cmd"},
		strOut:   "docker run foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args",
		command:  "run",
		argsIn:   []string{"--image=foo", "cmd"},
		strOut:   "docker run foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args (with label)",
		command:  "run",
		argsIn:   []string{"--image=foo", "--source=file://examples/fugu.labels.yml", "cmd"},
		strOut:   "docker run --name=my-redis foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args (without label)",
		command:  "run",
		argsIn:   []string{"--image=foo", "--source=file://examples/fugu.simple.yml", "cmd"},
		strOut:   "docker run --name=my-ubuntu foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args (without label) 2",
		command:  "run",
		argsIn:   []string{"--source=file://examples/fugu.simple.yml", "cmd"},
		strOut:   "docker run --name=my-ubuntu ubuntu cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker command is given via args and flags",
		command:  "run",
		argsIn:   []string{"--image=foo", "--command=cmdflag", "cmdarg"},
		strOut:   "docker run foo cmdarg",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from flag",
		command:  "run",
		argsIn:   []string{"--image=foo", "--arg=a", "--arg=b"},
		strOut:   "docker run foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from args",
		command:  "run",
		argsIn:   []string{"--image=foo", "cmd", "a", "b"},
		strOut:   "docker run foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker args is given via args and flags",
		command:  "run",
		argsIn:   []string{"--image=foo", "--arg=c", "--arg=d", "cmd", "a", "b"},
		strOut:   "docker run foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via args, args given via flags",
		command:  "run",
		argsIn:   []string{"--image=foo", "--arg=a", "--arg=b", "cmd"},
		strOut:   "docker run foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via flag, args given via args",
		command:  "run",
		argsIn:   []string{"--image=foo", "--command=cmd", "a", "b"},
		strOut:   "docker run foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "args given via args and flags",
		command:  "run",
		argsIn:   []string{"--image=foo", "--arg=c", "--arg=d", "--arg=e", "a", "b"},
		strOut:   "docker run foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "args given via args and flags (with label)",
		command:  "run",
		argsIn:   []string{"--image=foo", "--source=file://examples/fugu.labels.yml", "--arg=c", "--arg=d", "--arg=e", "a", "b"},
		strOut:   "docker run --name=my-redis foo a b",
		errOut:   nil,
	}).Test(t)

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	(&DockerCommandTest{
		testDesc: "command with volume flag",
		command:  "run",
		argsIn:   []string{"--image=foo", "--source=file://examples/fugu.volumes.yml"},
		strOut:   "docker run --name=my-ubuntu --volume=/tmp --volume=" + usr.HomeDir + "/Go:/root/Go foo",
		errOut:   nil,
	}).Test(t)
}

func TestCommandExec(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "exec",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingName,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain exec",
		command:  "exec",
		argsIn:   []string{"--name=foo"},
		strOut:   "docker exec foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain exec with label",
		command:  "exec",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker exec my-redis",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain exec with unknown label",
		command:  "exec",
		argsIn:   []string{"label-unknown", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker exec my-redis label-unknown",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from flag",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--command=cmd"},
		strOut:   "docker exec foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args",
		command:  "exec",
		argsIn:   []string{"--name=foo", "cmd"},
		strOut:   "docker exec foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker command is given via args and flags",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--command=cmdflag", "cmdarg"},
		strOut:   "docker exec foo cmdarg",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from flag",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--arg=a", "--arg=b"},
		strOut:   "docker exec foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from args",
		command:  "exec",
		argsIn:   []string{"--name=foo", "cmd", "a", "b"},
		strOut:   "docker exec foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker args is given via args and flags",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--arg=c", "--arg=d", "cmd", "a", "b"},
		strOut:   "docker exec foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via args, args given via flags",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--arg=a", "--arg=b", "cmd"},
		strOut:   "docker exec foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via flag, args given via args",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--command=cmd", "a", "b"},
		strOut:   "docker exec foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "args given via args and flags",
		command:  "exec",
		argsIn:   []string{"--name=foo", "--arg=c", "--arg=d", "--arg=e", "a", "b"},
		strOut:   "docker exec foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "args given via args and flags (with label)",
		command:  "exec",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml", "--arg=c", "--arg=d", "--arg=e", "a", "b"},
		strOut:   "docker exec my-redis a b",
		errOut:   nil,
	}).Test(t)
}

func TestCommandShell(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "shell",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingName,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain shell",
		command:  "shell",
		argsIn:   []string{"--name=foo"},
		strOut:   "docker exec --detach=false --interactive --tty foo /bin/bash",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "override shell",
		command:  "shell",
		argsIn:   []string{"--name=foo", "--shell=/bin/sh"},
		strOut:   "docker exec --detach=false --interactive --tty foo /bin/sh",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain shell with label",
		command:  "shell",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker exec --detach=false --interactive --tty my-redis /bin/bash",
		errOut:   nil,
	}).Test(t)
}

func TestCommandDestroy(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "destroy",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingName,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain destroy",
		command:  "destroy",
		argsIn:   []string{"--name=foo"},
		strOut:   "docker rm -f foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain destroy with label",
		command:  "destroy",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker rm -f my-redis",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain destroy with unknown label",
		command:  "destroy",
		argsIn:   []string{"label-unknown", "--source=file://examples/fugu.labels.yml"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "destroy: invalid number of args",
		command:  "destroy",
		argsIn:   []string{"--name=foo", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)
}

func TestCommandPush(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "push",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain push",
		command:  "push",
		argsIn:   []string{"--image=foo"},
		strOut:   "docker push foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain push with label",
		command:  "push",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker push redis",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain push with unknown label",
		command:  "push",
		argsIn:   []string{"label-unknown", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker push redis:label-unknown",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "push with tag",
		command:  "push",
		argsIn:   []string{"--image=foo", "--tag=bar"},
		strOut:   "docker push foo:bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "push with tag when given from args",
		command:  "push",
		argsIn:   []string{"--image=foo", "--tag=bar", "rab"},
		strOut:   "docker push foo:rab",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "push: invalid number of args",
		command:  "push",
		argsIn:   []string{"--image=foo", "tag", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "push: invalid number of args with label",
		command:  "push",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml", "tag", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)
}

func TestCommandPull(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "pull",
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain pull",
		command:  "pull",
		argsIn:   []string{"--image=foo"},
		strOut:   "docker pull foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain pull with label",
		command:  "pull",
		argsIn:   []string{"label1", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker pull redis",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain pull with unknown label",
		command:  "pull",
		argsIn:   []string{"label-unknown", "--source=file://examples/fugu.labels.yml"},
		strOut:   "docker pull redis:label-unknown",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "pull with tag",
		command:  "pull",
		argsIn:   []string{"--image=foo", "--tag=bar"},
		strOut:   "docker pull foo:bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "pull with tag when given from args",
		command:  "pull",
		argsIn:   []string{"--image=foo", "--tag=bar", "rab"},
		strOut:   "docker pull foo:rab",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "pull: invalid number of args",
		command:  "pull",
		argsIn:   []string{"--image=foo", "tag", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)
}

func TestShowData(t *testing.T) {
	(&CommandTest{
		testDesc:       "plain show-data call",
		command:        "show-data",
		argsIn:         []string{"label1", "--source=file://examples/fugu.labels.yml"},
		errOut:         nil,
		stdoutContains: []string{"image: redis", "name: my-redis"},
	}).Test(t)

	(&CommandTest{
		testDesc:       "show-data: invalid number of args",
		command:        "show-data",
		argsIn:         []string{"label1", "bogus", "--source=file://examples/fugu.labels.yml"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "show-data no output",
		command:        "show-data",
		argsIn:         []string{"label3", "--source=file://examples/fugu.labels.yml"},
		errOut:         nil,
		stdoutContains: []string{},
	}).Test(t)
}

func TestShowLabels(t *testing.T) {
	(&CommandTest{
		testDesc:       "plain show-labels call",
		command:        "show-labels",
		argsIn:         []string{"--source=file://examples/fugu.labels.yml"},
		errOut:         nil,
		stdoutContains: []string{"label1", "label2"},
	}).Test(t)

	(&CommandTest{
		testDesc:       "show-labels: invalid number of args",
		command:        "show-labels",
		argsIn:         []string{"bogus", "--source=file://examples/fugu.labels.yml"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "show-labels no output",
		command:        "show-labels",
		argsIn:         []string{""},
		errOut:         nil,
		stdoutContains: []string{},
	}).Test(t)
}

func TestListImages(t *testing.T) {
	(&CommandTest{
		testDesc:       "plain images call",
		command:        "images",
		argsIn:         []string{""},
		errOut:         nil,
		stdoutContains: []string{"<local docker images shown>"},
	}).Test(t)

	(&CommandTest{
		testDesc:       "images: invalid number of args",
		command:        "images",
		argsIn:         []string{"registry", "bogus"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "cannot use user when file given",
		command:        "images",
		argsIn:         []string{"registry", "--user=foobar", "--file=file"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "cannot use password when file given",
		command:        "images",
		argsIn:         []string{"registry", "--password=foobar", "--file=file"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "cannot use user or password when file given",
		command:        "images",
		argsIn:         []string{"registry", "--user=foobar", "--password=foobar", "--file=file"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "cannot use password-stdin when file given",
		command:        "images",
		argsIn:         []string{"registry", "--password-stdin", "--file=file"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "need user, too when password given",
		command:        "images",
		argsIn:         []string{"registry", "--password=foobar"},
		errOut:         ErrMissingFlag,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "need user, too when password-stin given",
		command:        "images",
		argsIn:         []string{"registry", "--password-stdin"},
		errOut:         ErrMissingFlag,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "user and password empty",
		command:        "images",
		argsIn:         []string{"registry", "--password=''", "--user=''"},
		errOut:         ErrNoCredentials,
		stdoutContains: []string{},
	}).Test(t)
}
