package docker

import (
	"path/filepath"
	"testing"

	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
	shlex "github.com/flynn/go-shlex"
	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	exp := []string{"sh", "-c", "exec /opt/bin/flanneld -logtostderr=true -iface=${NODE_IP}"}
	cmd, err := shlex.Split("sh -c 'exec /opt/bin/flanneld -logtostderr=true -iface=${NODE_IP}'")
	assert.Nil(t, err)
	assert.Equal(t, exp, cmd)
}

func TestParseBindsAndVolumes(t *testing.T) {
	ctx := &Context{}
	ctx.ComposeFiles = []string{"foo/docker-compose.yml"}
	ctx.ResourceLookup = &lookup.FileConfigLookup{}

	abs, err := filepath.Abs(".")
	assert.Nil(t, err)
	cfg, hostCfg, err := Convert(&project.ServiceConfig{
		Volumes: []string{"/foo", "/home:/home", "/bar/baz", ".:/home", "/usr/lib:/usr/lib:ro"},
	}, ctx.Context)
	assert.Nil(t, err)
	assert.Equal(t, map[string]struct{}{"/foo": {}, "/bar/baz": {}}, cfg.Volumes)
	assert.Equal(t, []string{"/home:/home", abs + "/foo:/home", "/usr/lib:/usr/lib:ro"}, hostCfg.Binds)
}

func TestParseLabels(t *testing.T) {
	ctx := &Context{}
	ctx.ComposeFiles = []string{"foo/docker-compose.yml"}
	ctx.ResourceLookup = &lookup.FileConfigLookup{}
	bashCmd := "bash"
	fooLabel := "foo.label"
	fooLabelValue := "service.config.value"
	sc := &project.ServiceConfig{
		Entrypoint: project.NewCommand(bashCmd),
		Labels:     project.NewSliceorMap(map[string]string{fooLabel: "service.config.value"}),
	}
	cfg, _, err := Convert(sc, ctx.Context)
	assert.Nil(t, err)

	cfg.Labels[fooLabel] = "FUN"
	cfg.Entrypoint[0] = "less"

	assert.Equal(t, fooLabelValue, sc.Labels.MapParts()[fooLabel])
	assert.Equal(t, "FUN", cfg.Labels[fooLabel])

	assert.Equal(t, []string{bashCmd}, sc.Entrypoint.Slice())
	assert.Equal(t, []string{"less"}, []string(cfg.Entrypoint))
}
