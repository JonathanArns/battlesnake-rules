package rules

import (
	"errors"
)

var royaleRulesetStages = []string{
	StageGameOverStandard,
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
	StageSpawnHazardsShrinkMap,
}

type RoyaleRuleset struct {
	StandardRuleset

	ShrinkEveryNTurns int
}

func (r *RoyaleRuleset) Name() string { return GameTypeRoyale }

func (r RoyaleRuleset) Execute(bs *BoardState, s Settings, sm []SnakeMove) (bool, *BoardState, error) {
	return NewPipeline(royaleRulesetStages...).Execute(bs, s, sm)
}

func (r *RoyaleRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	if r.StandardRuleset.HazardDamagePerTurn < 1 {
		return nil, errors.New("royale damage per turn must be greater than zero")
	}
	_, nextState, err := r.Execute(prevState, r.Settings(), moves)
	return nextState, err
}

func PopulateHazardsRoyale(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	b.Hazards = []Point{}

	// Royale uses the current turn to generate hazards, not the previous turn that's in the board state
	turn := b.Turn + 1

	if settings.RoyaleSettings.ShrinkEveryNTurns < 1 {
		return false, errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < settings.RoyaleSettings.ShrinkEveryNTurns {
		return false, nil
	}

	randGenerator := settings.GetRand(0)

	numShrinks := turn / settings.RoyaleSettings.ShrinkEveryNTurns
	minX, maxX := 0, b.Width-1
	minY, maxY := 0, b.Height-1
	for i := 0; i < numShrinks; i++ {
		switch randGenerator.Intn(4) {
		case 0:
			minX += 1
		case 1:
			maxX -= 1
		case 2:
			minY += 1
		case 3:
			maxY -= 1
		}
	}

	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				b.Hazards = append(b.Hazards, Point{x, y})
			}
		}
	}

	return false, nil
}

func (r *RoyaleRuleset) IsGameOver(b *BoardState) (bool, error) {
	return GameOverStandard(b, r.Settings(), nil)
}

func (r RoyaleRuleset) Settings() Settings {
	s := r.StandardRuleset.Settings()
	s.RoyaleSettings = RoyaleSettings{
		ShrinkEveryNTurns: r.ShrinkEveryNTurns,
	}
	return s
}
