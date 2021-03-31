package cmd

import (
	"log"
	"time"

	"github.com/powerpuffpenguin/xormcache/cache"
	"github.com/powerpuffpenguin/xormcache/cmd/internal/daemon"
	"github.com/spf13/cobra"
	"xorm.io/xorm/caches"
)

func init() {
	var (
		debug                   bool
		addr, certFile, keyFile string
		maxAge                  int64
		maxElementSize          int
	)

	cmd := &cobra.Command{
		Use:   `daemon`,
		Short: `run as daemon`,
		Run: func(cmd *cobra.Command, args []string) {
			if maxAge < 1 {
				maxAge = 3600
			}
			if maxElementSize < 1 {
				maxElementSize = 1000
			}
			expired := time.Second * time.Duration(maxAge)
			log.Println(`expired =`, expired)
			log.Println(`max-size =`, maxElementSize)
			cacher := caches.NewLRUCacher2(
				caches.NewMemoryStore(),
				expired,
				maxElementSize,
			)
			cache.SetDefaultCacher(cacher)
			// run
			daemon.Run(addr, certFile, keyFile, debug)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, `addr`,
		`a`,
		`:1234`,
		`listen address`,
	)
	flags.StringVar(&certFile, `cert-file`,
		``,
		`tls certfile`,
	)
	flags.StringVar(&keyFile, `key-file`,
		``,
		`tls keyfile`,
	)
	flags.Int64Var(&maxAge, `max-age`,
		3600,
		`expired seconds`,
	)
	flags.IntVar(&maxElementSize, `max-size`,
		1000,
		`max element size`,
	)
	flags.BoolVarP(&debug, `debug`,
		`d`,
		false,
		`run as debug`,
	)
	rootCmd.AddCommand(cmd)
}
