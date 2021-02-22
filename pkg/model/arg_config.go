package model

import "fmt"

var version *string

// ArgConfig is the settings for passing arguments to the command
type ArgConfig struct {
	SourceHost         string `arg:"--source-host,required" help:"source redis server hostname"`
	SourcePort         int    `arg:"--source-port" default:"6379" help:"source redis server port"`
	SourceDBIndex      int    `arg:"--source-db" default:"0" help:"source redis server db index"`
	SourceUser         string `arg:"--source-user" help:"source redis server auth username"`
	SourcePass         string `arg:"--source-pass" help:"source redis server auth password"`
	DestinationHost    string `arg:"--destination-host,required" help:"destination redis server hostname"`
	DestinationPort    int    `arg:"--destination-port" default:"6379" help:"destination redis server port"`
	DestinationDBIndex int    `arg:"--destination-db" default:"0" help:"destination redis server db index"`
	DestinationUser    string `arg:"--destination-user" help:"destination redis server auth username"`
	DestinationPass    string `arg:"--destination-pass" help:"destination redis server auth password"`
	KeyFilter          string `arg:"--source-filter" default:"*" help:"source keys filter string"`
	KeyPrefix          string `arg:"--destination-prefix" help:"destination key prefix to prepend"`
	Verbose            bool   `arg:"-v" help:"Print a lot more info"`
}

// SetVersion ...
func SetVersion(ver *string) {
	version = ver
}

// Version ...
func (ArgConfig) Version() string {
	return fmt.Sprintf("v%s", *version)
}

// Description returns a simple description of the command
func (ArgConfig) Description() string {
	return `Redis Migrator will take the keys from one server/db and
transfer them to another server/db. If connection details
are not passed in as arguments, you will be asked for
them interactively.
`
}

// Print will show the configuration given
func (c *ArgConfig) Print() {
	fmt.Printf(`SOURCE:
  Server:   %s:%d
  DB Index: %d
  Auth:     %s
DESTINATION:
  Server:   %s:%d
  DB Index: %d
  Auth:     %s
`,
		c.SourceHost,
		c.SourcePort,
		c.SourceDBIndex,
		determineAuth(c.SourceUser, c.SourcePass),
		c.DestinationHost,
		c.DestinationPort,
		c.DestinationDBIndex,
		determineAuth(c.DestinationUser, c.DestinationPass),
	)
}

// GetSource returns source redis server
func (c *ArgConfig) GetSource() RedisServerConfig {
	return RedisServerConfig{
		Hostname: c.SourceHost,
		Port:     c.SourcePort,
		DBIndex:  c.SourceDBIndex,
		Username: c.SourceUser,
		Password: c.SourcePass,
	}
}

// GetDestination returns destination redis server
func (c *ArgConfig) GetDestination() RedisServerConfig {
	return RedisServerConfig{
		Hostname: c.DestinationHost,
		Port:     c.DestinationPort,
		DBIndex:  c.DestinationDBIndex,
		Username: c.DestinationUser,
		Password: c.DestinationPass,
	}
}

// IsIdenticalServers returns if the source and destination are the same
func (c *ArgConfig) IsIdenticalServers() bool {
	return c.SourceHost == c.DestinationHost && c.SourcePort == c.DestinationPort && c.SourceDBIndex == c.DestinationDBIndex
}

func determineAuth(user, pass string) string {
	auth := "None"
	if pass != "" {
		auth = "Password only"
	}
	if user != "" {
		auth = "Username and Password"
	}
	return auth
}
