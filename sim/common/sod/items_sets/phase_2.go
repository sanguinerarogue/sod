package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

var ItemSetHyperconductiveMendersMeditation = core.NewItemSet(core.ItemSet{
	Name: "Hyperconductive Mender's Meditation",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Spirit, 14)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MP5, 7)
		},
	},
})

var ItemSetHyperconductiveWizardsAttire = core.NewItemSet(core.ItemSet{
	Name: "Hyperconductive Wizard's Attire",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
			c.AddStat(stats.BonusArmor, 100)
		},
		3: func(agent core.Agent) {
			character := agent.GetCharacter()

			procAura := character.NewTemporaryStatsAura("Energized Hyperconductor Proc", core.ActionID{SpellID: 435978}, stats.Stats{stats.SpellPower: 40}, time.Second*10)

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "Energized Hyperconductor",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage,
				ProcChance: 0.10,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					procAura.Activate(sim)
				},
			})
		},
	},
})

var ItemSetIrradiatedGarments = core.NewItemSet(core.ItemSet{
	Name: "Irradiated Garments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1)
			c.AddStat(stats.SpellCrit, 1)
			c.AddStat(stats.Stamina, -5)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 11)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetInsulatedLeather = core.NewItemSet(core.ItemSet{
	Name: "Insulated Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1)
			c.AddStat(stats.SpellCrit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.PseudoStats.DaggersSkill += 3
			c.AddStat(stats.FeralAttackPower, 20)
		},
	},
})

var ItemSetInsulatedSorcerorLeather = core.NewItemSet(core.ItemSet{
	Name: "Insulated Sorceror's Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 16)
		},
		3: func(agent core.Agent) {},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetElectromanticDevastator = core.NewItemSet(core.ItemSet{
	Name: "Electromantic Devastator's Mail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 24)
			c.AddStat(stats.RangedAttackPower, 24)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			if !c.HasManaBar() {
				return
			}
			metrics := c.NewManaMetrics(core.ActionID{SpellID: 435982})
			proc := c.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 435981},
				SpellSchool: core.SpellSchoolHoly,
				ApplyEffects: func(sim *core.Simulation, u *core.Unit, spell *core.Spell) {
					c.AddMana(sim, 100, metrics)
				},
			})
			procChance := 0.05
			c.RegisterAura(core.Aura{
				Label:    "Electromantic Devastator's Mail 3pc",
				ActionID: core.ActionID{SpellID: 435982},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				// Modeled after WotLK JoW https://github.com/wowsims/wotlk/blob/master/sim/core/debuffs.go#L202
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {

					if spell.ProcMask == core.ProcMaskEmpty {
						return // Can't use Matches for ProcMaskEmpty
					}

					if spell.ProcMask.Matches(core.ProcMaskProc|core.ProcMaskSpellDamageProc) && !spell.Flags.Matches(core.SpellFlagNotAProc) {
						return // Phantom spells don't proc
					}

					if spell.ProcMask.Matches(core.ProcMaskWhiteHit | core.ProcMaskRanged) { // Ranged/melee can proc on miss
						if sim.RandomFloat("Electromantic Devastator's Mail 3pc") > procChance {
							return
						}
					} else { // Spell Casting only procs on hits
						if !result.Landed() {
							return
						}
						if sim.RandomFloat("Electromantic Devastator's Mail 3pc") > procChance {
							return
						}
					}
					proc.Cast(sim, result.Target)
				},
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetHazardSuit = core.NewItemSet(core.ItemSet{
	Name: "H.A.Z.A.R.D. Suit",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Defense, 7)
			c.AddStat(stats.AttackPower, 16)
			c.AddStat(stats.RangedAttackPower, 16)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
	},
})
