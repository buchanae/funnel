package e2e

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/tests/testutils"
)

var log = logger.New("e2e")
var fun = testutils.NewFunnel()

func init() {
	fun.StartServer()
}
