package main

import (
	"fmt"
)

// UpdateLocation updates the internal location of the person
func (e *Entity) UpdateLocation(x, y int) {
	e.Location.X = x
	e.Location.Y = y
}

// ---------------- Check Body parts --------------

// HasRightArm returns true if the person has a right arm
func (e *Entity) HasRightArm() bool {
	return e.Body.RightArm != nil
}

// HasLeftArm returns true if the person has a left arm
func (e *Entity) HasLeftArm() bool {
	return e.Body.LeftArm != nil
}

// HasRightLeg returns true if the person has a right leg
func (e *Entity) HasRightLeg() bool {
	return e.Body.RightLeg != nil
}

// HasLeftLeg returns true if the person has a left leg
func (e *Entity) HasLeftLeg() bool {
	return e.Body.LeftLeg != nil
}

// HasRightFoot returns true if the person has a right foot
func (e *Entity) HasRightFoot() bool {
	return e.Body.RightLeg.Foot != nil
}

// HasLeftFoot returns true if the person has a left foot
func (e *Entity) HasLeftFoot() bool {
	return e.Body.LeftLeg.Foot != nil
}

// HasRightHand returns true if the person has a right hand
func (e *Entity) HasRightHand() bool {
	return e.Body.RightArm.Hand != nil
}

// HasLeftHand returns true if the person has a left hand
func (e *Entity) HasLeftHand() bool {
	return e.Body.LeftArm.Hand != nil
}

// HasHead returns true if the person has a head
func (e *Entity) HasHead() bool {
	return e.Body.Head != nil
}

// HasTorso returns true if the person has a torso
func (e *Entity) HasTorso() bool {
	return e.Body.Torso != nil
}

// HasTail returns true if the person has a tail
func (e *Entity) HasTail() bool {
	return e.Body.Tail != nil
}

// ---------------- Grab and Drop ----------------

// Grab in the right hand
func (e *Entity) GrabWithRightHand(item *Item) {
	e.Body.RightArm.Hand.Items = append(e.Body.RightArm.Hand.Items, item)
	e.OwnedItems = append(e.OwnedItems, item)
	// If the item has residues, add them to the limb
	if item.Residues != nil {
		for _, residue := range item.Residues {
			e.AddResidue("RightHand", residue)
		}
	}
}

// Grab in the left hand
func (e *Entity) GrabWithLeftHand(item *Item) {
	e.Body.LeftArm.Hand.Items = append(e.Body.LeftArm.Hand.Items, item)
	e.OwnedItems = append(e.OwnedItems, item)
	// If the item has residues, add them to the limb
	if item.Residues != nil {
		for _, residue := range item.Residues {
			e.AddResidue("RightHand", residue)
		}
	}
}

// Drop from the right hand
func (e *Entity) DropFromRightHand(item string) {
	for i, heldItem := range e.Body.RightArm.Hand.Items {
		if heldItem.Name == item {
			e.Body.RightArm.Hand.Items = append(e.Body.RightArm.Hand.Items[:i], e.Body.RightArm.Hand.Items[i+1:]...)
			e.WorldProvider.AddItem(e.Location.X, e.Location.Y, heldItem)
			e.Brain.AddMemoryToShortTerm("Dropped my " + item, e.FullName, e.Location)
			return
		}
	}
}

// Drop from the left hand
func (e *Entity) DropFromLeftHand(item string) {
	for i, heldItem := range e.Body.LeftArm.Hand.Items {
		if heldItem.Name == item {
			e.Body.RightArm.Hand.Items = append(e.Body.RightArm.Hand.Items[:i], e.Body.RightArm.Hand.Items[i+1:]...)
			e.WorldProvider.AddItem(e.Location.X, e.Location.Y, heldItem)
			e.Brain.AddMemoryToShortTerm("Dropped my " + item, e.FullName, e.Location)
			return
		}
	}
}

// ---------------- Residues & Limbs ----------------

// AddResidue adds a residue to the limb - Check if that limb exists first
func (e *Entity) AddResidue(limb BodyPartType, residue Residue) {
	switch limb {
	case "LeftFoot":
		// Should loop over the residues and add the new residue if it doesn't exist
		for _, r := range e.Body.LeftLeg.Foot.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.LeftLeg.Foot.Residues = append(e.Body.LeftLeg.Foot.Residues, residue)
			}
		}
	case "RightFoot":
		for _, r := range e.Body.RightLeg.Foot.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.RightLeg.Foot.Residues = append(e.Body.RightLeg.Foot.Residues, residue)
			}
		}
	case "Head":
		for _, r := range e.Body.Head.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.Head.Residues = append(e.Body.Head.Residues, residue)
			}
		}
	case "Torso":
		for _, r := range e.Body.Torso.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.Torso.Residues = append(e.Body.Torso.Residues, residue)
			}
		}
	case "RightLeg":
		for _, r := range e.Body.RightLeg.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.RightLeg.Residues = append(e.Body.RightLeg.Residues, residue)
			}
		}
	case "LeftLeg":
		for _, r := range e.Body.LeftLeg.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.LeftLeg.Residues = append(e.Body.LeftLeg.Residues, residue)
			}
		}
	case "RightHand":
		for _, r := range e.Body.RightArm.Hand.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.RightArm.Hand.Residues = append(e.Body.RightArm.Hand.Residues, residue)
			}
		}
	case "LeftHand":
		for _, r := range e.Body.LeftArm.Hand.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				e.Body.LeftArm.Hand.Residues = append(e.Body.LeftArm.Hand.Residues, residue)
			}
		}
	}
}

//RemoveLimb removes a limb from the entity
func (e *Entity) RemoveLimb(limb BodyPartType) {
	fmt.Println(limb, "has been SEVERED!!!")

	switch limb {
	case "Head":
		e.Brain.turnOff()
		e.Body.Head = nil
		return
	case "RightHand":
		e.Body.RightArm.Hand = nil 
		return
	case "LeftHand":
		e.Body.LeftArm.Hand = nil
		return
	case "RightFoot":
		e.Body.RightLeg.Foot = nil
		e.IsIncapacitated = true
		return
	case "LeftFoot":
		e.Body.LeftLeg.Foot = nil
		e.IsIncapacitated = true
		return
	case "RightLeg":
		e.Body.RightLeg = nil
		e.IsIncapacitated = true
		return
	case "LeftLeg":
		e.Body.LeftLeg = nil
		e.IsIncapacitated = true
		return
	}
}

// ---------------- Traversing ----------------

// WalkTo - Walk to a location - This assumes that the person can physically walk and that it's possible to walk to the location. This should be one tile away, so one stee.
func (e *Entity) WalkStepTo(x, y int) {
	requiredLimbs := []BodyPartType{"RightLeg", "LeftLeg"}
	e.IsMoving = TargetedAction{"Walk", string(x) + ", " + string(y), true, requiredLimbs, 10}
	e.WorldProvider.MoveEntity(e, x, y)
}

