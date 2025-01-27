package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var ItemSetEclipseOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Eclipse of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Balance2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Balance4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Balance6PBonus()
		},
	},
})

// Increases the damage done and damage radius of Starfall's stars and Hurricane by 25%.
func (druid *Druid) applyT2Balance2PBonus() {
	label := "S03 - Item - T2 - Druid - Balance 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
	}).AttachSpellMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_DruidHurricane | ClassSpellMask_DruidStarfall | ClassSpellMask_DruidStarfallTick | ClassSpellMask_DruidStarfallSplash,
		IntValue:  25,
	}))
}

// Your Wrath casts have a 10% chance to summon a stand of 3 Treants to attack your target for until cancelled.
func (druid *Druid) applyT2Balance4PBonus() {
	label := "S03 - Item - T2 - Druid - Balance 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidWrath|ClassSpellMask_DruidStarsurge) && !druid.t26pcTreants.IsActive() && sim.Proc(0.10, "Summon Treants") {
				druid.t26pcTreants.EnableWithTimeout(sim, druid.t26pcTreants, time.Second*15)
			}
		},
	}))
}

// Your Wrath critical strikes have a 50% chance to make your next Starfire deal 10% increased damage, stacking up to 3 times.
func (druid *Druid) applyT2Balance6PBonus() {
	label := "S03 - Item - T2 - Druid - Balance 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	damageMod := druid.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_DruidStarfire,
	})

	buffAura := druid.RegisterAura(core.Aura{
		ActionID:  core.ActionID{SpellID: 467088},
		Label:     "Astral Power",
		Duration:  time.Second * 15,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			damageMod.UpdateIntValue(int64(10 * newStacks))
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell.Matches(ClassSpellMask_DruidStarfire) {
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			damageMod.Deactivate()
		},
	})

	core.MakePermanent(druid.RegisterAura(core.Aura{
		Label: label,
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.Matches(ClassSpellMask_DruidWrath|ClassSpellMask_DruidStarsurge) && result.DidCrit() && sim.Proc(0.50, "Astral Power") {
				buffAura.Activate(sim)
				buffAura.AddStack(sim)
			}
		},
	}))
}

var ItemSetCunningOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Cunning of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Feral2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Feral4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Feral6PBonus()
		},
	},
})

// Increases the duration of Rake by 6 sec and its periodic damage by 50%.
func (druid *Druid) applyT2Feral2PBonus() {
	label := "S03 - Item - T2 - Druid - Feral 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 467207},
		Label:    "S03 - Item - T2- Druid - Feral 2P Bonus",
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, dot := range druid.Rake.Dots() {
				if dot == nil {
					continue
				}

				dot.NumberOfTicks += int32(6 / dot.TickLength.Seconds())
				dot.RecomputeAuraDuration()
				oldOnSnapshot := dot.OnSnapshot
				dot.OnSnapshot = func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					oldOnSnapshot(sim, target, dot, isRollover)
					dot.SnapshotAttackerMultiplier *= 1.50
				}
			}
		},
	})
}

// Your critical strike chance is increased by 15% while Tiger's Fury is active.
func (druid *Druid) applyT2Feral4PBonus() {
	label := "S03 - Item - T2 - Druid - Feral 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			oldOnGain := druid.TigersFuryAura.OnGain
			druid.TigersFuryAura.OnGain = func(aura *core.Aura, sim *core.Simulation) {
				oldOnGain(aura, sim)
				druid.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: 15 * core.CritRatingPerCritChance})
			}
			oldOnExpire := druid.TigersFuryAura.OnExpire
			druid.TigersFuryAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
				oldOnExpire(aura, sim)
				druid.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: -15 * core.CritRatingPerCritChance})
			}
		},
	})
}

// Your Shred and Mangle(Cat) abilities deal 10% increased damage per your Bleed effect on the target, up to a maximum of 20% increase.
func (druid *Druid) applyT2Feral6PBonus() {
	damageMod := druid.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Flat,
		ClassMask: ClassSpellMask_DruidShred | ClassSpellMask_DruidMangleCat | ClassSpellMask_DruidFerociousBite,
	})

	core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
		Name:           "S03 - Item - T2 - Druid - Feral 6P Bonus",
		Callback:       core.CallbackOnApplyEffects,
		ProcChance:     1,
		ClassSpellMask: ClassSpellMask_DruidShred | ClassSpellMask_DruidMangleCat | ClassSpellMask_DruidFerociousBite,
		Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			damageMod.Activate()
			damageMod.UpdateIntValue(int64(10 * druid.BleedsActive))
		},
	})
}

var ItemSetFuryOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Fury of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Guardian2PBonus()
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Guardian4PBonus()
		},
		6: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyT2Guardian6PBonus()
		},
	},
})

// Swipe(Bear) also causes your Maul to hit 1 additional target for the next 6 sec.
func (druid *Druid) applyT2Guardian2PBonus() {
	label := "S03 - Item - T2 - Druid - Guardian 2P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO
		},
	})
}

// Your Mangle(Bear), Swipe(Bear), Maul, and Lacerate abilities gain 5% increased critical strike chance against targets afflicted by your Lacerate.
func (druid *Druid) applyT2Guardian4PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}

	label := "S03 - Item - T2 - Druid - Guardian 4P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO
		},
	})
}

// Your Swipe now spreads your Lacerate from your primary target to other targets it strikes.
func (druid *Druid) applyT2Guardian6PBonus() {
	if !druid.HasRune(proto.DruidRune_RuneLegsLacerate) {
		return
	}

	label := "S03 - Item - T2 - Druid - Guardian 6P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			// TODO
		},
	})
}

var ItemSetBountyOfStormrage = core.NewItemSet(core.ItemSet{
	Name: "Bounty of Stormrage",
	Bonuses: map[int32]core.ApplyEffect{
		// Your healing spell critical strikes trigger the Dreamstate effect, granting you 50% of your mana regeneration while casting for 8 sec.
		2: func(agent core.Agent) {
		},
		// Your non-periodic spell critical strikes reduce the casting time of your next Healing Touch, Regrowth, or Nourish spell by 0.5 sec.
		4: func(agent core.Agent) {
		},
		// Increases healing from Wild Growth by 10%. In addition, Wild Growth can now be used in Moonkin Form, and its healing is increased by an additional 50% in that form.
		6: func(agent core.Agent) {
		},
	},
})

var ItemSetHaruspexsGarb = core.NewItemSet(core.ItemSet{
	Name: "Haruspex's Garb",
	Bonuses: map[int32]core.ApplyEffect{
		// Increases damage and healing done by magical spells and effects by up to 12.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
		3: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyZGBalance3PBonus()
		},
		5: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			druid.applyZGBalance5PBonus()
		},
	},
})

// Reduces the cast time and global cooldown of Starfire by 0.5 sec.
func (druid *Druid) applyZGBalance3PBonus() {
	label := "S03 - Item - ZG - Druid - Balance 3P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range druid.Starfire {
				if spell == nil {
					continue
				}

				spell.DefaultCast.CastTime -= time.Millisecond * 500
				spell.DefaultCast.GCD -= time.Millisecond * 500
			}
		},
	})
}

// Increases the critical strike chance of Wrath by 10%.
func (druid *Druid) applyZGBalance5PBonus() {
	label := "S03 - Item - ZG - Druid - Balance 5P Bonus"
	if druid.HasAura(label) {
		return
	}

	druid.RegisterAura(core.Aura{
		Label: label,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range druid.Wrath {
				if spell == nil {
					continue
				}

				spell.BonusCritRating += 10 * core.SpellCritRatingPerCritChance
			}
		},
	})
}
