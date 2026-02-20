package main

// ----------------- Memory Functions -----------------

// AddMemoryToShortTerm adds a memory to the short term memory
func (b *Brain) AddMemoryToShortTerm(event string, details string, location Location) {
	for _, memory := range b.Memories.ShortTermMemory {
		if memory.Event == event && memory.Details == details && memory.Location == location {
			return
		}
	}

	memory := Memory{event, details, location}
	b.Memories.ShortTermMemory = append(b.Memories.ShortTermMemory, memory)
}

// AddMemoryToLongTerm adds a memory to the long term memory
func (b *Brain) AddMemoryToLongTerm(event string, details string, location Location) {
	for _, memory := range b.Memories.LongTermMemory {
		if memory.Event == event && memory.Details == details && memory.Location == location {
			return
		}
	}

	memory := Memory{event, details, location}
	b.Memories.LongTermMemory = append(b.Memories.LongTermMemory, memory)
}

// RemoveMemoriesByEventAtLocation removes matching memories from both short and long term stores.
func (b *Brain) RemoveMemoriesByEventAtLocation(event string, location Location) {
	filter := func(memories []Memory) []Memory {
		filtered := memories[:0]
		for _, memory := range memories {
			if memory.Event == event && memory.Location == location {
				continue
			}
			filtered = append(filtered, memory)
		}
		return filtered
	}

	b.Memories.ShortTermMemory = filter(b.Memories.ShortTermMemory)
	b.Memories.LongTermMemory = filter(b.Memories.LongTermMemory)
}
