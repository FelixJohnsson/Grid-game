package main

import (
	"context"
	"errors"
	"time"
)

type worldIntentType string

const (
	worldIntentMoveEntity       worldIntentType = "MoveEntity"
	worldIntentAddItem          worldIntentType = "AddItem"
	worldIntentDestroyItem      worldIntentType = "DestroyItem"
	worldIntentConsumeRipeFruit worldIntentType = "ConsumeRipeFruit"
	worldIntentAddShelter       worldIntentType = "AddShelter"
)

const (
	worldIntentQueueBufferSize = 2048
	worldIntentTickInterval    = 50 * time.Millisecond
	worldIntentSubmitTimeout   = 2 * time.Second
	worldIntentResponseTimeout = 3 * time.Second
)

var (
	ErrWorldIntentSubmitTimeout   = errors.New("intent submit timed out")
	ErrWorldIntentResponseTimeout = errors.New("intent response timed out")
	ErrWorldIntentLoopStopped     = errors.New("intent loop stopped")
	ErrUnknownWorldIntentType     = errors.New("unknown world intent type")
)

type worldIntentRequest struct {
	intentType worldIntentType
	entity     *Entity
	item       *Item
	shelter    *Shelter
	x          int
	y          int
	response   chan worldIntentResult
}

type worldIntentResult struct {
	err   error
	fruit Fruit
}

func (w *World) EnableIntentEngine() {
	w.intentMu.Lock()
	defer w.intentMu.Unlock()

	if w.intentEngineEnabled {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.intentLoopCtx = ctx
	w.intentLoopCancel = cancel
	w.intentQueue = make(chan worldIntentRequest, worldIntentQueueBufferSize)
	w.intentEngineEnabled = true

	go w.runIntentLoop(ctx, w.intentQueue)
}

func (w *World) intentEngineIsEnabled() bool {
	w.intentMu.RLock()
	defer w.intentMu.RUnlock()
	return w.intentEngineEnabled && w.intentQueue != nil
}

func (w *World) submitWorldIntent(request worldIntentRequest) worldIntentResult {
	w.intentMu.RLock()
	enabled := w.intentEngineEnabled
	queue := w.intentQueue
	w.intentMu.RUnlock()

	if !enabled || queue == nil {
		return w.applyWorldIntent(request)
	}

	request.response = make(chan worldIntentResult, 1)

	select {
	case queue <- request:
	case <-time.After(worldIntentSubmitTimeout):
		return worldIntentResult{err: ErrWorldIntentSubmitTimeout}
	}

	select {
	case result := <-request.response:
		return result
	case <-time.After(worldIntentResponseTimeout):
		return worldIntentResult{err: ErrWorldIntentResponseTimeout}
	}
}

func (w *World) runIntentLoop(ctx context.Context, queue <-chan worldIntentRequest) {
	ticker := time.NewTicker(worldIntentTickInterval)
	defer ticker.Stop()

	pending := make([]worldIntentRequest, 0, 64)

	for {
		select {
		case <-ctx.Done():
			w.failPendingIntents(pending, ErrWorldIntentLoopStopped)
			return
		case request := <-queue:
			pending = append(pending, request)
		case <-ticker.C:
			if len(pending) > 0 {
				w.applyIntentBatch(pending)
				pending = pending[:0]
			}
		}

		// Drain queue quickly so each tick applies a full deterministic batch.
		for {
			select {
			case request := <-queue:
				pending = append(pending, request)
			default:
				goto drained
			}
		}
	drained:
	}
}

func (w *World) applyIntentBatch(batch []worldIntentRequest) {
	for _, request := range batch {
		result := w.applyWorldIntent(request)
		if request.response != nil {
			request.response <- result
		}
	}
}

func (w *World) failPendingIntents(batch []worldIntentRequest, err error) {
	for _, request := range batch {
		if request.response != nil {
			request.response <- worldIntentResult{err: err}
		}
	}
}

func (w *World) applyWorldIntent(request worldIntentRequest) worldIntentResult {
	switch request.intentType {
	case worldIntentMoveEntity:
		return worldIntentResult{err: w.moveEntityDirect(request.entity, request.x, request.y)}
	case worldIntentAddItem:
		return worldIntentResult{err: w.addItemDirect(request.x, request.y, request.item)}
	case worldIntentDestroyItem:
		return worldIntentResult{err: w.destroyItemDirect(request.item)}
	case worldIntentConsumeRipeFruit:
		fruit, err := w.consumeRipeFruitAtDirect(request.x, request.y)
		return worldIntentResult{err: err, fruit: fruit}
	case worldIntentAddShelter:
		return worldIntentResult{err: w.addShelterDirect(request.x, request.y, request.shelter)}
	default:
		return worldIntentResult{err: ErrUnknownWorldIntentType}
	}
}
