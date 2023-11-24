package warrior

import (
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (warrior *Warrior) registerBerserkerRageSpell() {
	actionID := core.ActionID{SpellID: 18499}

	warrior.BerserkerRage = warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 30,
			},
		},
	})
}

func (warrior *Warrior) ShouldBerserkerRage(sim *core.Simulation) bool {
	return warrior.Talents.ImprovedBerserkerRage > 0 && warrior.CurrentRage() < 80 && warrior.BerserkerRage.IsReady(sim)
}