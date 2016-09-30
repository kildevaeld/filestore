package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/kildevaeld/filestore"
	_ "github.com/kildevaeld/filestore/filesystem"
	_ "github.com/kildevaeld/filestore/memory"
	_ "github.com/kildevaeld/filestore/s3"
	"github.com/spf13/viper"
)

func getFileStore() (filestore.Store, error) {

	var config filestore.Options
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.Driver == "" {
		if err := viper.UnmarshalKey("filestore", &config); err != nil {
			return nil, err
		}
	}

	if config.Driver == "" {
		return nil, errors.New("No driver selected")
	}

	return filestore.New(config)
}

func exit(code int, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	os.Exit(code)
}
