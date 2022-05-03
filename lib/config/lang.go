package config

import "strconv"

type Strings struct {
	Lang     string
	Help     string
	Ranking  string
	CurrConf string
	Sync     syncstr
	Wiki     wikistr
	Usage    usagestr
	Error    errorstr
}

type errorstr struct {
	Title         string
	UnknownTitle  string
	UnknownDesc   string
	NoCmd         string
	SubCmd        string
	Syntax        string
	SyntaxDesc    string
	Onefile       string
	EmojiNotFound string
}

type usagestr struct {
	Title    string
	Config   configusagestr
	Ranking  rankingusagestr
	WikiDesc string
}

type syncstr struct {
	Title            string
	NewEmoji         string
	IDChangedEmoji   string
	NameChangedEmoji string
	DeletedEmoji     string
	OverEmoji        string
	ToCleanEmoji     string
	CleanLogTitle    string
	CleanLogDesc     string
	WeightTitle      string
	WeightDesc       string
}

type wikistr struct {
	Title       string
	IDMatched   string
	NameMatched string
	Ignored     string
	Updated     string
}

type configusagestr struct {
	Desc   string
	Prefix string
	Lang   string
	Bots   string
	Weight string
}

type rankingusagestr struct {
	Desc   string
	Page   string
	Invert string
	Num    string
	Period string
}

var (
	Lang map[string]Strings
)

func loadLang() {
	Lang = map[string]Strings{}
	Lang["japanese"] = Strings{
		Lang:     "japanese",
		Help:     "Botの使い方に関しては、Wikiをご覧ください。",
		Ranking:  "ランキング",
		CurrConf: "現在の設定",
		Sync: syncstr{
			Title:            "更新結果",
			NewEmoji:         "新しく検出された絵文字",
			IDChangedEmoji:   "登録され直した絵文字",
			NameChangedEmoji: "名前が変更された絵文字",
			DeletedEmoji:     "削除された絵文字",
			OverEmoji:        "...残り: ",
			ToCleanEmoji:     "これらの絵文字をDBから削除するには以下のコマンドを実行してください: ",
			CleanLogTitle:    "無効なログの削除",
			CleanLogDesc:     "件の無効なログを削除しました",
			WeightTitle:      "重みの変更",
			WeightDesc:       "件のログを更新しました",
		},
		Wiki: wikistr{
			Title:       "インポート成功",
			IDMatched:   "IDで一致した項目数",
			NameMatched: "名前で一致した項目数",
			Ignored:     "一致する絵文字が発見できなかった項目数",
			Updated:     "更新された件数",
		},
		Usage: usagestr{
			Title: "使い方: ",
			Config: configusagestr{
				Desc:   "各種設定を行います。\n設定項目と内容は以下を確認ください。",
				Prefix: "コマンドの接頭詞を指定します。\nデフォルトは" + CurrentConfig.Guild.Prefix + "です。",
				Lang:   "言語を指定します。\nデフォルトは" + CurrentConfig.Guild.Lang + "です。",
				Bots:   "Botを記録対象にするか指定します。\nデフォルトは" + strconv.FormatBool(CurrentConfig.Guild.Recordbots) + "です。",
				Weight: "それぞれの動作ごとのポイントを指定します。\nこのポイントはランキングに用いられます。\n変更したら必ずsyncを行ってください。\nデフォルトは" + strconv.Itoa(CurrentConfig.Guild.Weight.Message) + " " + strconv.Itoa(CurrentConfig.Guild.Weight.Reactnew) + " " + strconv.Itoa(CurrentConfig.Guild.Weight.Reactadd) + "です。",
			},
			Ranking: rankingusagestr{
				Desc:   "絵文字のランキングを表示します。",
				Page:   "何ページ目を見るか指定します。\n省略した場合1になります。",
				Invert: "指定すると逆順になり、使われていない絵文字が確認できます。",
				Num:    "1ページに何件入れるか指定します。\n省略した場合10になります。",
				Period: "集計対象の期間を秒数で指定します。\n省略した場合2592000(30日)になります。",
			},
			WikiDesc: "絵文字情報をポイント順にマークダウンで出力します",
		},
		Error: errorstr{
			UnknownTitle:  "不明なエラーが発生しました。",
			UnknownDesc:   "この問題は管理者に報告されます。",
			NoCmd:         "コマンドが見つかりません。",
			SubCmd:        "サブコマンドが不正です。",
			Syntax:        "構文エラー",
			SyntaxDesc:    "パラメータの解析に失敗しました。\nコマンドの構文が正しいか確認してください。",
			Onefile:       "1つのファイルを添付する必要があります",
			EmojiNotFound: "該当する絵文字が見つかりません",
		},
	}
	Lang["english"] = Strings{
		Lang:    "english",
		Help:    "Usage is available on the Wiki.",
		Ranking: "Ranking",
		Usage: usagestr{ // h2yk's task
			Title: "使い方: ",
			Ranking: rankingusagestr{
				Page:   "何ページ目を見るか指定します。\n省略した場合1になります。",
				Invert: "trueにするとpointの昇順になり、使われていない絵文字が確認できます。\n省略した場合falseになります。",
				Num:    "1ページに何件入れるか指定します。\n省略した場合10になります。",
				Period: "集計対象の期間を秒数で指定します。\n省略した場合2592000(30日)になります。",
			},
		},
		Error: errorstr{
			UnknownTitle: "Unknown Error!",
			UnknownDesc:  "This will be reported.",
			NoCmd:        "Command not found.",
			SubCmd:       "Invalid subcommand.",
			SyntaxDesc:   "パラメータの解析に失敗しました。\nコマンドの構文が正しいか確認してください。",
		},
	}
}
