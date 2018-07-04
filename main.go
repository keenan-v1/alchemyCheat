package main

import (
	"alchemyCheat/types"
)

func main() {
	// Requires a data.json
	// data.json format (acquired by source digging and some regex replace love):
	// [{ "id": 1, "base": true, "prime": true, "hidden": true, "condition": {"type":"progress","total":50}, "name": "element", "parents": [[2,3] ...], "children": [5,6 ...] }, ... ]
	types.ProcessElements()
	// Must be run before history
	types.PopulateCombinations()

	// Comment out if you don't have a history.json file
	// history.json format (taken from viewing frames in the firebase websocket and some regex replace love):
	// [{ "key": "12|345", "when": 123456789, "v": 123456789 }, ... ]
	types.ProcessHistory()

	// Uncomment the output you want
	//types.PrintCombinations()
	types.PrintElements()
	//types.DebugCombinations()
}
