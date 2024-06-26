package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) registerShadowflameSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneBootsShadowflame) {
		return
	}

	hasInvocationRune := warlock.HasRune(proto.WarlockRune_RuneBeltInvocation)
	hasPandemicRune := warlock.HasRune(proto.WarlockRune_RuneHelmPandemic)
	hasUnstableAffliction := warlock.HasRune(proto.WarlockRune_RuneBracerUnstableAffliction)

	// https://www.wowhead.com/classic/spell=18275/shadow-mastery
	// The Shadowflame periodic damage is NOT buffed by Shadow Mastery as it's not affected by haunt
	baseSpellCoeff := 0.20
	dotSpellCoeff := 0.107
	baseDamage := warlock.baseRuneAbilityDamage() * 2.26
	dotDamage := warlock.baseRuneAbilityDamage() * 0.61

	warlock.ShadowflameDot = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 426325},
		SpellSchool: core.SpellSchoolFire | core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       WarlockFlagAffliction | WarlockFlagDestruction,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Shadowflame" + warlock.Label,
			},

			NumberOfTicks:    4,
			TickLength:       time.Second * 2,
			BonusCoefficient: dotSpellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, dotDamage, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				var result *core.SpellResult
				if hasPandemicRune {
					// We add the crit damage bonus and remove it after the call to not affect the initial damage portion of the spell
					dot.Spell.CritDamageBonus += 1
					result = dot.CalcSnapshotDamage(sim, target, dot.OutcomeTickSnapshotCrit)
					dot.Spell.CritDamageBonus -= 1
				} else {
					result = dot.CalcSnapshotDamage(sim, target, dot.OutcomeTick)
				}
				dot.Spell.DealPeriodicDamage(sim, result)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			spell.Dot(target).Apply(sim)
		},
	})

	warlock.Shadowflame = warlock.RegisterSpell(core.SpellConfig{
		SpellCode:   SpellCode_WarlockShadowflame,
		ActionID:    core.ActionID{SpellID: 426320},
		SpellSchool: core.SpellSchoolFire | core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | core.SpellFlagResetAttackSwing | WarlockFlagAffliction | WarlockFlagDestruction,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.27,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 15 * time.Second,
			},
		},

		DamageMultiplier: 1 + warlock.improvedImmolateBonus(),
		ThreatMultiplier: 1,
		BonusCoefficient: baseSpellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				// UA, Immo, Shadowflame exclusivity
				if hasUnstableAffliction && warlock.UnstableAffliction.Dot(target).IsActive() {
					warlock.UnstableAffliction.Dot(target).Deactivate(sim)
				}
				immoDot := warlock.getActiveImmolateSpell(target)
				if immoDot != nil {
					immoDot.Dot(target).Deactivate(sim)
				}

				if hasInvocationRune && spell.Dot(target).IsActive() {
					warlock.InvocationRefresh(sim, spell.Dot(target))
				}

				warlock.ShadowflameDot.Cast(sim, target)
			}
		},
	})
}
