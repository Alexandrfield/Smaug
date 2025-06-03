package worker

import (
	"fmt"
	"os"

	clientCommon "github.com/Alexandrfield/Smaug/internal/client/common"
)

type OpenData struct {
	description string
	data        []byte
	metadata    *clientCommon.Metadata
}

func (dat *OpenData) Show() string {
	var res string
	res += dat.description
	if dat.metadata.TypeNote == "file" {
		err := os.WriteFile(dat.metadata.Filename, dat.data, 0644)
		if err != nil {
			res += fmt.Sprintf("\n   WARN: nproblem with save file: %s; err:%s\n", dat.metadata.Filename, err)
		} else {
			res += fmt.Sprintf("\n  file %s saved.\n", dat.metadata.Filename)
		}
	}
	if dat.metadata.TypeNote == "text" {
		res += fmt.Sprintf("\n  text: %s\n", string(dat.data))
	}
	if dat.metadata.TypeNote == "card" {
		res += fmt.Sprintf("\n  card: %s\n", string(dat.data))
	}
	return res
}
