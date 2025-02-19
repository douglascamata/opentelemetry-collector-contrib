// Code generated by mdatagen. DO NOT EDIT.

package processscraper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/scraper"
	"go.opentelemetry.io/collector/scraper/scrapertest"
)

var typ = component.MustNewType("process")

func TestComponentFactoryType(t *testing.T) {
	require.Equal(t, typ, NewFactory().Type())
}

func TestComponentConfigStruct(t *testing.T) {
	require.NoError(t, componenttest.CheckConfigStruct(NewFactory().CreateDefaultConfig()))
}

func TestComponentLifecycle(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		createFn func(ctx context.Context, set scraper.Settings, cfg component.Config) (component.Component, error)
		name     string
	}{

		{
			name: "metrics",
			createFn: func(ctx context.Context, set scraper.Settings, cfg component.Config) (component.Component, error) {
				return factory.CreateMetrics(ctx, set, cfg)
			},
		},
	}

	cm, err := confmaptest.LoadConf("metadata.yaml")
	require.NoError(t, err)
	cfg := factory.CreateDefaultConfig()
	sub, err := cm.Sub("tests::config")
	require.NoError(t, err)
	require.NoError(t, sub.Unmarshal(&cfg))

	for _, tt := range tests {
		t.Run(tt.name+"-shutdown", func(t *testing.T) {
			c, err := tt.createFn(context.Background(), scrapertest.NewNopSettingsWithType(typ), cfg)
			require.NoError(t, err)
			err = c.Shutdown(context.Background())
			require.NoError(t, err)
		})
		t.Run(tt.name+"-lifecycle", func(t *testing.T) {
			firstRcvr, err := tt.createFn(context.Background(), scrapertest.NewNopSettingsWithType(typ), cfg)
			require.NoError(t, err)
			host := componenttest.NewNopHost()
			require.NoError(t, err)
			require.NoError(t, firstRcvr.Start(context.Background(), host))
			require.NoError(t, firstRcvr.Shutdown(context.Background()))
			secondRcvr, err := tt.createFn(context.Background(), scrapertest.NewNopSettingsWithType(typ), cfg)
			require.NoError(t, err)
			require.NoError(t, secondRcvr.Start(context.Background(), host))
			require.NoError(t, secondRcvr.Shutdown(context.Background()))
		})
	}
}
