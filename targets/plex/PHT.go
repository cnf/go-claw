package plex

var pht = map[string]commander{
    // Player
    "MoveUp":           plainCommand{"/player/navigation/moveUp"},
    "MoveDown":         plainCommand{"/player/navigation/moveDown"},
    "MoveLeft":         plainCommand{"/player/navigation/moveLeft"},
    "MoveRight":        plainCommand{"/player/navigation/moveRight"},
    "Select":           plainCommand{"/player/navigation/select"},
    "Home":             plainCommand{"/player/navigation/home"},
    "Back":             plainCommand{"/player/navigation/back"},
    "Play":             plainCommand{"/player/playback/play"},
    "Pause":            plainCommand{"/player/playback/pause"},
    "Stop":             plainCommand{"/player/playback/stop"},
    "OSD":              plainCommand{"/player/navigation/toggleOSD"},
    //
    "StepForward":      plainCommand{"/player/playback/stepForward"},
    "StepBack":         plainCommand{"/player/playback/stepBack"},
    // Legacy
    "NextLetter":       plainCommand{"/player/navigation/nextLetter"},
    "PrevLetter":       plainCommand{"/player/navigation/previousLetter"},
    //
}
