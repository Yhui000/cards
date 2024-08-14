package cards

import (
	"strconv"
	"strings"
)

const (
	// Classes
	Neutral = "中立"
	Druid   = "德鲁伊"
	Hunter  = "猎人"
	Mage    = "法师"
	Paladin = "圣骑士"
	Priest  = "牧师"
	Rogue   = "潜行者"
	Shaman  = "萨满祭司"
	Warlock = "术士"
	Warrior = "战士"

	// Rarities
	Basic     = "基础"
	Common    = "普通"
	Rare      = "稀有"
	Epic      = "史诗"
	Legendary = "传说"

	// Card types
	Minion    = "随从"
	Spell     = "法术"
	Hero      = "英雄"
	Heropower = "英雄技能"
	Weapon    = "武器"
	Secret    = "奥秘"

	// Tribes
	Pirate  = "海盗"
	Murloc  = "鱼人"
	Beast   = "野兽"
	Element = "元素"

	// Keywords
	Charge       = "冲锋"
	Rush         = "突袭"
	Battlecry    = "战吼"
	Windfury     = "风怒"
	MegaWindfury = "超级风怒"
	Deathrattle  = "亡语"
	Taunt        = "嘲讽"
	BloodPayment = "鲜血支付"
	Reborn       = "复生"
	DivineShield = "圣盾"
	Spellpower   = "法术强度"
	EndOfTurn    = "回合结束"
	CannotAttack = "无法攻击"
)

type Card struct {
	Id           string
	Player       *Player
	SpellDamage  int
	Mana         int
	Name         string
	Attack       int
	Health       int
	MaxHealth    int
	Rarity       string
	Text         string
	Image        string
	Tags         []string
	Tribe        string
	Events       map[string]Event
	Sleeping     bool
	AttacksLeft  int
	Targets      func(*Board) []string
	Enchantments []*Enchantment
}

func (c *Card) AddEnchantment(e *Enchantment) {
	c.Health += e.Health
	c.Enchantments = append(c.Enchantments, e)
}

// Returns the formatted text with spell power, info for example
func (c *Card) GetFormattedText() string {
	text := c.Text
	text = strings.ReplaceAll(text, "%SpellDamage%", strconv.Itoa(c.SpellDamage))
	return text
}

// Returns the attack value with all enchantments if there are any
func (c *Card) GetAttack() int {
	attack := c.Attack
	for _, e := range c.Enchantments {
		attack += e.Attack
	}
	return attack
}

// Returns the max health value with all enchantments if there are any
func (c *Card) GetMaxHealth() int {
	maxHealth := c.MaxHealth
	for _, e := range c.Enchantments {
		maxHealth += e.Health
	}
	return maxHealth
}

func (c *Card) GetTag(tag string) int {
	for i, t := range c.GetTags() {
		if t == tag {
			return i
		}
	}
	return -1
}

func (c *Card) GetTags() []string {
	tags := c.Tags
	for _, e := range c.Enchantments {
		tags = append(tags, e.Tags...)
	}
	return tags
}

func (c *Card) HasTag(tag string) bool {
	for _, t := range c.GetTags() {
		if t == tag {
			return true
		}
	}
	return false
}

func (c *Card) DelTag(tag string) {
	for i, t := range c.GetTags() {
		if t == tag {
			c.Tags = append(c.Tags[:i], c.Tags[i+1:]...)
		}
	}
}

func (c *Card) GetManaCost() int {
	cost := c.Mana
	for _, e := range c.Enchantments {
		cost += e.ManaCost
	}
	if cost < 0 {
		return 0
	}
	return cost
}

func (c *Card) GetEnch(id string) *Enchantment {
	for _, e := range c.Enchantments {
		if e.Id == id {
			return e
		}
	}
	return nil
}

func (c *Card) DelEnchId(id string) {
	for i, e := range c.Enchantments {
		if e.Id == id {
			c.Enchantments = append(c.Enchantments[:i], c.Enchantments[i+1:]...)
		}
	}
}
