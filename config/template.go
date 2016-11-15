package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	// DefaultTemplateFilePerms are the default file permissions for templates
	// rendered onto disk when a specific file permission has not already been
	// specified.
	DefaultTemplateFilePerms = 0644

	// DefaultTemplateCommandTimeout is the amount of time to wait for a command
	// to return.
	DefaultTemplateCommandTimeout = 30 * time.Second
)

var (
	ErrTemplateStringEmpty   = errors.New("template: cannot be empty")
	ErrTemplateInvalidFormat = errors.New("template: invalid format")

	// configTemplateRe is the pattern to split the config template syntax.
	configTemplateRe = regexp.MustCompile("([a-zA-Z]:)?([^:]+)")
)

type TemplateConfig struct {
	// Backup determines if this template should retain a backup. The default
	// value is false.
	Backup *bool `mapstructure:"backup"`

	// Command is the arbitrary command to execute after a template has
	// successfully rendered. This is DEPRECATED. Use Exec instead.
	Command *string `mapstructure:"command"`

	// CommandTimeout is the amount of time to wait for the command to finish
	// before force-killing it. This is DEPRECATED. Use Exec instead.
	CommandTimeout *time.Duration `mapstructure:"command_timeout"`

	// Contents are the raw template contents to evaluate. Either this or Source
	// must be specified, but not both.
	Contents *string `mapstructure:"contents"`

	// Destination is the location on disk where the template should be rendered.
	// This is required unless running in debug/dry mode.
	Destination *string `mapstructure:"destination"`

	// Exec is the configuration for the command to run when the template renders
	// successfully.
	Exec *ExecConfig `mapstructure:"exec"`

	// Perms are the file system permissions to use when creating the file on
	// disk. This is useful for when files contain sensitive information, such as
	// secrets from Vault.
	Perms *os.FileMode `mapstructure:"perms"`

	// Source is the path on disk to the template contents to evaluate. Either
	// this or Contents should be specified, but not both.
	Source *string `mapstructure:"source"`

	// Wait configures per-template quiescence timers.
	Wait *WaitConfig `mapstructure:"wait"`

	// LeftDelim and RightDelim are optional configurations to control what
	// delimiter is utilized when parsing the template.
	LeftDelim  *string `mapstructure:"left_delimiter"`
	RightDelim *string `mapstructure:"right_delimiter"`
}

// DefaultTemplateConfig returns a configuration that is populated with the
// default values.
func DefaultTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		Exec: DefaultExecConfig(),
		Wait: DefaultWaitConfig(),
	}
}

// Copy returns a deep copy of this configuration.
func (c *TemplateConfig) Copy() *TemplateConfig {
	if c == nil {
		return nil
	}

	var o TemplateConfig

	o.Backup = c.Backup

	o.Command = c.Command

	o.CommandTimeout = c.CommandTimeout

	o.Contents = c.Contents

	o.Destination = c.Destination

	if c.Exec != nil {
		o.Exec = c.Exec.Copy()
	}

	o.Perms = c.Perms

	o.Source = c.Source

	if c.Wait != nil {
		o.Wait = c.Wait.Copy()
	}

	o.LeftDelim = c.LeftDelim
	o.RightDelim = c.RightDelim

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *TemplateConfig) Merge(o *TemplateConfig) *TemplateConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Backup != nil {
		r.Backup = o.Backup
	}

	if o.Command != nil {
		r.Command = o.Command
	}

	if o.CommandTimeout != nil {
		r.CommandTimeout = o.CommandTimeout
	}

	if o.Contents != nil {
		r.Contents = o.Contents
	}

	if o.Destination != nil {
		r.Destination = o.Destination
	}

	if o.Exec != nil {
		r.Exec = r.Exec.Merge(o.Exec)
	}

	if o.Perms != nil {
		r.Perms = o.Perms
	}

	if o.Source != nil {
		r.Source = o.Source
	}

	if o.Wait != nil {
		r.Wait = r.Wait.Merge(o.Wait)
	}

	if o.LeftDelim != nil {
		r.LeftDelim = o.LeftDelim
	}

	if o.RightDelim != nil {
		r.RightDelim = o.RightDelim
	}

	return r
}

// Finalize ensures the configuration has no nil pointers and sets default
// values.
func (c *TemplateConfig) Finalize() {
	if c.Backup == nil {
		c.Backup = Bool(false)
	}

	if c.Command == nil {
		c.Command = String("")
	}

	if c.CommandTimeout == nil {
		c.CommandTimeout = TimeDuration(DefaultTemplateCommandTimeout)
	}

	if c.Contents == nil {
		c.Contents = String("")
	}

	if c.Destination == nil {
		c.Destination = String("")
	}

	if c.Exec == nil {
		c.Exec = DefaultExecConfig()
	}

	// Backwards compat for specifying command directly
	if c.Exec.Command == nil && c.Command != nil {
		c.Exec.Command = c.Command
	}
	if c.Exec.Timeout == nil && c.CommandTimeout != nil {
		c.Exec.Timeout = c.CommandTimeout
	}
	c.Exec.Finalize()

	if c.Perms == nil {
		c.Perms = FileMode(DefaultTemplateFilePerms)
	}

	if c.Source == nil {
		c.Source = String("")
	}

	if c.Wait == nil {
		c.Wait = DefaultWaitConfig()
	}
	c.Wait.Finalize()

	if c.LeftDelim == nil {
		c.LeftDelim = String("")
	}

	if c.RightDelim == nil {
		c.RightDelim = String("")
	}
}

// GoString defines the printable version of this struct.
func (c *TemplateConfig) GoString() string {
	if c == nil {
		return "(*TemplateConfig)(nil)"
	}

	return fmt.Sprintf("&TemplateConfig{"+
		"Backup:%s, "+
		"Command:%s, "+
		"CommandTimeout:%s, "+
		"Contents:%s, "+
		"Destination:%s, "+
		"Exec:%#v, "+
		"Perms:%s, "+
		"Source:%s, "+
		"Wait:%#v, "+
		"LeftDelim:%s, "+
		"RightDelim:%s"+
		"}",
		BoolGoString(c.Backup),
		StringGoString(c.Command),
		TimeDurationGoString(c.CommandTimeout),
		StringGoString(c.Contents),
		StringGoString(c.Destination),
		c.Exec,
		FileModeGoString(c.Perms),
		StringGoString(c.Source),
		c.Wait,
		StringGoString(c.LeftDelim),
		StringGoString(c.RightDelim),
	)
}

// TemplateConfigs is a collection of TemplateConfigs
type TemplateConfigs []*TemplateConfig

// DefaultTemplateConfigs returns a configuration that is populated with the
// default values.
func DefaultTemplateConfigs() *TemplateConfigs {
	return &TemplateConfigs{}
}

// Copy returns a deep copy of this configuration.
func (c *TemplateConfigs) Copy() *TemplateConfigs {
	o := make(TemplateConfigs, len(*c))
	for i, t := range *c {
		o[i] = t.Copy()
	}
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *TemplateConfigs) Merge(o *TemplateConfigs) *TemplateConfigs {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	*r = append(*r, *o...)

	return r
}

// Finalize ensures the configuration has no nil pointers and sets default
// values.
func (c *TemplateConfigs) Finalize() {
	if c == nil {
		*c = *DefaultTemplateConfigs()
	}

	for _, t := range *c {
		t.Finalize()
	}
}

// GoString defines the printable version of this struct.
func (c *TemplateConfigs) GoString() string {
	if c == nil {
		return "(*TemplateConfigs)(nil)"
	}

	s := make([]string, len(*c))
	for i, t := range *c {
		s[i] = t.GoString()
	}

	return "{" + strings.Join(s, ", ") + "}"
}

// ParseTemplateConfig parses a string in the form source:destination:command
// into a TemplateConfig.
func ParseTemplateConfig(s string) (*TemplateConfig, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, ErrTemplateStringEmpty
	}

	var source, destination, command string
	parts := configTemplateRe.FindAllString(s, -1)

	switch len(parts) {
	case 1:
		source = parts[0]
	case 2:
		source, destination = parts[0], parts[1]
	case 3:
		source, destination, command = parts[0], parts[1], parts[2]
	default:
		return nil, ErrTemplateInvalidFormat
	}

	var sourcePtr, destinationPtr, commandPtr *string
	if source != "" {
		sourcePtr = String(source)
	}
	if destination != "" {
		destinationPtr = String(destination)
	}
	if command != "" {
		commandPtr = String(command)
	}

	return &TemplateConfig{
		Source:      sourcePtr,
		Destination: destinationPtr,
		Command:     commandPtr,
	}, nil
}
