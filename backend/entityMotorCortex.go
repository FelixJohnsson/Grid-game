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
                fmt.Println(b.Owner.FullName + "'s brain is not conscious but still alive.")
				continue
            }
            if b.MotorCortexCurrentTask.ActionType == "Walk" && !b.MotorCortexCurrentTask.Finished {
				fmt.Println(Yellow + "Motor cortex is executing the task: " + b.MotorCortexCurrentTask.ActionReason + Reset)

				if b.MotorCortexCurrentTask.TargetLocation.X == b.Owner.Location.X && b.MotorCortexCurrentTask.TargetLocation.Y == b.Owner.Location.Y {
                        fmt.Println("The motor cortex thinks we've arrived at the target location.")
                        b.MotorCortexCurrentTask.Finished = true
                        b.MotorCortexCurrentTask.IsActive = false
						continue
				}
                path := b.DecidePathTo(b.MotorCortexCurrentTask.TargetLocation.X, b.MotorCortexCurrentTask.TargetLocation.Y)
                if path == nil {
                    fmt.Println("I can't find a path to the target location.")
					continue
                } else {
                    b.TakeStepOverPath(b.MotorCortexCurrentTask)
                    fmt.Println("At: ", b.Owner.Location.X, b.Owner.Location.Y)
                }
            }

            // Sleep for 1 seconds
            time.Sleep(1000 * time.Millisecond)
        }
    }
}