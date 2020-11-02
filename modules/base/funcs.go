package base

import (
	"github.com/zincium/zinc/modules/errgroup"
)

// https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/
// https://golangcode.com/waiting-for-goroutines-to-finish/
// https://medium.com/swlh/using-goroutines-and-wait-groups-for-concurrency-in-golang-78ca7a069d28
// https://scene-si.org/2020/05/29/waiting-on-goroutines/ thanks

// GroupExecute funcs
func GroupExecute(funcs ...func() error) error {
	var g errgroup.Group
	for _, f := range funcs {
		g.Go(f)
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
