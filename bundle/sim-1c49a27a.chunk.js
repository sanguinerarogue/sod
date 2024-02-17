import{S as e}from"./shaman_shields-f1b80a9f.chunk.js";import{dG as t,c as a,at as l,dp as s,dH as n,dA as i,r as o,bD as p,bF as c,ae as d,dx as r,af as S,E as I,a as u,P as h,b2 as m,d as f,q as v,x as k,S as P,aF as A,l as g,C as y,au as w}from"./detailed_results-34fd0d97.chunk.js";import{k as O,m as H,a as M,r as E,b as T,l as b,M as N,n as D,o as L,T as C,p as V,I as q}from"./preset_utils-828125b5.chunk.js";const B=O({fieldName:"syncType",label:"Sync/Stagger Setting",labelTooltip:"Choose your sync or stagger option Perfect\n\t\t<ul>\n\t\t\t<li><div>Auto: Will auto pick sync options based on your weapons attack speeds</div></li>\n\t\t\t<li><div>None: No Sync or Staggering, used for mismatched weapon speeds</div></li>\n\t\t\t<li><div>Perfect Sync: Makes your weapons always attack at the same time, for match weapon speeds</div></li>\n\t\t\t<li><div>Delayed Offhand: Adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window</div></li>\n\t\t</ul>",values:[{name:"Automatic",value:t.Auto},{name:"None",value:t.NoSync},{name:"Perfect Sync",value:t.SyncMainhandOffhandSwings},{name:"Delayed Offhand",value:t.DelayOffhandSwings}]}),R={items:[{id:211789},{id:209422},{id:14573},{id:213087,enchant:247},{id:211512,enchant:847,rune:408496},{id:209524,enchant:823},{id:211423,rune:408507},{id:16659},{id:10410,rune:425336},{id:211511},{id:211467},{id:13097},{id:211449},{id:4381},{id:209436,enchant:241},{id:1454,enchant:241},{id:209575}]},x={items:[{id:215166},{id:213344},{id:213304},{id:16336,enchant:849},{id:213314,enchant:866,rune:408496},{id:213317,enchant:856},{id:211423,enchant:856,rune:408507},{id:213325,rune:408498},{id:213333,rune:425336},{id:213339,enchant:849,rune:425874},{id:19512},{id:213284},{id:211449},{id:213348},{id:213409,enchant:7210},{id:213442,enchant:7210},{id:209575}]},W={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:5394,rank:1}}},doAtValue:{const:{val:"-6s"}}},{action:{castSpell:{spellId:{spellId:8160,rank:2}}},doAtValue:{const:{val:"-4.5s"}}},{action:{castSpell:{spellId:{spellId:6363,rank:2}}},doAtValue:{const:{val:"-3s"}}},{action:{castSpell:{spellId:{spellId:20572}}},doAtValue:{const:{val:"-1.5s"}}}],priorityList:[{action:{castSpell:{spellId:{spellId:20572}}}},{action:{condition:{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"75%"}}}},castSpell:{spellId:{spellId:425336}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:8052,rank:2}}}}},castSpell:{spellId:{spellId:8052,rank:2}}}},{action:{castSpell:{spellId:{spellId:408507}}}},{action:{castSpell:{spellId:{spellId:8056,rank:1}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:6363,rank:2}}}}},castSpell:{spellId:{spellId:6363,rank:2}}}}]},F={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:10495,rank:2}}},doAtValue:{const:{val:"-4.75s"}}},{action:{castSpell:{spellId:{spellId:8161,rank:3}}},doAtValue:{const:{val:"-3.25s"}}},{action:{castSpell:{spellId:{spellId:6365,rank:4}}},doAtValue:{const:{val:"-1.75s"}}},{action:{castSpell:{spellId:{spellId:20572}}},doAtValue:{const:{val:"-.25s"}}}],priorityList:[{action:{castSpell:{spellId:{spellId:437349}}}},{action:{autocastOtherCooldowns:{}}},{action:{condition:{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"65%"}}}},castSpell:{spellId:{spellId:425336}}}},{action:{condition:{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"55%"}}}},castSpell:{spellId:{itemId:6149}}}},{action:{castSpell:{spellId:{spellId:17364,rank:1}}}},{action:{castSpell:{spellId:{spellId:408507}}}},{action:{condition:{cmp:{op:"OpEq",lhs:{const:{val:"5"}},rhs:{auraNumStacks:{auraId:{spellId:408505}}}}},castSpell:{spellId:{spellId:408490}}}},{action:{condition:{cmp:{op:"OpEq",lhs:{const:{val:"5"}},rhs:{auraNumStacks:{auraId:{spellId:408505}}}}},castSpell:{spellId:{spellId:930,rank:2}}}},{action:{condition:{cmp:{op:"OpEq",lhs:{const:{val:"5"}},rhs:{auraNumStacks:{auraId:{spellId:408505}}}}},castSpell:{spellId:{spellId:10391,rank:7}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:10447,rank:4}}}}},castSpell:{spellId:{spellId:10447,rank:4}}}},{action:{castSpell:{spellId:{spellId:10412,rank:5}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:6365,rank:4}}}}},castSpell:{spellId:{spellId:6365,rank:4}}}}]},U={type:"TypeAPL",prepullActions:[{action:{castSpell:{spellId:{spellId:10495,rank:2}}},doAtValue:{const:{val:"-4.75s"}}},{action:{castSpell:{spellId:{spellId:8161,rank:3}}},doAtValue:{const:{val:"-3.25s"}}},{action:{castSpell:{spellId:{spellId:6365,rank:4}}},doAtValue:{const:{val:"-1.75s"}}},{action:{castSpell:{spellId:{spellId:20572}}},doAtValue:{const:{val:"-.25s"}}}],priorityList:[{action:{castSpell:{spellId:{spellId:437362}}}},{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:437362}}}}},autocastOtherCooldowns:{}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"65%"}}}},{spellIsReady:{spellId:{spellId:425336}}}]}},sequence:{name:"mana",actions:[{itemSwap:{swapSet:"Swap1"}},{castSpell:{spellId:{spellId:425336}}},{itemSwap:{swapSet:"Main"}},{resetSequence:{sequenceName:"mana"}}]}}},{action:{condition:{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"55%"}}}},castSpell:{spellId:{itemId:6149}}}},{action:{castSpell:{spellId:{spellId:17364,rank:1}}}},{action:{castSpell:{spellId:{spellId:408507}}}},{action:{condition:{cmp:{op:"OpEq",lhs:{const:{val:"5"}},rhs:{auraNumStacks:{auraId:{spellId:408505}}}}},castSpell:{spellId:{spellId:408490}}}},{action:{condition:{cmp:{op:"OpEq",lhs:{const:{val:"5"}},rhs:{auraNumStacks:{auraId:{spellId:408505}}}}},castSpell:{spellId:{spellId:930,rank:2}}}},{action:{condition:{cmp:{op:"OpEq",lhs:{const:{val:"5"}},rhs:{auraNumStacks:{auraId:{spellId:408505}}}}},castSpell:{spellId:{spellId:10391,rank:7}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:10447,rank:4}}}}},castSpell:{spellId:{spellId:10447,rank:4}}}},{action:{castSpell:{spellId:{spellId:10412,rank:5}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:6365,rank:4}}}}},castSpell:{spellId:{spellId:6365,rank:4}}}}]},j=H("Blank",{items:[]}),_=H("Phase 1",R),G=H("Phase 2",x),z={[a.Phase1]:[_],[a.Phase2]:[G]},J=z[l][0],K=M("Phase 1",W),Q=M("Phase 2",F),X=M("Phase 2 (Swap)",U),Y={[a.Phase1]:[K],[a.Phase2]:[Q,X]},Z={25:Y[a.Phase1][0],40:Y[a.Phase2][0]},$={name:"Phase 1",data:s.create({talentsString:"-5005202101"})},ee={name:"Phase 2",data:s.create({talentsString:"-5005202105023051"})},te={[a.Phase1]:[$],[a.Phase2]:[ee]},ae=te[l][0],le=n.create({shield:i.LightningShield,syncType:t.Auto}),se=o.create({mainHandImbue:p.WindfuryWeapon,offHandImbue:p.WindfuryWeapon,enchantedSigil:c.InnovationSigil}),ne=d.create({arcaneBrilliance:!0,aspectOfTheLion:!0,battleShout:r.TristateEffectImproved,divineSpirit:!0,giftOfTheWild:r.TristateEffectImproved,leaderOfThePack:!0,manaSpringTotem:r.TristateEffectImproved,moonkinAura:!0,strengthOfEarthTotem:r.TristateEffectImproved,trueshotAura:!0}),ie=S.create({curseOfElements:!0,curseOfRecklessness:!0,dreamstate:!0,faerieFire:!0,improvedScorch:!0,sunderArmor:!0}),oe={profession1:I.Enchanting,profession2:I.Leatherworking},pe=E(P.SpecEnhancementShaman,{cssClass:"enhancement-shaman-sim-ui",cssScheme:"shaman",knownIssues:[],epStats:[u.StatIntellect,u.StatAgility,u.StatStrength,u.StatAttackPower,u.StatMeleeHit,u.StatMeleeCrit,u.StatMeleeHaste,u.StatArmorPenetration,u.StatExpertise,u.StatSpellPower,u.StatFirePower,u.StatNaturePower,u.StatSpellCrit,u.StatSpellHit,u.StatSpellHaste],epPseudoStats:[h.PseudoStatMainHandDps,h.PseudoStatOffHandDps],epReferenceStat:u.StatAttackPower,displayStats:[u.StatHealth,u.StatStamina,u.StatStrength,u.StatAgility,u.StatIntellect,u.StatAttackPower,u.StatMeleeHit,u.StatMeleeCrit,u.StatMeleeHaste,u.StatExpertise,u.StatArmorPenetration,u.StatSpellPower,u.StatSpellHit,u.StatSpellCrit,u.StatSpellHaste],defaults:{race:m.RaceTroll,gear:J.gear,epWeights:T.fromMap({[u.StatIntellect]:.02,[u.StatAgility]:1.12,[u.StatStrength]:2.29,[u.StatSpellPower]:1.15,[u.StatFirePower]:.63,[u.StatNaturePower]:.48,[u.StatSpellHit]:.03,[u.StatSpellCrit]:1.94,[u.StatSpellHaste]:2.97,[u.StatAttackPower]:1,[u.StatMeleeHit]:9.62,[u.StatMeleeCrit]:14.8,[u.StatMeleeHaste]:11.84,[u.StatArmorPenetration]:.35,[u.StatExpertise]:1.92},{[h.PseudoStatMainHandDps]:8.15,[h.PseudoStatOffHandDps]:5.81}),consumes:se,talents:ae.data,specOptions:le,other:oe,raidBuffs:ne,partyBuffs:f.create({}),individualBuffs:v.create({}),debuffs:ie},playerIconInputs:[e()],includeBuffDebuffInputs:[b,N,D],excludeBuffDebuffInputs:[L],otherInputs:{inputs:[B,C,V]},itemSwapConfig:{itemSlots:[k.ItemSlotMainHand,k.ItemSlotOffHand],note:"Swap items are given the highest available rank of Rockbiter Weapon"},customSections:[],encounterPicker:{showExecuteProportion:!1},presets:{talents:[...te[a.Phase1],...te[l]],rotations:[...Y[a.Phase1],...Y[l]],gear:[j,...z[a.Phase1],...z[l]]},autoRotation:e=>Z[e.getLevel()].rotation.rotation,raidSimPresets:[{spec:P.SpecBalanceDruid,tooltip:A[P.SpecBalanceDruid],defaultName:"Balance",iconUrl:g(y.ClassDruid,0),talents:ae.data,specOptions:le,consumes:se,otherDefaults:oe,defaultFactionRaces:{[w.Unknown]:m.RaceUnknown,[w.Alliance]:m.RaceNightElf,[w.Horde]:m.RaceTauren},defaultGear:{[w.Unknown]:{},[w.Alliance]:{1:z[a.Phase1][0].gear},[w.Horde]:{1:z[a.Phase1][0].gear}}}]});class ce extends q{constructor(e,t){super(e,t,pe)}}export{ce as E};
//# sourceMappingURL=sim-1c49a27a.chunk.js.map
