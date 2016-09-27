package plugin

import (
	"github.com/kmaehashi/sensorbee-traceback-formatter"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
)

func init() {
	udf.MustRegisterGlobalUDSCreator("traceback", udf.UDSCreatorFunc(traceback.NewTracebackUDS))
}
