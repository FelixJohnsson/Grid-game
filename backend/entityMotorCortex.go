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
                        fmt.Println("The motor cortex thinks we've arrived at the target location.")
                        b.MotorCortexCurrentTask.Finished = true
                        b.MotorCortexCurrentTask.IsActive = false
						continue
				}
                path := b.DecidePathTo(b.MotorCortexCurrentTask.TargetLocation.X, b.MotorCortexCurrentTask.TargetLocation.Y)
                if path == nil {
					continue
                } else {
                    b.TakeStepOverPath(b.MotorCortexCurrentTask)
                }
            }

            // Sleep for 1 seconds
            time.Sleep(250 * time.Millisecond)
        }
    }
}