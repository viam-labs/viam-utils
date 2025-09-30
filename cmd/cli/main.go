package main

import (
	"context"
	"viamutils"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	arm "go.viam.com/rdk/components/arm"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	ctx := context.Background()
	logger := logging.NewLogger("cli")

	deps := resource.Dependencies{}
	// can load these from a remote machine if you need

	cfg := viamutils.Config{}

	thing, err := viamutils.NewArm(ctx, deps, arm.Named("foo"), &cfg, logger)
	if err != nil {
		return err
	}
	defer thing.Close(ctx)

	return nil
}
