package config

import (
	"fmt"
	"os"
	"time"

	"github.com/hectorgimenez/d2go/pkg/data/area"
	"github.com/hectorgimenez/d2go/pkg/data/item"
	"github.com/hectorgimenez/d2go/pkg/nip"

	"github.com/hectorgimenez/d2go/pkg/data/difficulty"
	"github.com/hectorgimenez/d2go/pkg/data/stat"
	"gopkg.in/yaml.v3"
)

var (
	Config *StructConfig
)

type StructConfig struct {
	Debug struct {
		Log       bool `yaml:"log"`
		RenderMap bool `yaml:"renderMap"`
	} `yaml:"debug"`
	LogFilePath   string `yaml:"logFilePath"`
	MaxGameLength int    `yaml:"maxGameLength"`
	D2LoDPath     string `yaml:"D2LoDPath"`
	Discord       struct {
		Enabled   bool   `yaml:"enabled"`
		ChannelID string `yaml:"channelId"`
		Token     string `yaml:"token"`
	} `yaml:"discord"`
	Telegram struct {
		Enabled bool   `yaml:"enabled"`
		ChatID  int64  `yaml:"chatId"`
		Token   string `yaml:"token"`
	}
	Health struct {
		HealingPotionAt     int `yaml:"healingPotionAt"`
		ManaPotionAt        int `yaml:"manaPotionAt"`
		RejuvPotionAtLife   int `yaml:"rejuvPotionAtLife"`
		RejuvPotionAtMana   int `yaml:"rejuvPotionAtMana"`
		MercHealingPotionAt int `yaml:"mercHealingPotionAt"`
		MercRejuvPotionAt   int `yaml:"mercRejuvPotionAt"`
		ChickenAt           int `yaml:"chickenAt"`
		MercChickenAt       int `yaml:"mercChickenAt"`
	} `yaml:"health"`
	Bindings struct {
		OpenInventory       string `yaml:"openInventory"`
		OpenCharacterScreen string `yaml:"openCharacterScreen"`
		OpenSkillTree       string `yaml:"openSkillTree"`
		OpenQuestLog        string `yaml:"openQuestLog"`
		Potion1             string `yaml:"potion1"`
		Potion2             string `yaml:"potion2"`
		Potion3             string `yaml:"potion3"`
		Potion4             string `yaml:"potion4"`
		ForceMove           string `yaml:"forceMove"`
		StandStill          string `yaml:"standStill"`
		SwapWeapon          string `yaml:"swapWeapon"`
		Teleport            string `yaml:"teleport"`
		TP                  string `yaml:"tp"`
		CTABattleCommand    string `yaml:"CTABattleCommand"`
		CTABattleOrders     string `yaml:"CTABattleOrders"`

		// Class Specific bindings
		Sorceress struct {
			Blizzard     string `yaml:"blizzard"`
			StaticField  string `yaml:"staticField"`
			FrozenArmor  string `yaml:"frozenArmor"`
			FireBall     string `yaml:"fireBall"`
			Nova         string `yaml:"nova"`
			EnergyShield string `yaml:"energyShield"`
		} `yaml:"sorceress"`
		Paladin struct {
			Concentration string `yaml:"concentration"`
			HolyShield    string `yaml:"holyShield"`
			Vigor         string `yaml:"vigor"`
			Redemption    string `yaml:"redemption"`
			Cleansing     string `yaml:"cleansing"`
		} `yaml:"paladin"`
	} `yaml:"bindings"`
	Inventory struct {
		InventoryLock [][]int `yaml:"inventoryLock"`
		BeltColumns   struct {
			Healing      int `yaml:"healing"`
			Mana         int `yaml:"mana"`
			Rejuvenation int `yaml:"rejuvenation"`
		} `yaml:"beltColumns"`
	} `yaml:"inventory"`
	Character struct {
		Class         string `yaml:"class"`
		Build         string `yaml:"build"`
		CastingFrames int    `yaml:"castingFrames"`
		UseMerc       bool   `yaml:"useMerc"`
		BuyPots       bool   `yaml:"buyPots"`
	} `yaml:"character"`
	Game struct {
		ClearTPArea   bool                  `yaml:"clearTPArea"`
		Difficulty    difficulty.Difficulty `yaml:"difficulty"`
		RandomizeRuns bool                  `yaml:"randomizeRuns"`
		Runs          []string              `yaml:"runs"`
		Pindleskin    struct {
			SkipOnImmunities []stat.Resist `yaml:"skipOnImmunities"`
		} `yaml:"pindleskin"`
		Mephisto struct {
			KillCouncilMembers bool `yaml:"killCouncilMembers"`
			OpenChests         bool `yaml:"openChests"`
		} `yaml:"mephisto"`
		Tristram struct {
			ClearPortal       bool `yaml:"clearPortal"`
			FocusOnElitePacks bool `yaml:"focusOnElitePacks"`
		} `yaml:"tristram"`
		Nihlathak struct {
			ClearArea bool `yaml:"clearArea"`
		} `yaml:"nihlathak"`
		Baal struct {
			OpenTP   bool `yaml:"openTP"`
			KillBaal bool `yaml:"killBaal"`
		} `yaml:"baal"`
		TerrorZone struct {
			FocusOnElitePacks bool          `yaml:"focusOnElitePacks"`
			SkipOnImmunities  []stat.Resist `yaml:"skipOnImmunities"`
			SkipOtherRuns     bool          `yaml:"skipOtherRuns"`
			Areas             []area.Area   `yaml:"areas"`
		} `yaml:"terrorZone"`
	} `yaml:"game"`
	Companion struct {
		Enabled          bool   `yaml:"enabled"`
		Leader           bool   `yaml:"leader"`
		LeaderName       string `yaml:"leaderName"`
		Remote           string `yaml:"remote"`
		GameNameTemplate string `yaml:"gameNameTemplate"`
		GamePass         string `yaml:"gamePass"`
	} `yaml:"companion"`
	Gambling struct {
		Enabled bool        `yaml:"enabled"`
		Items   []item.Name `yaml:"items"`
	} `yaml:"gambling"`
	Runtime struct {
		CastDuration time.Duration
		Rules        []nip.Rule
	}
	Db struct {
		Enabled     bool   `yaml:"enabled"`
		Host        string `yaml:"hostIp"`
		Port        string `yaml:"portNumber"`
		User        string `yaml:"userName"`
		Pass        string `yaml:"password"`
		DbName      string `yaml:"dbName"`
		BotName     string `yaml:"botName"`
		BotCharName string `yaml:"botCharName"`
	} `yaml:"db"`
}

// Load reads the config.ini file and returns a Config struct filled with data from the ini file
func Load() error {
	r, err := os.Open("config/config.yaml")
	if err != nil {
		return fmt.Errorf("error loading config.yaml: %w", err)
	}

	d := yaml.NewDecoder(r)
	if err = d.Decode(&Config); err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	rules, err := nip.ReadDir("config/pickit/")
	if err != nil {
		return err
	}

	if Config.Game.Runs[0] == "leveling" {
		levelingRules, err := nip.ReadDir("config/pickit_leveling/")
		if err != nil {
			return err
		}
		rules = append(rules, levelingRules...)
	}

	Config.Runtime.Rules = rules

	secs := float32(Config.Character.CastingFrames)*0.04 + 0.01
	Config.Runtime.CastDuration = time.Duration(secs*1000) * time.Millisecond

	return nil
}
