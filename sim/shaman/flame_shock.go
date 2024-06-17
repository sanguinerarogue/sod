package shaman

import (
	"fmt"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const FlameShockRanks = 6

var FlameShockSpellId = [FlameShockRanks + 1]int32{0, 8050, 8052, 8053, 10447, 10448, 29228}
var FlameShockBaseDamage = [FlameShockRanks + 1]float64{0, 25, 51, 95, 164, 245, 292}
var FlameShockBaseDotDamage = [FlameShockRanks + 1]float64{0, 28, 48, 96, 168, 256, 320}
var FlameShockBaseSpellCoef = [FlameShockRanks + 1]float64{0, .134, .198, .214, .214, .214, .214}
var FlameShockDotSpellCoef = [FlameShockRanks + 1]float64{0, .063, .093, .1, .1, .1, .1}
var FlameShockManaCost = [FlameShockRanks + 1]float64{0, 55, 95, 160, 250, 345, 410}
var FlameShockLevel = [FlameShockRanks + 1]int{0, 10, 18, 28, 40, 52, 60}

func (shaman *Shaman) registerFlameShockSpell(shockTimer *core.Timer) {
	shaman.FlameShock = make([]*core.Spell, FlameShockRanks+1)
	shaman.FlameShockDots = make([]*core.Spell, FlameShockRanks+1)

	for rank := 1; rank <= FlameShockRanks; rank++ {
		if FlameShockLevel[rank] <= int(shaman.Level) {
			shaman.FlameShock[rank] = shaman.RegisterSpell(shaman.newFlameShockSpell(rank, shockTimer))
		}
	}
}

func (shaman *Shaman) newFlameShockSpell(rank int, shockTimer *core.Timer) core.SpellConfig {
	spellId := FlameShockSpellId[rank]
	baseDamage := FlameShockBaseDamage[rank]
	baseSpellCoeff := FlameShockBaseSpellCoef[rank]
	manaCost := FlameShockManaCost[rank]
	level := FlameShockLevel[rank]

	hasBurnRune := shaman.HasRune(proto.ShamanRune_RuneHelmBurn)
	hasPowerSurgeRune := shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge)

	spell := shaman.newShockSpellConfig(
		core.ActionID{SpellID: spellId},
		core.SpellSchoolFire,
		manaCost,
		shockTimer,
	)

	spell.SpellCode = SpellCode_ShamanFlameShockDirect
	spell.RequiredLevel = level
	spell.Rank = rank

	spell.Cast.IgnoreHaste = true

	spell.DamageMultiplier += shaman.burnFlameShockDamageMultiplier()
	spell.BonusCoefficient = baseSpellCoeff

	shaman.FlameShockDots[rank] = shaman.RegisterSpell(shaman.newFlameShockDotSpellConfig(rank))

	results := make([]*core.SpellResult, min(core.TernaryInt32(hasBurnRune, BurnFlameShockTargetCount, 1), shaman.Env.GetNumTargets()))
	spell.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		for idx := range results {
			results[idx] = spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			target = sim.Environment.NextTargetUnit(target)
		}

		for _, result := range results {
			spell.DealDamage(sim, result)
			if result.Landed() {
				shaman.FlameShockDots[rank].Cast(sim, target)

				if shaman.HasRune(proto.ShamanRune_RuneLegsAncestralGuidance) {
					shaman.lastFlameShockTarget = target
				}

				if hasPowerSurgeRune && sim.Proc(ShamanPowerSurgeProcChance, "Power Surge Proc") {
					shaman.PowerSurgeAura.Activate(sim)
				}
			}
		}
	}

	return spell
}

func (shaman *Shaman) newFlameShockDotSpellConfig(rank int) core.SpellConfig {
	hasBurnRune := shaman.HasRune(proto.ShamanRune_RuneHelmBurn)
	hasLavaBurstRune := shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst)
	hasPowerSurgeRune := shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge)
	hasMoltenBlastRune := shaman.HasRune(proto.ShamanRune_RuneHandsMoltenBlast)
	hasOverloadRune := shaman.HasRune(proto.ShamanRune_RuneChestOverload)
	hasStormEarthAndFireRune := shaman.HasRune(proto.ShamanRune_RuneCloakStormEarthAndFire)

	numTicks := 4 + core.TernaryInt32(hasBurnRune, BurnFlameShockBonusTicks, 0)
	tickDuration := time.Second * 3

	spellId := FlameShockSpellId[rank]
	baseDotDamage := FlameShockBaseDotDamage[rank] / float64(numTicks)
	dotSpellCoeff := FlameShockDotSpellCoef[rank]
	level := FlameShockLevel[rank]

	return core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId}.WithTag(1),
		SpellSchool: core.SpellSchoolFire,
		SpellCode:   SpellCode_ShamanFlameShockDot,
		ProcMask:    core.ProcMaskSpellDamage,
		DefenseType: core.DefenseTypeMagic,
		Flags:       core.SpellFlagPureDot,

		Rank:          rank,
		RequiredLevel: level,

		DamageMultiplier: 1 + shaman.burnFlameShockDamageMultiplier() + core.TernaryFloat64(hasStormEarthAndFireRune, .6, 0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Flame Shock (Rank %d)", rank),
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if hasLavaBurstRune {
						shaman.LavaBurst.BonusCritRating += 100 * core.CritRatingPerCritChance
						if hasOverloadRune {
							shaman.LavaBurstOverload.BonusCritRating += 100 * core.CritRatingPerCritChance
						}
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if hasLavaBurstRune {
						shaman.LavaBurst.BonusCritRating -= 100 * core.CritRatingPerCritChance
						if hasOverloadRune {
							shaman.LavaBurstOverload.BonusCritRating -= 100 * core.CritRatingPerCritChance
						}
					}
				},
			},

			NumberOfTicks:    int32(numTicks),
			TickLength:       tickDuration,
			BonusCoefficient: dotSpellCoeff,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, baseDotDamage, isRollover)
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)

				if hasMoltenBlastRune && sim.Proc(ShamanMoltenBlastResetChance, "Molten Blast Reset") {
					shaman.MoltenBlast.CD.Reset()
				}

				if hasPowerSurgeRune && sim.Proc(ShamanPowerSurgeProcChance, "Power Surge Proc") {
					shaman.PowerSurgeAura.Activate(sim)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},

		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.Spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDamage := baseDotDamage / float64(numTicks)
				return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}
}
