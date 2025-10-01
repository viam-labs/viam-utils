package viamutils

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	generic "go.viam.com/rdk/services/generic"
)

var (
	Arm              = resource.NewModel("viam", "viam-utils", "arm")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterService(generic.API, Arm,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newViamUtilsArm,
		},
	)
}

type Config struct {
	Arm string `json:"arm"`
}

func (cfg *Config) Validate(path string) ([]string, []string, error) {
	var deps []string
	if cfg.Arm == "" {
		return nil, nil, resource.NewConfigValidationFieldRequiredError(path, "arm")
	}
	deps = append(deps, cfg.Arm)

	return deps, nil, nil
}

type viamUtilsArm struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	arm      arm.Arm
	armModel referenceframe.Model

	mu sync.Mutex
}

func newViamUtilsArm(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewArm(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewArm(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &viamUtilsArm{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}

	var err error

	s.arm, err = arm.FromDependencies(deps, conf.Arm)
	if err != nil {
		return nil, err
	}

	s.armModel, err = s.arm.Kinematics(ctx)
	if err != nil {
		return nil, err
	}

	s.cfg = conf

	return s, nil
}

func (s *viamUtilsArm) Name() resource.Name {
	return s.name
}

func (s *viamUtilsArm) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resp := map[string]any{}

	for key, value := range cmd {
		switch key {
		case "transform":
			values, ok := value.([]any)
			if !ok {
				return nil, fmt.Errorf("input must be an array")
			}

			var jointPos []referenceframe.Input
			for i, val := range values {
				floatVal, ok := val.(float64)
				if !ok {
					return nil, fmt.Errorf("joint position at index %d must be a number", i)
				}
				jointPos = append(jointPos, referenceframe.Input{Value: floatVal})
			}

			pose, err := s.armModel.Transform(jointPos)
			if err != nil {
				s.logger.Error(err)
				return nil, err
			}
			return map[string]any{
				"orientation": pose.Orientation(),
			}, nil
		default:
			return resp, nil
		}
	}

	if len(resp) == 0 {
		return nil, errors.New("no valid DoCommand submitted")
	}
	return resp, nil
}

func (s *viamUtilsArm) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
