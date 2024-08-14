package cards

func defaultPlayer(hero, heropower *Card, deck []*Card) *Player {
	return &Player{
		Image:           "",
		Hero:            hero,
		HeroPower:       heropower,
		UsedHeroPower:   false,
		Weapon:          nil,
		MaxMinions:      7,
		Minions:         []*Card{},
		MaxHand:         10,
		Hand:            []*Card{},
		MaxDeck:         60,
		Deck:            deck,
		Graveyard:       []*Card{},
		DiscardPile:     []*Card{},
		Crematorium:     []*Card{},
		MaxSecrets:      5,
		Secrets:         []*Card{},
		MaxMaxMana:      10,
		MaxMana:         0,
		Mana:            0,
		OverchargedMana: 0,
	}
}

func WarlockMurloc() *Player {
	return defaultPlayer(GuldanHero(), DrawCardFor2Life(),
		[]*Card{
			RockBottom(),
			RockBottom(),
			AzsharanScavenger(),
			AzsharanScavenger(),
			ChumBucket(),
			ChumBucket(),
			Voidgill(),
			Voidgill(),
			BloodscentVilefin(),
			BloodscentVilefin(),
			SeadevilStinger(),
			SeadevilStinger(),
			Gigafin(),
			MurlocTinyfin(),
			MurlocTinyfin(),
			Murmy(),
			Murmy(),
			BluegillWarrior(),
			BluegillWarrior(),
			Crabrider(),
			Crabrider(),
			LushwaterScout(),
			LushwaterScout(),
			MurlocWarleader(),
			MurlocWarleader(),
			TwinfinFinTwin(),
			TwinfinFinTwin(),
			OldMurkEye(),
			GorlocRavager(),
			GorlocRavager(),
		},
	)
}

func PaladinControlZoth() *Player {
	return defaultPlayer(UtherHero(), SummonRecruiter(),
		[]*Card{
			Redemption(),
			Redemption(),
			Equality(),
			Equality(),
			AldorPeacekeeper(),
			AldorPeacekeeper(),
			Consecration(),
			Consecration(),
			DivineFavor(),
			DivineFavor(),
			SwordOfJustice(),
			BlessingOfKings(),
			KeeperOfUldaman(),
			KeeperOfUldaman(),
			StandAgainstDarkness(),
			StandAgainstDarkness(),
			TruesilverChampion(),
			TruesilverChampion(),
			LayOnHands(),
			TirionFordring(),
			BloodmageThalnos(),
			BloodmageThalnos(), //remove
			LootHoarder(),
			LootHoarder(),
			SpawnOfNZoth(),
			SludgeBelcher(),
			EmperorThaurissan(),
			SylvanasWindrunner(),
			RagnarosTheFirelord(),
			NZothTheCorruptor(),
		},
	)
}
