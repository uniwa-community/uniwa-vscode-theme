// Replicates what uniwa.nvim does with Lush/Shipwright, written from scratch since VS Code has no equivalent library.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type palette struct {
	Name, Type string

	Fg, Bg, Mid, BgContrast, TabActive, TabInactiveFg string
	Comment, LineNr, Grey                             string

	Blue, Magenta, Red, Crimson, Brown string
	Azure, Lilac, Salmon, Orange       string
	Yellow, Cabbage, Green, Teal, Cyan string

	Black, White, BrightBlack string
}

var palettes = []palette{
	{
		Name: "uniwa", Type: "dark",
		Fg: "#BED6FA", Bg: "#121C2A", Mid: "#334F72", BgContrast: "#182436", TabActive: "#212F43", TabInactiveFg: "#8CA8CE",
		Comment: "#6B6B5D", LineNr: "#426389", Grey: "#9C9C88",
		Blue: "#215E99", Magenta: "#9723A9", Red: "#B41B30", Crimson: "#EA2641", Brown: "#7B541A",
		Azure: "#68ADF4", Lilac: "#D58DEC", Salmon: "#F68895", Orange: "#E19A3C",
		Yellow: "#E1CA46", Cabbage: "#83B83B", Green: "#4DB554", Teal: "#79B3AF", Cyan: "#40BAB4",
		Black: "#080E18", White: "#B3CBED", BrightBlack: "#293342",
	},
	{
		Name: "uniwa-dark", Type: "dark",
		Fg: "#BED6FA", Bg: "#04080E", Mid: "#29415F", BgContrast: "#0B141F", TabActive: "#152030", TabInactiveFg: "#8CA8CE",
		Comment: "#6B6B5D", LineNr: "#426389", Grey: "#9C9C88",
		Blue: "#215E99", Magenta: "#9723A9", Red: "#B41B30", Crimson: "#EA2641", Brown: "#7B541A",
		Azure: "#68ADF4", Lilac: "#D58DEC", Salmon: "#F68895", Orange: "#E19A3C",
		Yellow: "#E1CA46", Cabbage: "#83B83B", Green: "#4DB554", Teal: "#79B3AF", Cyan: "#40BAB4",
		Black: "#020408", White: "#B3CBED", BrightBlack: "#1C242F",
	},
	{
		Name: "uniwa-light", Type: "light",
		Fg: "#121C2A", Bg: "#BED6FA", Mid: "#539AE3", BgContrast: "#ACCCF7", TabActive: "#C7DCFA", TabInactiveFg: "#14181E",
		Comment: "#6B6B5D", LineNr: "#3C5776", Grey: "#9C9C88",
		Blue: "#245487", Magenta: "#852694", Red: "#9F1E2E", Crimson: "#D82C41", Brown: "#6D4C1C",
		Azure: "#4795DD", Lilac: "#C86AE2", Salmon: "#F05F75", Orange: "#C0843C",
		Yellow: "#D1BC4B", Cabbage: "#729D3B", Green: "#52AE57", Teal: "#6A9996", Cyan: "#409F9A",
		Black: "#1E5380", White: "#293342", BrightBlack: "#B3CBED",
	},
	{
		Name: "uniwa-paper", Type: "light",
		Fg: "#121C2A", Bg: "#FFF5E8", Mid: "#CBA975", BgContrast: "#FEE8C9", TabActive: "#FFF5E8", TabInactiveFg: "#14181E",
		Comment: "#6B6B5D", LineNr: "#354A66", Grey: "#9C9C88",
		Blue: "#234872", Magenta: "#70257C", Red: "#871E29", Crimson: "#C22F3F", Brown: "#5D411D",
		Azure: "#427DB9", Lilac: "#B150C9", Salmon: "#DE3A54", Orange: "#9F7038",
		Yellow: "#C0AE4F", Cabbage: "#5F8437", Green: "#55A75A", Teal: "#5B807E", Cyan: "#3C8581",
		Black: "#745800", White: "#293342", BrightBlack: "#FBDFB0",
	},
}

type obj = map[string]any

func tok(fg, style string, scopes ...string) obj {
	s := obj{"foreground": fg}
	if style != "" {
		s["fontStyle"] = style
	}
	return obj{"scope": scopes, "settings": s}
}

func theme(p palette) obj {
	colors := obj{
		"editor.background":                   p.Bg,
		"editor.foreground":                   p.Fg,
		"editor.lineHighlightBackground":      p.BgContrast,
		"editor.selectionBackground":          p.Mid + "aa",
		"editor.inactiveSelectionBackground":  p.Mid + "55",
		"editor.wordHighlightBackground":      p.Mid + "55",
		"editor.findMatchBackground":          p.Brown,
		"editor.findMatchHighlightBackground": p.Orange + "55",
		"editorCursor.foreground":             p.Azure,
		"editorLineNumber.foreground":         p.LineNr,
		"editorLineNumber.activeForeground":   p.Azure,
		"editorIndentGuide.background1":       p.BgContrast,
		"editorIndentGuide.activeBackground1": p.Mid,
		"editorWhitespace.foreground":         p.Mid,
		"editorBracketMatch.background":       p.Mid + "55",
		"editorBracketMatch.border":           p.Yellow,

		"editorError.foreground":   p.Crimson,
		"editorWarning.foreground": p.Orange,
		"editorInfo.foreground":    p.Lilac,
		"editorHint.foreground":    p.Azure,

		"editorWidget.background":                      p.BgContrast,
		"editorWidget.border":                          p.Azure,
		"editorHoverWidget.background":                 p.BgContrast,
		"editorHoverWidget.border":                     p.Azure,
		"editorSuggestWidget.background":               p.BgContrast,
		"editorSuggestWidget.border":                   p.Azure,
		"editorSuggestWidget.selectedBackground":       p.Mid,
		"editorSuggestWidget.highlightForeground":      p.Azure,
		"quickInput.background":                        p.BgContrast,
		"input.background":                             p.BgContrast,
		"dropdown.background":                          p.BgContrast,
		"list.activeSelectionBackground":               p.Mid,
		"list.activeSelectionForeground":               p.Fg,
		"list.inactiveSelectionBackground":             p.BgContrast,
		"list.hoverBackground":                         p.BgContrast,
		"list.highlightForeground":                     p.Azure,
		"list.focusHighlightForeground":                p.Fg,
		"editorSuggestWidget.focusHighlightForeground": p.Fg,
		"quickInputList.focusBackground":               p.Mid,

		"foreground":                      p.Fg,
		"focusBorder":                     p.Azure,
		"button.background":               p.Azure,
		"button.foreground":               p.Bg,
		"badge.background":                p.Azure,
		"badge.foreground":                p.Bg,
		"activityBar.background":          p.Bg,
		"activityBar.foreground":          p.Azure,
		"activityBar.inactiveForeground":  p.LineNr,
		"activityBarBadge.background":     p.Azure,
		"activityBarBadge.foreground":     p.Bg,
		"sideBar.background":              p.Bg,
		"sideBar.foreground":              p.Fg,
		"sideBarSectionHeader.background": p.BgContrast,
		"titleBar.activeBackground":       p.Bg,
		"titleBar.activeForeground":       p.Fg,
		"titleBar.inactiveBackground":     p.Bg,
		"panel.background":                p.Bg,
		"panel.border":                    p.Mid,
		"panelTitle.activeBorder":         p.Azure,
		"editorGroup.border":              p.Mid,
		"sideBar.border":                  p.BgContrast,

		"editorGroupHeader.tabsBackground": p.Bg,
		"tab.activeBackground":             p.TabActive,
		"tab.activeForeground":             p.Fg,
		"tab.activeBorderTop":              p.Azure,
		"tab.inactiveBackground":           p.Bg,
		"tab.inactiveForeground":           p.TabInactiveFg,
		"tab.border":                       p.Bg,

		"statusBar.background":           p.Mid,
		"statusBar.foreground":           p.Fg,
		"statusBar.noFolderBackground":   p.Mid,
		"statusBar.debuggingBackground":  p.Salmon,
		"statusBar.debuggingForeground":  p.Bg,
		"statusBarItem.remoteBackground": p.Azure,
		"statusBarItem.remoteForeground": p.Bg,

		"gitDecoration.addedResourceForeground":     p.Cabbage,
		"gitDecoration.untrackedResourceForeground": p.Cabbage,
		"gitDecoration.modifiedResourceForeground":  p.Cyan,
		"gitDecoration.deletedResourceForeground":   p.Red,
		"gitDecoration.ignoredResourceForeground":   p.Comment,
		"editorGutter.addedBackground":              p.Cabbage,
		"editorGutter.modifiedBackground":           p.Cyan,
		"editorGutter.deletedBackground":            p.Red,
		"diffEditor.insertedTextBackground":         p.Cabbage + "33",
		"diffEditor.removedTextBackground":          p.Red + "33",

		"terminal.background":        p.Bg,
		"terminal.foreground":        p.Fg,
		"terminalCursor.foreground":  p.Fg,
		"terminal.ansiBlack":         p.Black,
		"terminal.ansiRed":           p.Salmon,
		"terminal.ansiGreen":         p.Cabbage,
		"terminal.ansiYellow":        p.Yellow,
		"terminal.ansiBlue":          p.Azure,
		"terminal.ansiMagenta":       p.Lilac,
		"terminal.ansiCyan":          p.Cyan,
		"terminal.ansiWhite":         p.White,
		"terminal.ansiBrightBlack":   p.BrightBlack,
		"terminal.ansiBrightRed":     p.Crimson,
		"terminal.ansiBrightGreen":   p.Green,
		"terminal.ansiBrightYellow":  p.Orange,
		"terminal.ansiBrightBlue":    p.Blue,
		"terminal.ansiBrightMagenta": p.Lilac,
		"terminal.ansiBrightCyan":    p.Teal,
		"terminal.ansiBrightWhite":   p.Grey,
	}

	tokens := []obj{
		tok(p.Comment, "", "comment", "punctuation.definition.comment"),
		tok(p.Cabbage, "", "string", "punctuation.definition.string"),
		tok(p.Yellow, "", "constant.character.escape", "punctuation", "meta.brace"),
		tok(p.Salmon, "", "constant.character"),
		tok(p.Lilac, "", "constant.numeric"),
		tok(p.Orange, "", "constant.language", "constant.other", "variable.other.constant", "variable.other.enummember"),
		tok(p.Azure, "", "variable", "variable.parameter", "variable.other", "support.variable", "entity.other.attribute-name", "meta.object-literal.key"),
		tok(p.Blue, "", "entity.name.function", "support.function", "meta.function-call"),
		tok(p.Blue, "bold", "keyword", "storage.modifier", "keyword.control"),
		tok(p.Fg, "", "keyword.operator"),
		tok(p.Magenta, "bold", "meta.preprocessor", "keyword.control.directive", "keyword.control.import", "keyword.import", "keyword.package"),
		tok(p.Teal, "italic", "entity.name.type", "entity.name.class", "support.type", "support.class", "storage.type", "entity.name.tag"),
		tok(p.Crimson, "", "invalid"),

		tok(p.Fg, "bold", "markup.heading"),
		tok(p.Blue, "bold", "markup.heading.1", "heading.1"),
		tok(p.Azure, "bold", "markup.heading.2", "heading.2"),
		tok(p.Cyan, "bold", "markup.heading.3", "heading.3"),
		tok(p.Lilac, "bold", "markup.heading.4", "heading.4"),
		tok(p.Blue, "bold", "markup.bold"),
		tok(p.Azure, "italic", "markup.italic", "markup.quote"),
		tok(p.Red, "strikethrough", "markup.strikethrough"),
		tok(p.Orange, "underline", "markup.underline"),
		tok(p.Cyan, "", "markup.underline.link", "string.other.link"),
		tok(p.Teal, "", "markup.inline.raw", "markup.fenced_code", "markup.math"),
		tok(p.Cabbage, "", "markup.inserted"),
		tok(p.Red, "", "markup.deleted"),
		tok(p.Cyan, "", "markup.changed"),
	}

	return obj{
		"$schema":              "vscode://schemas/color-theme",
		"name":                 p.Name,
		"type":                 p.Type,
		"semanticHighlighting": true,
		"semanticTokenColors": obj{
			"variable":          p.Azure,
			"parameter":         p.Azure,
			"property":          p.Azure,
			"function":          p.Blue,
			"method":            p.Blue,
			"enumMember":        p.Orange,
			"variable.readonly": p.Orange,
			"namespace":         obj{"foreground": p.Magenta, "bold": true},
			"type":              obj{"foreground": p.Teal, "italic": true},
			"class":             obj{"foreground": p.Teal, "italic": true},
			"interface":         obj{"foreground": p.Teal, "italic": true},
			"struct":            obj{"foreground": p.Teal, "italic": true},
		},
		"colors":      colors,
		"tokenColors": tokens,
	}
}

func main() {
	for _, p := range palettes {
		b, err := json.MarshalIndent(theme(p), "", "  ")
		if err != nil {
			panic(err)
		}
		path := filepath.Join("themes", p.Name+"-color-theme.json")
		if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
			panic(err)
		}
		fmt.Println("wrote", path)
	}
}
