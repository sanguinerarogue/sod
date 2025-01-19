package core

import (
	"fmt"
	"slices"

	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

type OnItemSwap func(*Simulation, proto.ItemSlot)

type ItemSwap struct {
	character       *Character
	onSwapCallbacks [NumItemSlots][]OnItemSwap

	isFuryWarrior        bool
	mhCritMultiplier     float64
	ohCritMultiplier     float64
	rangedCritMultiplier float64

	// Which slots to actually swap.
	slots []proto.ItemSlot

	// Holds the original equip
	originalEquip Equipment
	// Holds the items that are selected for swapping
	swapEquip Equipment
	// Holds items that are currently not equipped
	unEquippedItems Equipment
	swapSet         proto.APLActionItemSwap_SwapSet
	equipmentStats  ItemSwapStats

	initialized bool
}

type ItemSwapStats struct {
	allSlots    stats.Stats
	weaponSlots stats.Stats
}

func (character *Character) enableItemSwap(itemSwap *proto.ItemSwap) {
	var swapItems Equipment
	hasItemSwap := make(map[proto.ItemSlot]bool)

	for idx, itemSpec := range itemSwap.Items {
		itemSlot := proto.ItemSlot(idx)
		if !slices.Contains(AllWeaponSlots(), itemSlot) {
			panic(fmt.Sprintf("Slot %d is not supported. Currently only Mainhand, Offhand and Ranged are supported.", itemSlot))
		}
		hasItemSwap[itemSlot] = itemSpec != nil && itemSpec.Id != 0
		swapItems[itemSlot] = toItem(itemSpec)
	}

	has2HSwap := swapItems[proto.ItemSlot_ItemSlotMainHand].HandType == proto.HandType_HandTypeTwoHand
	hasMhEquipped := character.HasMHWeapon()
	hasOhEquipped := character.HasOHWeapon()

	// Handle MH and OH together, because present MH + empty OH --> swap MH and unequip OH
	if hasItemSwap[proto.ItemSlot_ItemSlotOffHand] && hasMhEquipped {
		hasItemSwap[proto.ItemSlot_ItemSlotMainHand] = true
	}

	if has2HSwap && hasOhEquipped {
		hasItemSwap[proto.ItemSlot_ItemSlotOffHand] = true
	}

	slots := SetToSortedSlice(hasItemSwap)

	if len(slots) == 0 {
		return
	}

	var prepullBonusStats stats.Stats
	if itemSwap.PrepullBonusStats != nil {
		prepullBonusStats = stats.FromFloatArray(itemSwap.PrepullBonusStats.Stats)
	}

	equipmentStats := calcItemSwapStatsOffset(character.Equipment, swapItems, prepullBonusStats, slots)

	character.ItemSwap = ItemSwap{
		slots:           slots,
		originalEquip:   character.Equipment,
		swapEquip:       swapItems,
		unEquippedItems: swapItems,
		equipmentStats:  equipmentStats,
		swapSet:         proto.APLActionItemSwap_Main,
		initialized:     false,
	}
}

func (swap *ItemSwap) initialize(character *Character) {
	swap.character = character
}

func (character *Character) RegisterItemSwapCallback(slots []proto.ItemSlot, callback OnItemSwap) {
	if character == nil || !character.ItemSwap.IsEnabled() || len(slots) == 0 {
		return
	}

	if (character.Env != nil) && character.Env.IsFinalized() {
		panic("Tried to add a new item swap callback in a finalized environment!")
	}

	for _, slot := range slots {
		character.ItemSwap.onSwapCallbacks[slot] = append(character.ItemSwap.onSwapCallbacks[slot], callback)
	}
}

// Helper for handling Effects that use PPMManager to toggle the aura on/off
func (swap *ItemSwap) RegisterOnSwapItemForEffectWithPPMManager(effectID int32, ppm float64, ppmm *PPMManager, aura *Aura) {
	slots := swap.EligibleSlotsForEffect(effectID)
	character := swap.character
	character.RegisterItemSwapCallback(slots, func(sim *Simulation, _ proto.ItemSlot) {
		procMask := character.GetProcMaskForEnchant(effectID)
		*ppmm = character.AutoAttacks.NewPPMManager(ppm, procMask)

		if ppmm.Chance(procMask) == 0 {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})

}

// Helper for handling Effects that use the effectID to toggle the aura on and off
func (swap *ItemSwap) RegisterOnSwapItemForEffect(effectID int32, aura *Aura) {
	slots := swap.EligibleSlotsForEffect(effectID)
	character := swap.character
	character.RegisterItemSwapCallback(slots, func(sim *Simulation, _ proto.ItemSlot) {
		procMask := character.GetProcMaskForEnchant(effectID)

		if procMask == ProcMaskUnknown {
			aura.Deactivate(sim)
		} else {
			aura.Activate(sim)
		}
	})
}

func (swap *ItemSwap) IsEnabled() bool {
	return swap.character != nil && len(swap.slots) > 0
}

func (swap *ItemSwap) IsValidSwap(swapSet proto.APLActionItemSwap_SwapSet) bool {
	return swap.swapSet != swapSet
}

func (swap *ItemSwap) IsSwapped() bool {
	return swap.swapSet == proto.APLActionItemSwap_Swap1
}

func (character *Character) hasItemEquipped(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return character.Equipment.containsItemInSlots(itemID, possibleSlots)
}

func (character *Character) hasEnchantEquipped(effectID int32, possibleSlots []proto.ItemSlot) bool {
	return character.Equipment.containsEnchantInSlots(effectID, possibleSlots)
}

func (swap *ItemSwap) GetEquippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.character.Equipment[slot]
}

func (swap *ItemSwap) GetUnequippedItemBySlot(slot proto.ItemSlot) *Item {
	return &swap.unEquippedItems[slot]
}

func (swap *ItemSwap) ItemExistsInMainEquip(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.originalEquip.containsItemInSlots(itemID, possibleSlots)
}

func (swap *ItemSwap) ItemExistsInSwapEquip(itemID int32, possibleSlots []proto.ItemSlot) bool {
	return swap.swapEquip.containsItemInSlots(itemID, possibleSlots)
}

func (swap *ItemSwap) EligibleSlotsForItem(itemID int32) []proto.ItemSlot {
	eligibleSlots := eligibleSlotsForItem(GetItemByID(itemID))

	if len(eligibleSlots) == 0 {
		return []proto.ItemSlot{}
	}

	if !swap.IsEnabled() {
		return eligibleSlots
	} else {
		return FilterSlice(eligibleSlots, func(slot proto.ItemSlot) bool {
			return (swap.originalEquip[slot].ID == itemID) || (swap.swapEquip[slot].ID == itemID)
		})
	}
}

func (swap *ItemSwap) EligibleSlotsForEffect(effectID int32) []proto.ItemSlot {
	var eligibleSlots []proto.ItemSlot

	if swap.IsEnabled() {
		for itemSlot := proto.ItemSlot(0); itemSlot < NumItemSlots; itemSlot++ {
			if swap.originalEquip.containsEnchantInSlot(effectID, itemSlot) || swap.swapEquip.containsEnchantInSlot(effectID, itemSlot) {
				eligibleSlots = append(eligibleSlots, itemSlot)
			}
		}
	}

	return eligibleSlots
}

func (swap *ItemSwap) SwapItems(sim *Simulation, swapSet proto.APLActionItemSwap_SwapSet, isReset bool) {
	if !swap.IsEnabled() || (!swap.IsValidSwap(swapSet) && !isReset) {
		return
	}

	character := swap.character
	weaponSlotSwapped := false
	isPrepull := sim.CurrentTime < 0

	for _, slot := range swap.slots {
		if (slot >= proto.ItemSlot_ItemSlotMainHand) && (slot <= proto.ItemSlot_ItemSlotRanged) {
			weaponSlotSwapped = true
		} else if !isReset && !isPrepull {
			continue
		}

		swap.swapItem(sim, slot, isPrepull, isReset)

		for _, onSwap := range swap.onSwapCallbacks[slot] {
			onSwap(sim, slot)
		}
	}

	if !swap.IsValidSwap(swapSet) {
		return
	}

	statsToSwap := Ternary(isPrepull, swap.equipmentStats.allSlots, swap.equipmentStats.weaponSlots)
	if swap.IsSwapped() {
		statsToSwap = statsToSwap.Invert()
	}

	if sim.Log != nil {
		sim.Log("Item Swap - Stats Change: %v", statsToSwap.FlatString())
	}
	character.AddStatsDynamic(sim, statsToSwap)

	if !isPrepull && !isReset && weaponSlotSwapped {
		character.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime, false)
		character.AutoAttacks.StopRangedUntil(sim, sim.CurrentTime)
		character.SetGCDTimer(sim, max(character.NextGCDAt(), sim.CurrentTime+GCDDefault))
	}

	swap.swapSet = swapSet
}

func (swap *ItemSwap) swapItem(sim *Simulation, slot proto.ItemSlot, isPrepull bool, isReset bool) {
	oldItem := *swap.GetEquippedItemBySlot(slot)

	if isReset {
		swap.character.Equipment[slot] = swap.originalEquip[slot]
	} else {
		swap.character.Equipment[slot] = swap.unEquippedItems[slot]
	}

	swap.unEquippedItems[slot] = oldItem

	if isPrepull {
		return
	}

	character := swap.character

	switch slot {
	case proto.ItemSlot_ItemSlotMainHand:
		if character.AutoAttacks.AutoSwingMelee {
			character.AutoAttacks.SetMH(character.WeaponFromMainHand())
		}
	case proto.ItemSlot_ItemSlotOffHand:
		// OH slot handling is more involved because we need to dynamically toggle the OH weapon attack on/off
		// depending on the updated DW status after the swap.
		if character.AutoAttacks.AutoSwingMelee {
			weapon := character.WeaponFromOffHand()
			character.AutoAttacks.SetOH(weapon)
			character.AutoAttacks.IsDualWielding = weapon.SwingSpeed != 0
			character.AutoAttacks.EnableMeleeSwing(sim)
			character.PseudoStats.CanBlock = character.OffHand().WeaponType == proto.WeaponType_WeaponTypeShield
		}
	case proto.ItemSlot_ItemSlotRanged:
		if character.AutoAttacks.AutoSwingRanged {
			character.AutoAttacks.SetRanged(character.WeaponFromRanged())
		}
	}
}

func (swap *ItemSwap) reset(sim *Simulation) {
	swap.initialized = false
	if !swap.IsEnabled() {
		return
	}

	swap.SwapItems(sim, proto.APLActionItemSwap_Main, true)

	swap.unEquippedItems = swap.swapEquip

	// This is used to set the initial spell flags for unequipped items.
	// Reset is called before the first iteration.
	swap.initialized = true
}

func (swap *ItemSwap) doneIteration(sim *Simulation) {
	swap.reset(sim)
}

func calcItemSwapStatsOffset(originalEquipment Equipment, swapEquipment Equipment, prepullBonusStats stats.Stats, slots []proto.ItemSlot) ItemSwapStats {
	allSlotStats := prepullBonusStats
	weaponSlotStats := stats.Stats{}
	allWeaponSlots := AllWeaponSlots()

	for _, slot := range slots {
		slotStats := ItemEquipmentStats(swapEquipment[slot], false).Subtract(ItemEquipmentStats(originalEquipment[slot], false))
		allSlotStats = allSlotStats.Add(slotStats)

		if slices.Contains(allWeaponSlots, slot) {
			weaponSlotStats = weaponSlotStats.Add(slotStats)
		}
	}

	return ItemSwapStats{
		allSlots:    allSlotStats,
		weaponSlots: weaponSlotStats,
	}
}

func toItem(itemSpec *proto.ItemSpec) Item {
	if itemSpec == nil || itemSpec.Id == 0 {
		return Item{}
	}

	return NewItem(ItemSpec{
		ID:           itemSpec.Id,
		Enchant:      itemSpec.Enchant,
		RandomSuffix: itemSpec.RandomSuffix,
		Rune:         itemSpec.Rune,
	})
}
