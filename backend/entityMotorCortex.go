package main

import (
	"fmt"
	"time"
)

// MotorCortex is handling all the task that require motor function, for example walking. This is it's own task list.
func (b *Brain) MotorCortex() {
    for {
        select {
        case <-b.Ctx.Done():
            return
        default:
            if !b.Active {
                continue
            }

            if !b.IsConscious {
				continue
            }
            
            if b.MotorCortexCurrentTask.ActionType == "Walk" && !b.MotorCortexCurrentTask.Finished {
				if b.MotorCortexCurrentTask.TargetLocation.X == b.Owner.Location.X && b.MotorCortexCurrentTask.TargetLocation.Y == b.Owner.Location.Y {
                        b.FinishMotorCortexTask()
						return
				}
                path := b.DecidePathTo(b.MotorCortexCurrentTask.TargetLocation.X, b.MotorCortexCurrentTask.TargetLocation.Y)
                if path == nil {
                    fmt.Println(Yellow + "Motor cortex: I can't find a path to the target location.", b.MotorCortexCurrentTask.TargetLocation.X, b.MotorCortexCurrentTask.TargetLocation.Y)
                    fmt.Print(Reset)
                    return
                } else {
                    fmt.Println(Yellow + "I will take a step over the path.", b.MotorCortexCurrentTask.TargetLocation.X, b.MotorCortexCurrentTask.TargetLocation.Y)
                    fmt.Print(Reset)
                    b.TakeStepOverPath(b.MotorCortexCurrentTask)
                }
            } else {
                fmt.Println("Action type: ", b.MotorCortexCurrentTask.ActionType, "Finished", b.MotorCortexCurrentTask.Finished)
            }
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func (b *Brain) AddMotorCortexTask(reason string, ActionType string, targetLocation Location) {
    fmt.Println(Yellow + "Adding a motor cortex task: ", reason, ActionType, targetLocation)
    fmt.Print(Reset)
    b.MotorCortexCurrentTask = MotorCortexAction{reason, ActionType, targetLocation, false, false}
}

func (b *Brain) FinishMotorCortexTask() {
    fmt.Println(Yellow + "Finishing motor cortex task.", b.MotorCortexCurrentTask.ActionType)
    fmt.Print(Reset)
    b.MotorCortexCurrentTask.Finished = true
    b.MotorCortexCurrentTask.IsActive = false
}