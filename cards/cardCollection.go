package cards

import (
	"cards/utils"

	"math/rand"

	"github.com/google/uuid"
)

/*
 * Card creation rules:
 * - The first tag in the card is the type of the card (Spell, Minion, etc).
 */

func TheCoin() *Card {
	return &Card{
		Mana:   0,
		Name:   "硬币",
		Rarity: Basic,
		Text:   "在本回合中，获得一个法力水晶。",
		Image:  "/imgs/TheCoin.png",
		Tags:   []string{Spell, Neutral},
		Events: map[string]Event{
			EventSpellCast: GainXManaCrystalEvent(1),
		},
	}
}

func ElvenArcher() *Card {
	return &Card{
		Mana:      1,
		Name:      "精灵弓箭手",
		Attack:    1,
		MaxHealth: 1,
		Rarity:    Basic,
		Text:      "<b>战吼：</b>造成1点伤害。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/2f/Elven_Archer_full.jpg",
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventBattlecry: DealXDamageEvent(1),
		},
		Targets: func(b *Board) []string {
			options := []string{}
			for _, c := range b.AllCharacters() {
				options = append(options, c.Id)
			}
			return options
		},
	}
}

func RockBottom() *Card {
	return &Card{
		Mana:   1,
		Name:   "岩石海底",
		Rarity: Rare,
		Text:   "召唤一个1/1的鱼人，然后<b>探底</b>。如果选中的是鱼人牌，则再召唤一个。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/74/Rock_Bottom_full.jpg",
		Tags:   []string{Spell, Warlock},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				p := ctx.Source.Player
				murloc := RockBottomMurloc()
				murloc.Player = p
				ctx.Board.SummonMinion(ctx.This, murloc)
				card := ctx.Board.Dredge(p)
				if card == nil {
					return nil
				}
				if card.Tribe == Murloc {
					murloc := RockBottomMurloc()
					ctx.Board.SummonMinion(ctx.This, murloc)
				}
				return nil
			},
		},
	}
}

func AzsharanScavenger() *Card {
	return &Card{
		Mana:      2,
		Name:      "艾萨拉的拾荒者",
		Attack:    2,
		MaxHealth: 3,
		Rarity:    Common,
		Text:      "<b>战吼：</b> 将一张沉没的拾荒者置于你的牌库底。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f8/Azsharan_Scavenger_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				card := SunkenScavenger()
				card.Player = ctx.This.Player
				p := ctx.This.Player
				// TODO: create method to add to deck and assign player
				p.Deck = append(p.Deck, card)
				return nil
			},
		},
	}
}

func ChumBucket() *Card {
	return &Card{
		Mana:   2,
		Name:   "鱼饵桶",
		Rarity: Epic,
		Text:   "使你手牌中的所有鱼人牌获得+1/+1。你每控制一个鱼人，重复一次。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/61/Chum_Bucket_full.jpg",
		Tags:   []string{Spell, Warlock},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				p := ctx.This.Player
				for _, c := range p.Hand {
					if c.Tribe == Murloc {
						ctx.Target = c
						ctx.Board.EnchantCard(c, &Enchantment{
							Name:   "+1/+1",
							Attack: 1,
							Health: 1,
						})
					}
				}
				for _, m := range p.Minions {
					if m.Tribe == Murloc {
						for _, c := range p.Hand {
							if c.Tribe == Murloc {
								ctx.Target = c
								ctx.Board.EnchantCard(c, &Enchantment{
									Name:   "+1/+1",
									Attack: 1,
									Health: 1,
								})
							}
						}
					}
				}
				return nil
			},
		},
	}
}

func Voidgill() *Card {
	return &Card{
		Mana:      2,
		Name:      "虚鳃鱼人",
		Attack:    3,
		MaxHealth: 2,
		Rarity:    Rare,
		Text:      "<b>亡语：</b> 使你手牌中的所有鱼人牌获得+1/+1",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/e3/Voidgill_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				for _, c := range ctx.Target.Player.Hand {
					if c.Tribe == Murloc {
						ctx.Board.EnchantCard(c, &Enchantment{
							Name:   "+1/+1",
							Attack: 1,
							Health: 1,
						})
					}
				}
				return nil
			},
		},
	}
}

func BloodscentVilefin() *Card {
	return &Card{
		Mana:      3,
		Name:      "血腥恶鳍鱼人",
		Attack:    3,
		MaxHealth: 4,
		Rarity:    Rare,
		Text:      "<b>战吼：探底。</b>如果选中的是鱼人牌，则使其改为消耗生命值，而非法力值。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/9/9a/Bloodscent_Vilefin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				p, b := ctx.This.Player, ctx.Board
				card := ctx.Board.Dredge(p)
				if card != nil && card.Tags[0] == Minion && card.Tribe == Murloc {
					b.EnchantCard(card, &Enchantment{
						Tags: []string{BloodPayment},
					})
				}
				return nil
			},
		},
	}
}

func SeadevilStinger() *Card {
	return &Card{
		Mana:      4,
		Name:      "海魔钉刺者",
		Attack:    4,
		MaxHealth: 2,
		Rarity:    Rare,
		Text:      "<b>战吼：</b>在本回合中，你使用的下一张鱼人牌不再消耗法力值，转而消耗生命值。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/e8/Seadevil_Stinger_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock, Battlecry},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				h, b := ctx.This.Player.Hero, ctx.Board

				blpmid := uuid.NewString()
				for _, c := range h.Player.Hand {
					b.EnchantCard(c, &Enchantment{
						Id:   blpmid,
						Tags: []string{BloodPayment},
					})
				}
				enchid := uuid.NewString()
				b.EnchantCard(h, &Enchantment{
					Id: enchid,
					Events: map[string]Event{
						EventAfterAddToHand: func(ctx *EventContext) error {
							if ctx.Source.Player != ctx.This.Player {
								return nil
							}
							c := ctx.Target
							b.EnchantCard(c, &Enchantment{
								Id:   blpmid,
								Tags: []string{BloodPayment},
							})
							return nil
						},
						EventAfterSummon: func(ctx *EventContext) error {
							if ctx.Target.Player != ctx.This.Player {
								return nil
							}
							if ctx.Target.Tribe == Murloc {
								for _, c := range h.Player.Hand {
									c.DelEnchId(blpmid)
								}
								ctx.This.DelEnchId(enchid)
								return nil
							}
							return nil
						},
						EventEndOfTurn: func(ctx *EventContext) error {
							for _, c := range h.Player.Hand {
								c.DelEnchId(blpmid)
							}
							ctx.This.DelEnchId(enchid)
							return nil
						},
					},
				})
				return nil
			},
		},
	}
}

func Gigafin() *Card {
	return &Card{
		Mana:      8,
		Name:      "老巨鳍",
		Attack:    7,
		MaxHealth: 4,
		Rarity:    Legendary,
		Text:      "<b>巨型+1</br>战吼：</b>吞食所有敌方随从。亡语：吐出来。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/2/26/Gigafin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Warlock, Deathrattle},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				p := ctx.Source.Player
				opp := ctx.Board.getOpponent(p)
				oppminions := append([]*Card{}, opp.Minions...)
				opp.Minions = []*Card{}
				enchId := uuid.NewString()
				ench := &Enchantment{
					Id: enchId,
					Events: map[string]Event{
						EventDeathrattle: func(ctx *EventContext) error {
							for _, c := range oppminions {
								ctx.Board.SummonMinion(ctx.This, c)
							}
							return nil
						},
					},
				}
				ctx.Board.EnchantCard(ctx.Source, ench)
				maw := GenerateGigafinMaw(ctx.Source.Id, enchId)
				ctx.Board.SummonMinion(ctx.This, maw)
				return nil
			},
		},
	}
}

func MurlocTinyfin() *Card {
	return &Card{
		Mana:      0,
		Name:      "鱼人宝宝",
		Attack:    1,
		MaxHealth: 1,
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f5/Murloc_Tinyfin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
	}
}

func Murmy() *Card {
	return &Card{
		Mana:      1,
		Name:      "鱼人木乃伊",
		Attack:    1,
		MaxHealth: 1,
		Text:      "<b>复生</b>",
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/ad/Murmy_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Reborn},
	}
}

func BluegillWarrior() *Card {
	return &Card{
		Mana:      2,
		Name:      "蓝腮战士",
		Attack:    2,
		MaxHealth: 1,
		Text:      "<b>冲锋</b>",
		Rarity:    Basic,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f2/Bluegill_Warrior_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Charge},
	}
}

func Crabrider() *Card {
	return &Card{
		Mana:      2,
		Name:      "螃蟹骑士",
		Attack:    1,
		MaxHealth: 4,
		Text:      "<b>突袭，风怒。</b>",
		Rarity:    Common,
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/5/59/Crabrider_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Rush, Windfury},
	}
}

func LushwaterScout() *Card {
	return &Card{
		Mana:      2,
		Name:      "甜水鱼人斥候",
		Attack:    1,
		MaxHealth: 3,
		Rarity:    Common,
		Text:      "在你召唤一个鱼人后，使其获得+1攻击力和<b>突袭</b>。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/6b/Lushwater_Scout_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventAfterSummon: func(ctx *EventContext) error {
				card := ctx.Target
				if card.Player != ctx.This.Player || card == ctx.This || card.Tribe != Murloc {
					return nil
				}
				ench := &Enchantment{
					Attack: 1,
					Tags:   []string{Rush},
				}
				ctx.Board.EnchantCard(card, ench)
				return nil
			},
		},
	}
}

func MurlocWarleader() *Card {
	return &Card{
		Mana:      3,
		Name:      "鱼人领军",
		Attack:    3,
		MaxHealth: 3,
		Rarity:    Epic,
		Text:      "你的其他鱼人拥有+2攻击力。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/8/82/Murloc_Warleader_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventEndOfAction: func(ctx *EventContext) error {
				for _, c := range ctx.This.Player.Minions {
					if c.Id != ctx.This.Id && c.Tribe == Murloc && c.GetEnch(ctx.This.Id) == nil {
						ctx.Board.EnchantCard(c, &Enchantment{
							Id:     ctx.This.Id,
							Attack: 2,
						})
					}
				}
				return nil
			},
			EventDestroyMinion: func(ctx *EventContext) error {
				for _, c := range ctx.This.Player.Minions {
					if c.GetEnch(ctx.This.Id) != nil {
						c.DelEnchId(ctx.This.Id)
					}
				}
				return nil
			},
		},
	}
}

func TwinfinFinTwin() *Card {
	return &Card{
		Mana:      3,
		Name:      "并鳍奇兵",
		Attack:    2,
		MaxHealth: 1,
		Rarity:    Rare,
		Text:      "<b>突袭。战吼：</b>召唤一个本随从的复制。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/4/4d/Twin-fin_Fin_Twin_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Rush},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				copy := *ctx.This
				copyptr := &copy
				copyptr.Id = uuid.NewString()
				ctx.Board.SummonMinion(ctx.This, copyptr)
				return nil
			},
		},
	}
}

func OldMurkEye() *Card {
	return &Card{
		Mana:      4,
		Name:      "老瞎眼",
		Attack:    2,
		MaxHealth: 4,
		Rarity:    Legendary,
		Text:      "<b>冲锋</b>，在战场上每有一个其他鱼人便拥有+1攻击力。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/62/Murloc_Raid_art.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral, Charge},
		Events: map[string]Event{
			EventEndOfAction: func(ctx *EventContext) error {
				murlocCount := 0
				for _, c := range append(ctx.This.Player.Minions, ctx.Board.getOpponent(ctx.This.Player).Minions...) {
					if c.Tribe == Murloc && c.Id != ctx.This.Id {
						murlocCount++
					}
				}
				ench := ctx.This.GetEnch(ctx.This.Id)
				if ench == nil {
					ctx.This.AddEnchantment(&Enchantment{
						Id:     ctx.This.Id,
						Attack: murlocCount,
					})
					return nil
				}
				ench.Attack = murlocCount
				return nil
			},
		},
	}
}

func GorlocRavager() *Card {
	return &Card{
		Mana:      5,
		Name:      "鳄鱼人掠夺者",
		Attack:    4,
		MaxHealth: 3,
		Rarity:    Common,
		Text:      "<b>战吼：</b>抽三张鱼人牌。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/a6/Gorloc_Ravager_full.jpg",
		Tribe:     Murloc,
		Tags:      []string{Minion, Neutral},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				count := 3
				p := ctx.Source.Player
				for i, card := range p.Deck {
					if card.Tribe == Murloc {
						ctx.Board.DrawCard(ctx.This, p, byte(i))
						count--
					}
					if count <= 0 {
						break
					}
				}
				return nil
			},
		},
	}
}

func Redemption() *Card {
	return &Card{
		Mana:   1,
		Name:   "救赎",
		Rarity: Common,
		Text:   "<b>奥秘：</b>当一个友方随从死亡时，使其回到战场，并具有1点生命值。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/a/a8/Redemption_full.jpg",
		Tags:   []string{Spell, Paladin, Secret},
		Events: map[string]Event{
			EventAfterDestroyMinion: func(ctx *EventContext) error {
				if ctx.Target.Player == ctx.This.Player {
					ctx.Target.Health = 1
					ctx.Board.SummonMinion(ctx.This, ctx.Target)
				}
				return nil
			},
		},
	}
}

func Equality() *Card {
	return &Card{
		Mana:   2,
		Name:   "生而平等",
		Rarity: Rare,
		Text:   "将所有随从的生命值变为1。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/b/bf/Equality_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				for _, c := range append(ctx.This.Player.Minions, ctx.Board.getOpponent(ctx.This.Player).Minions...) {
					ctx.Board.EnchantCard(c, &Enchantment{
						Id:     ctx.This.Id,
						Health: -ctx.Target.GetMaxHealth() + 1,
					})
					c.Health = 1
				}
				return nil
			},
		},
	}
}

func AldorPeacekeeper() *Card {
	return &Card{
		Mana:      3,
		Name:      "奥尔多卫士",
		Attack:    3,
		MaxHealth: 3,
		Rarity:    Rare,
		Text:      "<b>战吼：</b>使一个敌方随从的攻击力变为1。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/0/0a/Aldor_Peacekeeper_full.jpg",
		Tags:      []string{Minion, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllMinionCards() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				ctx.Board.EnchantCard(ctx.Target, &Enchantment{
					Id:     ctx.This.Id,
					Attack: -ctx.Target.GetAttack() + 1,
				})
				return nil
			},
		},
	}
}

func Consecration() *Card {
	return &Card{
		Mana:   4,
		Name:   "奉献",
		Rarity: Basic,
		Text:   "对所有敌人造成2点伤害。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/b/b4/Consecration_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				op := ctx.Board.getOpponent(ctx.This.Player)
				ctx.Board.DealDamage(ctx.This, op.Hero, 2)
				// create a new array for opp minions
				minions := append([]*Card{}, op.Minions...)
				// because if loop through op.Minions, the slice
				// will be updated as minions are destroyed
				for _, c := range minions {
					ctx.Board.DealDamage(ctx.This, c, 2)
				}
				return nil
			},
		},
	}
}

func DivineFavor() *Card {
	return &Card{
		Mana:   3,
		Name:   "神恩术",
		Rarity: Rare,
		Text:   "抽若干数量的牌，直到你的手牌数量等同于你对手的手牌数量。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/1/19/Divine_Favor_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				cardDiff := len(ctx.Board.getOpponent(ctx.This.Player).Hand) - len(ctx.This.Player.Hand)
				for i := 0; i < cardDiff; i++ {
					ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				}
				return nil
			},
		},
	}
}

func SwordOfJustice() *Card {
	return &Card{
		Mana:      3,
		Name:      "公正之剑",
		Rarity:    Epic,
		Attack:    1,
		MaxHealth: 5,
		Text:      "每当你召唤一个随从，使它获得+1/+1，这把武器失去1点耐久度。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/4/4e/Sword_of_Justice_full.jpg",
		Tags:      []string{Weapon, Paladin},
		Events: map[string]Event{
			EventAfterSummon: func(ctx *EventContext) error {
				if ctx.Target.Player != ctx.This.Player {
					return nil
				}
				ctx.Target.AddEnchantment(&Enchantment{
					Id:     ctx.This.Id,
					Attack: 1,
					Health: 1,
				})
				ctx.Board.LoseDurability(ctx.This, ctx.This, 1)
				return nil
			},
		},
	}
}

func BlessingOfKings() *Card {
	return &Card{
		Mana:   4,
		Name:   "王者祝福",
		Rarity: Basic,
		Text:   "使一个随从获得+4/+4。（+4攻击力/+4生命值）",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/72/Blessing_of_Kings_full.jpg",
		Tags:   []string{Spell, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllMinionCards() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				ctx.Board.EnchantCard(ctx.Target, &Enchantment{
					Id:     ctx.This.Id,
					Attack: 4,
					Health: 4,
				})
				return nil
			},
		},
	}
}

func KeeperOfUldaman() *Card {
	return &Card{
		Mana:      4,
		Name:      "奥达曼守护者",
		Rarity:    Common,
		Attack:    3,
		MaxHealth: 3,
		Text:      "<b>战吼：</b>将一个随从的攻击力和生命值变为3。",
		Image:     "https://huiji-public.huijistatic.com/hearthstone/uploads/3/3d/Art_LOE_017.png",
		Tags:      []string{Minion, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllMinionCards() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				ctx.Board.EnchantCard(ctx.Target, &Enchantment{
					Id:     ctx.This.Id,
					Attack: -ctx.Target.GetAttack() + 3,
					Health: -ctx.Target.GetMaxHealth() + 3,
				})
				return nil
			},
		},
	}
}

func StandAgainstDarkness() *Card {
	return &Card{
		Mana:   4,
		Name:   "惩黑除恶",
		Rarity: Common,
		Text:   "召唤五个1/1的白银之手新兵。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/ff/Stand_Against_Darkness_full.jpg",
		Tags:   []string{Spell, Paladin},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				for i := 0; i < 5; i++ {
					m := SilverHandRecruiter()
					ctx.Board.SummonMinion(ctx.This, m)
				}
				return nil
			},
		},
	}
}

func TruesilverChampion() *Card {
	return &Card{
		Mana:      4,
		Name:      "真银圣剑",
		Rarity:    Basic,
		Attack:    4,
		MaxHealth: 2,
		Text:      "每当你的英雄进攻，便为其恢复3点生命值。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/9/99/Truesilver_Champion_full.jpg",
		Tags:      []string{Weapon, Paladin},
		Events: map[string]Event{
			EventAfterAttack: func(ctx *EventContext) error {
				if ctx.Source.HasTag(Hero) {
					ctx.Board.Heal(ctx.This, ctx.This.Player.Hero, 3)
				}
				return nil
			},
		},
	}
}

func LayOnHands() *Card {
	return &Card{
		Mana:   8,
		Name:   "圣疗术",
		Rarity: Epic,
		Text:   "恢复8点生命值，抽3张牌。",
		Image:  "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/b/bf/Lay_on_Hands_full.jpg",
		Tags:   []string{Spell, Paladin},
		Targets: func(b *Board) []string {
			targets := []string{}
			for _, c := range b.AllCharacters() {
				targets = append(targets, c.Id)
			}
			return targets
		},
		Events: map[string]Event{
			EventSpellCast: func(ctx *EventContext) error {
				ctx.Board.Heal(ctx.This, ctx.Target, 8)
				for i := 0; i < 3; i++ {
					ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				}
				return nil
			},
		},
	}
}

func Ashbringer() *Card {
	return &Card{
		Mana:      5,
		Name:      "灰烬使者",
		Rarity:    Epic,
		Attack:    5,
		MaxHealth: 3,
		Image:     "https://wow.gamepedia.com/media/wow.gamepedia.com/a/a6/Ashbringer_TCG.jpg",
		Tags:      []string{Weapon, Paladin},
	}
}

func TirionFordring() *Card {
	return &Card{
		Mana:      8,
		Name:      "提里奥·弗丁",
		Attack:    6,
		MaxHealth: 6,
		Rarity:    Legendary,
		Text:      "<b>圣盾，嘲讽，亡语：</b>装备一把5/3的灰烬使者。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/6/63/Tirion_Fordring_full.jpg",
		Tags:      []string{Minion, Paladin, DivineShield, Taunt, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				w := Ashbringer()
				w.Player = ctx.This.Player
				ctx.Board.EquipWeapon(w)
				return nil
			},
		},
	}
}

func BloodmageThalnos() *Card {
	return &Card{
		Id:          uuid.NewString(),
		Mana:        2,
		Name:        "血法师萨尔诺斯",
		SpellDamage: 1,
		Attack:      1,
		MaxHealth:   1,
		Rarity:      Legendary,
		Text:        "<b>法术伤害+1，亡语：</b>抽一张牌。",
		Image:       "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/e/ed/Bloodmage_Thalnos_full.jpg",
		Tags:        []string{Minion, Neutral, Spellpower, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				return nil
			},
		},
	}
}

func LootHoarder() *Card {
	return &Card{
		Mana:      2,
		Name:      "战利品贮藏者",
		Attack:    2,
		MaxHealth: 1,
		Rarity:    Common,
		Text:      "<b>亡语：</b>抽一张牌。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/d/d6/Loot_Hoarder_full.jpg",
		Tags:      []string{Minion, Neutral, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				ctx.Board.DrawCard(ctx.This, ctx.This.Player, 0)
				return nil
			},
		},
	}
}

func SpawnOfNZoth() *Card {
	return &Card{
		Mana:      3,
		Name:      "恩佐斯的子嗣",
		Attack:    2,
		MaxHealth: 2,
		Rarity:    Common,
		Text:      "<b>亡语：</b>使你的所有随从获得+1/+1。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/7/73/Spawn_of_N%27Zoth_full.jpg",
		Tags:      []string{Minion, Neutral, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				for _, m := range ctx.This.Player.Minions {
					m.AddEnchantment(&Enchantment{
						Id:     ctx.This.Id,
						Name:   "Spawn of N'Zoth",
						Attack: 1,
						Health: 1,
					})
				}
				return nil
			},
		},
	}
}

func PutridSlime() *Card {
	return &Card{
		Mana:      1,
		Name:      "腐臭软泥",
		Attack:    1,
		MaxHealth: 2,
		Rarity:    Basic,
		Text:      "<b>嘲讽</b>",
		Image:     "https://gamepedia.cursecdn.com/hearthstone_gamepedia/thumb/f/fe/Slime_full.png/800px-Slime_full.png",
		Tags:      []string{Minion, Neutral, Taunt},
	}
}

func SludgeBelcher() *Card {
	return &Card{
		Mana:      5,
		Name:      "淤泥喷射者",
		Attack:    3,
		MaxHealth: 5,
		Rarity:    Common,
		Text:      "<b>嘲讽，亡语：</b>召唤一个1/2并具有嘲讽的泥浆怪。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/0/01/Sludge_Belcher_full.jpg",
		Tags:      []string{Minion, Neutral, Taunt, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				m := PutridSlime()
				ctx.Board.SummonMinion(ctx.This, m)
				return nil
			},
		},
	}
}

func EmperorThaurissan() *Card {
	return &Card{
		Mana:      6,
		Name:      "索瑞森大帝",
		Attack:    5,
		MaxHealth: 5,
		Rarity:    Legendary,
		Text:      "在你的回合结束时，你所有手牌的法力值消耗减少（1）点。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/5/5a/Emperor_Thaurissan_full.jpg",
		Tags:      []string{Minion, Neutral, EndOfTurn},
		Events: map[string]Event{
			EventEndOfTurn: func(ctx *EventContext) error {
				for _, c := range ctx.This.Player.Hand {
					ctx.Board.EnchantCard(c, &Enchantment{
						ManaCost: -1,
					})
				}
				return nil
			},
		},
	}
}

func SylvanasWindrunner() *Card {
	return &Card{
		Mana:      6,
		Name:      "希尔瓦娜斯·风行者",
		Attack:    5,
		MaxHealth: 5,
		Rarity:    Legendary,
		Text:      "<b>亡语：</b>随机获得一个敌方随从的控制权。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/f/f9/Sylvanas_Windrunner_full.jpg",
		Tags:      []string{Minion, Neutral, Deathrattle},
		Events: map[string]Event{
			EventDeathrattle: func(ctx *EventContext) error {
				myMinionCount := len(ctx.This.Player.Minions)
				if myMinionCount == 0 || myMinionCount == ctx.This.Player.MaxMinions {
					return nil
				}
				opp := ctx.Board.getOpponent(ctx.This.Player)
				minions := opp.Minions
				mcount := len(minions)
				utils.RandInt(0, mcount)
				m := minions[mcount-1]
				opp.Minions = append(opp.Minions[:mcount-1], opp.Minions[mcount:]...)
				ctx.This.Player.Minions = append(ctx.This.Player.Minions, m)
				return nil
			},
		},
	}
}

func NZothTheCorruptor() *Card {
	return &Card{
		Mana:      10,
		Name:      "恩佐斯",
		Attack:    5,
		MaxHealth: 7,
		Rarity:    Legendary,
		Text:      "<b>战吼：</b>召唤所有你在本局对战中死亡的，并具有亡语的随从。",
		Image:     "https://static.wikia.nocookie.net/hearthstone_gamepedia/images/1/13/N%27Zoth%2C_the_Corruptor_full.jpg",
		Tags:      []string{Minion, Neutral, Battlecry},
		Events: map[string]Event{
			EventBattlecry: func(ctx *EventContext) error {
				for _, m := range ctx.This.Player.Graveyard {
					if !m.HasTag(Deathrattle) {
						continue
					}
					ress := *m
					ressptr := &ress
					ressptr.Id = ""
					ressptr.Enchantments = nil
					ctx.Board.SummonMinion(ctx.This, ressptr)
				}
				return nil
			},
		},
	}
}

func RagnarosTheFirelord() *Card {
	return &Card{
		Mana:      8,
		Name:      "炎魔之王拉格纳罗斯",
		Attack:    8,
		MaxHealth: 8,
		Rarity:    Legendary,
		Text:      "<b>无法攻击。在你的回合结束时，随机对一个敌人造成8点伤害。",
		Image:     "/imgs/炎魔之王拉格纳罗斯.jpg",
		Tags:      []string{Minion, Neutral, Element, CannotAttack, EndOfTurn},
		Events: map[string]Event{
			EventEndOfTurn: func(ctx *EventContext) error {
				player := ctx.Board.getOpponent(ctx.This.Player)
				targes := append(player.Minions, player.Hero)

				targeIndex := rand.Intn(len(targes))
				ctx.Board.DealDamage(ctx.This, targes[targeIndex], 8)
				return nil
			},
		},
	}
}
