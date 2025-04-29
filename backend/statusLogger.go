package main

import (
	"fmt"
)

const CONSOLE_LOGGER_ENABLED = false

func (b *Brain) StatusLogger() {
    entity := b.Owner
    
    // Log basic info to the InfoPanel
    if GlobalInfoPanel != nil {
        if b.PainLevel > b.PainTolerance/2 {
            LogHealth(fmt.Sprintf("Pain level critical: %d/%d", b.PainLevel, b.PainTolerance), entity.FullName)
        }
        
        if b.PhysiologicalNeeds.Hunger > 70 {
            LogFood(fmt.Sprintf("Hunger critical: %d", b.PhysiologicalNeeds.Hunger), entity.FullName)
        } else if b.PhysiologicalNeeds.Hunger > 30 {
            LogFood(fmt.Sprintf("Getting hungry: %d", b.PhysiologicalNeeds.Hunger), entity.FullName)
        }
        
        if b.PhysiologicalNeeds.Thirst > 70 {
            LogWater(fmt.Sprintf("Thirst critical: %d", b.PhysiologicalNeeds.Thirst), entity.FullName)
        } else if b.PhysiologicalNeeds.Thirst > 30 {
            LogWater(fmt.Sprintf("Getting thirsty: %d", b.PhysiologicalNeeds.Thirst), entity.FullName)
        }
        
        if entity.IsBleeding {
            LogHealth("Bleeding", entity.FullName)
        }
        
        if !b.IsConscious {
            LogHealth("Unconscious", entity.FullName)
        }
        
        if b.CurrentTask.IsActive {
            LogInfo(fmt.Sprintf("Task: %s → %s", b.CurrentTask.Action, b.CurrentTask.Target), entity.FullName)
        }
    }
    
    if CONSOLE_LOGGER_ENABLED {
        // Original console output
        fmt.Println("\n" + Green + "=== STATUS REPORT: " + entity.FullName + " ===" + Reset)
        fmt.Println(Cyan + "Basic Info:" + Reset)
        fmt.Printf("  Age: %d | Gender: %s | Species: %s | Occupation: %s\n", 
            entity.Age, entity.Gender, entity.Species, entity.Occupation)
		fmt.Printf("  Location: X:%d, Y:%d\n", entity.Location.X, entity.Location.Y)
		
		fmt.Println(Cyan + "Physical Status:" + Reset)
		healthStatus := Green
		if b.PainLevel > b.PainTolerance/2 {
			healthStatus = Yellow
		}
		if b.PainLevel > b.PainTolerance*3/4 {
			healthStatus = Red
		}
		fmt.Printf("  Conscious: %t | Pain Level: %s%d/%d%s | Oxygen: %d | Bleeding: %t\n", 
			b.IsConscious, healthStatus, b.PainLevel, b.PainTolerance, Reset, b.OxygenLevel, entity.IsBleeding)
		fmt.Printf("  Incapacitated: %t | Can Breathe: %t | Brain Damage: %d\n", 
			entity.IsIncapacitated, b.CanBreath, b.BrainDamage)
		
		fmt.Println(Cyan + "Attributes:" + Reset)
		fmt.Printf("  STR: %d | AGI: %d | INT: %d | CHA: %d | STA: %d | Combat Skill: %d\n", 
			entity.Strength, entity.Agility, entity.Intelligence, entity.Charisma, entity.Stamina, entity.CombatSkill)
		
		fmt.Println(Cyan + "Needs:" + Reset)
		
		// Color-code hunger level
		hungerColor := Green
		if b.PhysiologicalNeeds.Hunger > 30 {
			hungerColor = Yellow
		} 
		if b.PhysiologicalNeeds.Hunger > 70 {
			hungerColor = Red
		}
		
		// Color-code thirst level
		thirstColor := Green
		if b.PhysiologicalNeeds.Thirst > 30 {
			thirstColor = Yellow
		}
		if b.PhysiologicalNeeds.Thirst > 70 {
			thirstColor = Red
		}
		
		fmt.Printf("  Hunger: %s%d%s | Thirst: %s%d%s | Rest: %d\n", 
			hungerColor, b.PhysiologicalNeeds.Hunger, Reset,
			thirstColor, b.PhysiologicalNeeds.Thirst, Reset,
			b.PhysiologicalNeeds.Rested)
		fmt.Printf("  Safe: %d | Scared: %d | In Safe Area: %t\n", 
			entity.FeelingSafe, entity.FeelingScared, b.PhysiologicalNeeds.IsInSafeArea)
		
		fmt.Println(Cyan + "Current State:" + Reset)
		if b.CurrentTask.IsActive {
			fmt.Printf("  Current Task: %s%s → %s%s (Priority: %d)\n", 
				Yellow, b.CurrentTask.Action, b.CurrentTask.Target, Reset, b.CurrentTask.Priority)
			
			fmt.Println(Cyan + "Current Action Details:" + Reset)
			fmt.Printf("    Action: %s%s → %s%s (Priority: %d)\n", 
				Yellow, b.CurrentTask.Action, b.CurrentTask.Target, Reset, b.CurrentTask.Priority)
		} else {
			fmt.Println("  No active task")
		}
		
		fmt.Println("  Thinking: " + Magenta + entity.Thinking + Reset)
		
		if len(entity.WantsTo) > 0 {
			fmt.Println(Cyan + "Wants:" + Reset)
			for _, want := range entity.WantsTo {
				fmt.Println("  • " + Blue + want + Reset)
			}
		}
		
		if len(entity.Relationships) > 0 {
			fmt.Println(Cyan + "Relationships:" + Reset)
			for _, rel := range entity.Relationships {
				relationshipColor := Green
				if rel.Intensity < 0 {
					relationshipColor = Red
				}
				fmt.Printf("  • %s: %s%s (%d)%s\n", 
					rel.WithEntity, relationshipColor, rel.Relationship, rel.Intensity, Reset)
			}
		}
		
		if len(entity.OwnedItems) > 0 {
			fmt.Println(Cyan + "Owned Items:" + Reset)
			for _, item := range entity.OwnedItems {
				fmt.Printf("  • %s\n", item.Name)
			}
		}
		
		fmt.Println(Cyan + "Body Status:" + Reset)
		if entity.Body.Head != nil {
			headStatus := Green
			if entity.Body.Head.IsBleeding || entity.Body.Head.IsBroken {
				headStatus = Red
			}
			fmt.Printf("  Head: %sBleeding=%t, Broken=%t%s\n", 
				headStatus, entity.Body.Head.IsBleeding, entity.Body.Head.IsBroken, Reset)
		}
		if entity.Body.Torso != nil {
			torsoStatus := Green
			if entity.Body.Torso.IsBleeding || entity.Body.Torso.IsBroken {
				torsoStatus = Red
			}
			fmt.Printf("  Torso: %sBleeding=%t, Broken=%t%s\n", 
				torsoStatus, entity.Body.Torso.IsBleeding, entity.Body.Torso.IsBroken, Reset)
		}
		if entity.Body.RightArm != nil {
			armStatus := Green
			if entity.Body.RightArm.IsBleeding || entity.Body.RightArm.IsBroken {
				armStatus = Red
			}
			fmt.Printf("  Right Arm: %sBleeding=%t, Broken=%t%s\n", 
				armStatus, entity.Body.RightArm.IsBleeding, entity.Body.RightArm.IsBroken, Reset)
			if entity.Body.RightArm.Hand != nil {
				holding := ""
				if len(entity.Body.RightArm.Hand.Items) > 0 {
					for _, item := range entity.Body.RightArm.Hand.Items {
						holding += item.Name + " "
					}
				} else {
					holding = "nothing"
				}
				fmt.Printf("    Hand: Holding=%s\n", holding)
			}
		}
		if entity.Body.LeftArm != nil {
			armStatus := Green
			if entity.Body.LeftArm.IsBleeding || entity.Body.LeftArm.IsBroken {
				armStatus = Red
			}
			fmt.Printf("  Left Arm: %sBleeding=%t, Broken=%t%s\n", 
				armStatus, entity.Body.LeftArm.IsBleeding, entity.Body.LeftArm.IsBroken, Reset)
			if entity.Body.LeftArm.Hand != nil {
				holding := ""
				if len(entity.Body.LeftArm.Hand.Items) > 0 {
					for _, item := range entity.Body.LeftArm.Hand.Items {
						holding += item.Name + " "
					}
				} else {
					holding = "nothing"
				}
				fmt.Printf("    Hand: Holding=%s\n", holding)
			}
        }
        
        fmt.Println(Green + "=============================" + Reset + "\n")
    }
}