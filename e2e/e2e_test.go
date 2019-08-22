package e2e_test

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("cpu-entitlement-admin-plugin", func() {
	var (
		org string
		uid string
	)

	BeforeEach(func() {
		uid = uuid.New().String()
		org = "org-" + uid
		space := "space-" + uid

		Expect(Cmd("cf", "create-org", org).Run()).To(gexec.Exit(0))
		Expect(Cmd("cf", "target", "-o", org).Run()).To(gexec.Exit(0))
		Expect(Cmd("cf", "create-space", space).Run()).To(gexec.Exit(0))
		Expect(Cmd("cf", "target", "-o", org, "-s", space).Run()).To(gexec.Exit(0))
	})

	AfterEach(func() {
		Expect(Cmd("cf", "delete-org", "-f", org).WithTimeout("1m").Run()).To(gexec.Exit(0))
	})

	Describe("with an app", func() {
		var (
			appName string
			appURL  string
		)

		BeforeEach(func() {
			appName = "spinner-" + uid

			Expect(Cmd("cf", "push", appName).WithDir("../../spinner").WithTimeout("2m").Run()).To(gexec.Exit(0))

			appURL = strings.Replace(cfApi, "api.", appName+".", 1)
		})

		It("prints the list of apps that are over entitlement", func() {
			httpGet(appURL + "/spin")
			Eventually(Cmd("cf", "over-entitlement-instances").Run, "20s", "1s").Should(gbytes.Say(appName))
		})
	})

	It("prints a no apps over messages if no apps over entitlement", func() {
		Consistently(Cmd("cf", "over-entitlement-instances").Run).Should(gbytes.Say("No apps over entitlement"))
	})

})

type Command struct {
	cmd     string
	args    []string
	dir     string
	timeout string
}

func Cmd(cmd string, args ...string) Command {
	return Command{
		cmd:     cmd,
		args:    args,
		timeout: "1s",
	}
}

func (c Command) WithDir(dir string) Command {
	return Command{
		cmd:     c.cmd,
		args:    c.args,
		dir:     dir,
		timeout: c.timeout,
	}
}

func (c Command) WithTimeout(timeout string) Command {
	return Command{
		cmd:     c.cmd,
		args:    c.args,
		dir:     c.dir,
		timeout: timeout,
	}
}

func (c Command) Run() *gexec.Session {
	session, err := gexec.Start(c.build(), GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	EventuallyWithOffset(1, session, c.timeout).Should(gexec.Exit())
	return session
}

func (c Command) build() *exec.Cmd {
	command := exec.Command(c.cmd, c.args...)
	if c.dir != "" {
		cwd, err := os.Getwd()
		ExpectWithOffset(2, err).NotTo(HaveOccurred())
		command.Dir = filepath.Join(cwd, c.dir)
	}
	return command
}

func httpGet(url string) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(url)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, resp.StatusCode).To(Equal(http.StatusOK))
}
