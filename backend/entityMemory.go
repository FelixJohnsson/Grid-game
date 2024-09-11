package main

// ----------------- Memory Functions -----------------

// AddMemoryToShortTerm adds a memory to the short term memory
func (b *Brain) AddMemoryToShortTerm(event string, details string, location Location) {
    memory := Memory{event, details, location}
    b.Memories.ShortTermMemory = append(b.Memories.ShortTermMemory, memory)
}

// AddMemoryToLongTerm adds a memory to the long term memory
func (b *Brain) AddMemoryToLongTerm(event string, details string, location Location) {
    memory := Memory{event, details, location}
    b.Memories.LongTermMemory = append(b.Memories.LongTermMemory, memory)
}