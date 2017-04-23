package cafe

import (
	"log"
	"os"
	"os/exec"
	"testing"

	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/ginkgo"
	. "github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega"

	"github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega/gbytes"
	"github.com/duskhacker/cqrsnu/internal/github.com/onsi/gomega/gexec"
)

var (
	serverSession *gexec.Session
	suite         = "cafe"
)

func TestMainSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, suite+" Suite")
}

func RemoveDataFiles() {
	dirName := os.ExpandEnv("${GOPATH}/src/github.com/duskhacker/cqrsnu/data")
	dir, err := os.Open(dirName)
	if err != nil {
		log.Fatalf("error opening %s: %s", dirName, err)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatalf("error reading dir %s: %s\n", dir.Name(), err)
	}

	for _, file := range files {
		os.Remove(dir.Name() + "/" + file.Name())
	}
}

var _ = BeforeSuite(func() {
	var err error

	RemoveDataFiles()
	dataPath := os.ExpandEnv("${GOPATH}/src/github.com/duskhacker/cqrsnu/data")

	command := exec.Command("nsqd", "--data-path="+dataPath, "--tcp-address=localhost:4150", "--http-address=localhost:4151", "--broadcast-address=localhost")
	serverSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(serverSession.Err, "2s").Should(gbytes.Say(`TCP: listening on`))

	connectToNSQD = true
	SetNsqdTCPAddr("localhost:4150")
	InitConsumers()
})

var _ = AfterSuite(func() {
	StopAllConsumers()
	serverSession.Interrupt()
	gexec.CleanupBuildArtifacts()
})
