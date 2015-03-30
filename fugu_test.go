package fugu

import (
	// "bytes"
	"github.com/mattes/go-collect"
	"github.com/mattes/go-collect/data"
	"github.com/stretchr/testify/assert"
	// "io"
	"io/ioutil"
	"os"
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
	dataIn   *data.Data
	argsIn   []string
	strOut   string
	errOut   error
}

func (dct *DockerCommandTest) Test(t *testing.T) {
	if RunOnly != "" && RunOnly != dct.testDesc {
		return
	}
	str, err := DockerCommands[dct.command](dct.dataIn, dct.argsIn)
	assert.Equal(t, dct.errOut, err, dct.testDesc)
	if err == nil {
		assert.Equal(t, dct.strOut, str, dct.testDesc)
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

func TestCommandBuild(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "image is missing",
		command:  "build",
		dataIn:   data.New(),
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "tag flag",
		command:  "build",
		dataIn:   data.New().Set("image", "foo").Set("tag", "bar"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo:bar .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "no tag given",
		command:  "build",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "given path",
		command:  "build",
		dataIn:   data.New().Set("image", "foo").Set("path", "bar"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "given url",
		command:  "build",
		dataIn:   data.New().Set("image", "foo").Set("url", "bar"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "both given path and url",
		command:  "build",
		dataIn:   data.New().Set("image", "foo").Set("url", "url").Set("path", "path"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo url",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "path or url given via args",
		command:  "build",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{"pathOrUrl"},
		strOut:   "docker build --tag=foo pathOrUrl",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "invalid number of args",
		command:  "build",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{"pathOrUrl", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "tag-git-branch",
		command:  "build",
		dataIn:   data.New().Set("image", "foo").SetTrue("tag-git-branch"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo:current-branch .",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "tag-git-branch overwrites given tag flag",
		command:  "build",
		dataIn:   data.New().Set("image", "foo").Set("tag", "bar").SetTrue("tag-git-branch"),
		argsIn:   []string{},
		strOut:   "docker build --tag=foo:current-branch .",
		errOut:   nil,
	}).Test(t)
}

func TestCommandRun(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "image is missing",
		command:  "run",
		dataIn:   data.New(),
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain run",
		command:  "run",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{},
		strOut:   "docker run foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from flag",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Set("command", "cmd"),
		argsIn:   []string{},
		strOut:   "docker run foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args",
		command:  "run",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{"cmd"},
		strOut:   "docker run foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker command is given via args and flags",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Set("command", "cmdflag"),
		argsIn:   []string{"cmdarg"},
		strOut:   "docker run foo cmdarg",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from flag",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Add("arg", "a").Add("arg", "b"),
		argsIn:   []string{},
		strOut:   "docker run foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from args",
		command:  "run",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{"cmd", "a", "b"},
		strOut:   "docker run foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker args is given via args and flags",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Add("arg", "c").Add("arg", "d"),
		argsIn:   []string{"cmd", "a", "b"},
		strOut:   "docker run foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via args, args given via flags",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Add("arg", "a").Add("arg", "b"),
		argsIn:   []string{"cmd"},
		strOut:   "docker run foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via flag, args given via args",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Set("command", "cmd"),
		argsIn:   []string{"a", "b"},
		strOut:   "docker run foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "args given via args and flags",
		command:  "run",
		dataIn:   data.New().Set("image", "foo").Add("arg", "c").Add("arg", "d").Add("arg", "e"),
		argsIn:   []string{"a", "b"},
		strOut:   "docker run foo a b",
		errOut:   nil,
	}).Test(t)
}

func TestCommandExec(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "exec",
		dataIn:   data.New(),
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingName,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain exec",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo"),
		argsIn:   []string{},
		strOut:   "docker exec --interactive --tty foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain exec",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo"),
		argsIn:   []string{},
		strOut:   "docker exec --interactive --tty foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "allow overwrite interactive",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").SetFalse("interactive"),
		argsIn:   []string{},
		strOut:   "docker exec --tty foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "allow overwrite tty",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").SetFalse("tty"),
		argsIn:   []string{},
		strOut:   "docker exec --interactive foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from flag",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Set("command", "cmd"),
		argsIn:   []string{},
		strOut:   "docker exec --interactive --tty foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker command from args",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo"),
		argsIn:   []string{"cmd"},
		strOut:   "docker exec --interactive --tty foo cmd",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker command is given via args and flags",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Set("command", "cmdflag"),
		argsIn:   []string{"cmdarg"},
		strOut:   "docker exec --interactive --tty foo cmdarg",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from flag",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Add("arg", "a").Add("arg", "b"),
		argsIn:   []string{},
		strOut:   "docker exec --interactive --tty foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "get docker args from args",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo"),
		argsIn:   []string{"cmd", "a", "b"},
		strOut:   "docker exec --interactive --tty foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "docker args is given via args and flags",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Add("arg", "c").Add("arg", "d"),
		argsIn:   []string{"cmd", "a", "b"},
		strOut:   "docker exec --interactive --tty foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via args, args given via flags",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Add("arg", "a").Add("arg", "b"),
		argsIn:   []string{"cmd"},
		strOut:   "docker exec --interactive --tty foo cmd a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "command given via flag, args given via args",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Set("command", "cmd"),
		argsIn:   []string{"a", "b"},
		strOut:   "docker exec --interactive --tty foo a b",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "args given via args and flags",
		command:  "exec",
		dataIn:   data.New().Set("name", "foo").Add("arg", "c").Add("arg", "d").Add("arg", "e"),
		argsIn:   []string{"a", "b"},
		strOut:   "docker exec --interactive --tty foo a b",
		errOut:   nil,
	}).Test(t)
}

func TestCommandDestroy(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "destroy",
		dataIn:   data.New(),
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingName,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain destroy",
		command:  "destroy",
		dataIn:   data.New().Set("name", "foo"),
		argsIn:   []string{},
		strOut:   "docker rm -f foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "invalid number of args",
		command:  "destroy",
		dataIn:   data.New().Set("name", "foo"),
		argsIn:   []string{"bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)
}

func TestCommandPush(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "push",
		dataIn:   data.New(),
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain push",
		command:  "push",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{},
		strOut:   "docker push foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "push with tag",
		command:  "push",
		dataIn:   data.New().Set("image", "foo").Set("tag", "bar"),
		argsIn:   []string{},
		strOut:   "docker push foo:bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "push with tag when given from args",
		command:  "push",
		dataIn:   data.New().Set("image", "foo").Set("tag", "bar"),
		argsIn:   []string{"rab"},
		strOut:   "docker push foo:rab",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "invalid number of args",
		command:  "push",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{"tag", "bogus"},
		strOut:   "",
		errOut:   ErrTooManyArgs,
	}).Test(t)
}

func TestCommandPull(t *testing.T) {
	(&DockerCommandTest{
		testDesc: "name is missing",
		command:  "pull",
		dataIn:   data.New(),
		argsIn:   []string{},
		strOut:   "",
		errOut:   ErrMissingImage,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "plain pull",
		command:  "pull",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{},
		strOut:   "docker pull foo",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "pull with tag",
		command:  "pull",
		dataIn:   data.New().Set("image", "foo").Set("tag", "bar"),
		argsIn:   []string{},
		strOut:   "docker pull foo:bar",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "pull with tag when given from args",
		command:  "pull",
		dataIn:   data.New().Set("image", "foo").Set("tag", "bar"),
		argsIn:   []string{"rab"},
		strOut:   "docker pull foo:rab",
		errOut:   nil,
	}).Test(t)

	(&DockerCommandTest{
		testDesc: "invalid number of args",
		command:  "pull",
		dataIn:   data.New().Set("image", "foo"),
		argsIn:   []string{"tag", "bogus"},
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
		testDesc:       "invalid number of args",
		command:        "show-data",
		argsIn:         []string{"label1", "bogus", "--source=file://examples/fugu.labels.yml"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "no output",
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
		testDesc:       "invalid number of args",
		command:        "show-labels",
		argsIn:         []string{"bogus", "--source=file://examples/fugu.labels.yml"},
		errOut:         ErrTooManyArgs,
		stdoutContains: []string{},
	}).Test(t)

	(&CommandTest{
		testDesc:       "no output",
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
		testDesc:       "invalid number of args",
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
