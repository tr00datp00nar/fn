/*
A collection of functions I got tired of writing over and over again.
Currently not organized in any kind of way. Eventually, I may move
related functions into subdirectories under tr00datp00nar/fn.

I have made an attempt to document these functions in-line so that
go-doc keybinds can be used for function documentation.
*/

package fn

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//
// http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func OpenURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}

// RawUrlEncode takes any string and conactenates each word with "%20"
// as the delimeter. Used when encoding a URL.
// Often used in conjunction with Z.ArgsorIn()
//
//	More information available at: https://github.com/rwxrob/bonzai/z
func RawUrlEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// GetViperConfig reads a Viper Configuration into your application
//
// n is the name of your configuration file without the extension
//
//	for a configuration file `config.yaml`, use `config`
//
// t is the filetype of the configuration file
//
//	can be one of yaml, YAML, json, toml, ini, hcl
//	viper also supports envfile and java properties files
//
// p is the path to the configuration directory
//
//	usually `$XDG_CONFIG_HOME/some/path`
func GetViperConfig(n, t, p string) {
	viper.SetConfigName(n)
	viper.SetConfigType(t)
	viper.AddConfigPath(p)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// NewCmd runs the provided command using `exec.Command()`
func NewCmd(command string, args ...string) *cmd {
	return &cmd{exec.Command(command, args...)}
}

// ExecBash runs the provided command using `exec.Command("bash", "-c", cmd)`
func ExecBash(cmd string) (o string, e string, err error) {
	return NewCmd("bash", "-c", cmd).Output()
}

type buf []byte

func (b *buf) Read(p []byte) (n int, err error) {
	n = copy(p, *b)
	*b = (*b)[n:]
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (b *buf) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *buf) String() string { return string(*b) }

type cmd struct {
	cmd *exec.Cmd
}

func (c *cmd) Output() (stdout, stderr string, err error) {
	o, e := &buf{}, &buf{}
	c.cmd.Stdout = o
	c.cmd.Stderr = e
	err = c.cmd.Run()
	return o.String(), e.String(), err
}
