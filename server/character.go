package main

import gamepb "server/pkg/grpc"

type Position struct {
	x float32
	y float32
	z float32
}

type Character struct {
	userID        string
	position      Position
	rotationZ     float32
	syncTimestamp float32
}

type CharacterList struct {
	characters map[string]*Character
}

func (c *CharacterList) GetPbCharactersExceptSelf(userID string) []*gamepb.Character {
	pbCharacters := make([]*gamepb.Character, 0, len(c.characters))
	for _, character := range c.characters {
		if character.userID == userID {
			continue
		}
		pbCharacters = append(pbCharacters, &gamepb.Character{
			UserID:    character.userID,
			PositionX: character.position.x,
			PositionY: character.position.y,
			PositionZ: character.position.z,
			RotationZ: character.rotationZ,
			Timestamp: character.syncTimestamp,
		})
	}
	return pbCharacters
}

func (c *CharacterList) UpdateCharacter(pbc *gamepb.Character) {
	if c.characters == nil {
		c.characters = make(map[string]*Character)
	}
	character, exists := c.characters[pbc.GetUserID()]
	if !exists {
		character = &Character{
			userID: pbc.GetUserID(),
			position: Position{
				x: pbc.GetPositionX(),
				y: pbc.GetPositionY(),
				z: pbc.GetPositionZ(),
			},
			rotationZ:     pbc.GetRotationZ(),
			syncTimestamp: pbc.GetTimestamp(),
		}

		c.characters[pbc.GetUserID()] = character
		return
	}

	character.position = Position{
		x: pbc.GetPositionX(),
		y: pbc.GetPositionY(),
		z: pbc.GetPositionZ(),
	}
	character.rotationZ = pbc.GetRotationZ()
	character.syncTimestamp = pbc.GetTimestamp()
}