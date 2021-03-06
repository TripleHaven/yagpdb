package featureflags

import (
	"sync"

	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/eventsystem"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/pubsub"
)

var logger = common.GetPluginLogger(&Plugin{})

// Plugin represents the mqueue plugin
type Plugin struct {
	stopBGWorker chan *sync.WaitGroup
}

// PluginInfo implements common.Plugin
func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "featureflags",
		SysName:  "featureflags",
		Category: common.PluginCategoryCore,
	}
}

// RegisterPlugin registers the mqueue plugin into the plugin system and also initializes it
func RegisterPlugin() {
	p := &Plugin{
		stopBGWorker: make(chan *sync.WaitGroup),
	}

	common.RegisterPlugin(p)

	pubsub.AddHandler("feature_flags_updated", handleInvalidateCacheFor, nil)
}

// Invalidate the cache when the rules have changed
func handleInvalidateCacheFor(event *pubsub.Event) {
	cacheID := (event.TargetGuildInt >> 22) % int64(len(caches))
	caches[cacheID].invalidateGuild(event.TargetGuildInt)
}

var _ bot.BotInitHandler = (*Plugin)(nil)

// BotInit implements bot.BotInitHandler
func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLastLegacy(p, func(evt *eventsystem.EventData) {
		cacheID := (evt.GuildDelete().ID >> 22) % int64(len(caches))
		caches[cacheID].invalidateGuild(evt.GuildDelete().ID)
	}, eventsystem.EventGuildDelete)
}
