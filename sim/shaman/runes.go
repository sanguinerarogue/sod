package shaman

import (
	"slices"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (shaman *Shaman) ApplyRunes() {
	// Chest
	shaman.applyDualWieldSpec()
	// shaman.applyHealingRain()
	// shaman.applyOverload()
	shaman.applyShieldMastery()
	shaman.applyTwoHandedMastery()

	// Hands
	shaman.applyLavaBurst()
	shaman.applyLavaLash()
	shaman.applyMoltenBlast()

	// Waist
	// shaman.applyFireNova()
	shaman.applyMaelstromWeapon()
	shaman.applyPowerSurge()

	// Legs
	shaman.applyAncestralGuidance()
	// shaman.applyEarthShield()
	shaman.applyShamanisticRage()
	shaman.applyWayOfEarth()

	// Feet
	shaman.applyAncestralAwakening()
	// shaman.applyDecoyTotem()
	// shaman.applySpiritOfTheAlpha()
}

func (shaman *Shaman) applyDualWieldSpec() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestDualWieldSpec) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Dual Wield Specialization",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestDualWieldSpec)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyShieldMastery() {
	if !shaman.HasRune(proto.ShamanRune_RuneChestShieldMastery) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Shield Mastery",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneChestShieldMastery)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyTwoHandedMastery() {
	if !shaman.HasRune(proto.ShamanRune_RuneTwoHandedMastery) {
		return
	}

	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery Proc",
		ActionID: core.ActionID{SpellID: 436365},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyMeleeSpeed(sim, 1.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.MultiplyAttackSpeed(sim, 1/1.3)
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Two-Handed Mastery",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneTwoHandedMastery)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMeleeMH) {
				return
			}

			if shaman.MainHand().HandType == proto.HandType_HandTypeTwoHand {
				procAura.Activate(sim)
			} else {
				procAura.Deactivate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyMaelstromWeapon() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistMaelstromWeapon) {
		return
	}

	shaman.MaelstromWeaponAura = shaman.RegisterAura(core.Aura{
		Label:    "Maelstrom Weapon",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneWaistMaelstromWeapon)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

const ShamanPowerSurgeProcChance = .05

func (shaman *Shaman) applyPowerSurge() {
	if !shaman.HasRune(proto.ShamanRune_RuneWaistPowerSurge) {
		return
	}

	var affectedSpells []*core.Spell

	procAura := shaman.RegisterAura(core.Aura{
		Label:    "Power Surge Proc",
		ActionID: core.ActionID{SpellID: 440285},
		Duration: time.Second * 10,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			affectedSpells = core.FilterSlice(
				core.Flatten([][]*core.Spell{
					shaman.ChainLightning,
					{shaman.ChainHeal},
					{shaman.LavaBurst},
				}), func(spell *core.Spell) bool { return spell != nil },
			)
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.LavaBurst.CD.Reset()
			for _, spell := range affectedSpells {
				spell.CD.Reset()
			}
			core.Each(affectedSpells, func(spell *core.Spell) { spell.CastTimeMultiplier -= 1 })
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) { spell.CastTimeMultiplier += 1 })
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.SpellCode != int(SpellCode_ShamanLightningBolt) && spell.SpellCode != int(SpellCode_ShamanChainLightning) && spell != shaman.LavaBurst {
				return
			}
			aura.Deactivate(sim)
		},
	})

	shaman.RegisterAura(core.Aura{
		Label:    "Power Surge",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneWaistPowerSurge)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && slices.Contains(shaman.FlameShock, spell) && sim.RandomFloat("Power Surge Proc") < ShamanPowerSurgeProcChance {
				procAura.Activate(sim)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() && slices.Contains(shaman.FlameShock, spell) && sim.RandomFloat("Power Surge Proc") < ShamanPowerSurgeProcChance {
				procAura.Activate(sim)
			}
		},
	})
}

func (shaman *Shaman) applyWayOfEarth() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsWayOfEarth) {
		return
	}

	shaman.RegisterAura(core.Aura{
		Label:    "Way of Earth",
		ActionID: core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsWayOfEarth)},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
	})
}

func (shaman *Shaman) applyShamanisticRage() {
	if !shaman.HasRune(proto.ShamanRune_RuneLegsShamanisticRage) {
		return
	}

	actionID := core.ActionID{SpellID: int32(proto.ShamanRune_RuneLegsShamanisticRage)}
	manaMetrics := shaman.NewManaMetrics(actionID)
	srAura := shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			srAura.Activate(sim)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				NumTicks: 15,
				Period:   time.Second * 1,
				OnAction: func(sim *core.Simulation) {
					mana := max(
						shaman.GetStat(stats.AttackPower)*0.15,
						shaman.GetStat(stats.SpellPower)*0.10,
						shaman.GetStat(stats.Healing)*0.06,
					)
					shaman.AddMana(sim, mana, manaMetrics)
				},
			})
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.2
		},
	})
}
