package cards

const (
	EventEndOfAction          = "动作结束后触发"
	EventBeforeDrawCard       = "抽卡前触发"
	EventAfterDrawCard        = "抽卡后触发"
	EventBeforeAddToHand      = "添加到手中前触发"
	EventAfterAddToHand       = "添加到手中后触发"
	EventBeforeShuffleDeck    = "洗牌前触发"
	EventAfterShuffleDeck     = "洗牌后触发"
	EventStartOfGame          = "游戏开始时触发"
	EventStartOfTurn          = "回合开始时触发"
	EventEndOfTurn            = "回合结束时触发"
	EventAfterCardPlay        = "卡牌打出后触发"
	EventBeforeCardPlay       = "卡牌打出前触发"
	EventBeforeHeroPower      = "使用英雄技能前触发"
	EventAfterHeroPower       = "使用英雄技能后触发"
	EventBattlecry            = "战吼效果触发"
	EventSpellCast            = "法术施放时触发"
	EventBeforeSummon         = "召唤随从前触发"
	EventAfterSummon          = "召唤随从后触发"
	EventBeforeDamage         = "受到伤害前触发"
	EventAfterDamage          = "受到伤害后触发"
	EventBeforeHeal           = "恢复生命前触发"
	EventAfterHeal            = "恢复生命后触发"
	EventBeforeDestroyMinion  = "摧毁随从前触发"
	EventDestroyMinion        = "摧毁随从时触发"
	EventDeathrattle          = "亡语效果触发"
	EventAfterDestroyMinion   = "摧毁随从后触发"
	EventBeforeAttack         = "攻击前触发"
	EventAfterAttack          = "攻击后触发"
	EventBeforeLoseDurability = "失去耐久前触发"
	EventAfterLoseDurability  = "失去耐久后触发"
	EventBeforeWeaponDestroy  = "武器摧毁前触发"
	EventAfterWeaponDestroy   = "武器摧毁后触发"
	EventBeforeWeaponEquip    = "武器装备前触发"
	EventAfterWeaponEquip     = "武器装备后触发"
	EventHeroPower            = "英雄能力"
	EventSummon               = "召唤"

	// Historic
	Draw    = "draw"
	Attack  = "attack"
	Heal    = "heal"
	Play    = "play"
	EndTurn = "endturn"
	Summon  = "summon"
)

type Event func(ctx *EventContext) error

type EventContext struct {
	Board        *Board
	This         *Card
	Source       *Card
	Target       *Card
	HealAmount   int
	DamageAmount int
}
