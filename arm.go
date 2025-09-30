package viamutils

import (
	"context"
	"errors"
	"fmt"
	"sync"

	arm "go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

var (
	Arm              = resource.NewModel("viam", "viam-utils", "arm")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(arm.API, Arm,
		resource.Registration[arm.Arm, *Config]{
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
	resource.TriviallyCloseable

	name resource.Name

	logger logging.Logger
	cfg    *Config

	arm      arm.Arm
	armModel referenceframe.Model

	mu sync.Mutex
}

func newViamUtilsArm(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (arm.Arm, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewArm(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewArm(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (arm.Arm, error) {

	s := &viamUtilsArm{
		name:   name,
		logger: logger,
		cfg:    conf,
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

// EndPosition returns the current position of the arm.
func (s *viamUtilsArm) EndPosition(ctx context.Context, extra map[string]any) (spatialmath.Pose, error) {
	var poseRetVal spatialmath.Pose

	return poseRetVal, fmt.Errorf("not implemented")
}

// MoveToPosition moves the arm to the given absolute position.
// This will block until done or a new operation cancels this one.
func (s *viamUtilsArm) MoveToPosition(ctx context.Context, pose spatialmath.Pose, extra map[string]any) error {
	return fmt.Errorf("not implemented")
}

// MoveToJointPositions moves the arm's joints to the given positions.
// This will block until done or a new operation cancels this one.
func (s *viamUtilsArm) MoveToJointPositions(ctx context.Context, positions []referenceframe.Input, extra map[string]any) error {
	return fmt.Errorf("not implemented")
}

// MoveThroughJointPositions moves the arm's joints through the given positions in the order they are specified.
// This will block until done or a new operation cancels this one.
func (s *viamUtilsArm) MoveThroughJointPositions(ctx context.Context, positions [][]referenceframe.Input, options *arm.MoveOptions, extra map[string]any) error {
	return fmt.Errorf("not implemented")
}

// JointPositions returns the current joint positions of the arm.
func (s *viamUtilsArm) JointPositions(ctx context.Context, extra map[string]any) ([]referenceframe.Input, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) Stop(ctx context.Context, extra map[string]any) error {
	return fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) Kinematics(ctx context.Context) (referenceframe.Model, error) {
	var modelRetVal referenceframe.Model

	return modelRetVal, fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) GoToInputs(ctx context.Context, inputSteps ...[]referenceframe.Input) error {
	return fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) DoCommand(ctx context.Context, cmd map[string]any) (map[string]any, error) {
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

			if len(values) != 6 {
				return nil, fmt.Errorf("input must contain exactly 6 joint positions, got %d", len(values))
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
			return map[string]any{"command": pose}, nil
		default:
			return resp, nil
		}
	}

	if len(resp) == 0 {
		return nil, errors.New("no valid DoCommand submitted")
	}
	return resp, nil
}

func (s *viamUtilsArm) IsMoving(ctx context.Context) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) Geometries(ctx context.Context, extra map[string]any) ([]spatialmath.Geometry, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *viamUtilsArm) Close(context.Context) error {
	// Put close code here
	// s.cancelFunc()
	return nil
}
