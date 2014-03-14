package plex

var PHT = map[string]Commander{
    // Player
    "MoveUp":           PlainCommand{"/player/navigation/moveUp"},
    "MoveDown":         PlainCommand{"/player/navigation/moveDown"},
    "MoveLeft":         PlainCommand{"/player/navigation/moveLeft"},
    "MoveRight":        PlainCommand{"/player/navigation/moveRight"},
    "Select":           PlainCommand{"/player/navigation/select"},
    "Home":             PlainCommand{"/player/navigation/home"},
    "Back":             PlainCommand{"/player/navigation/back"},
    "Play":             PlainCommand{"/player/playback/play"},
    "Pause":            PlainCommand{"/player/playback/pause"},
    "Stop":             PlainCommand{"/player/playback/stop"},
    "OSD":              PlainCommand{"/player/navigation/toggleOSD"},
    //
    "StepForward":      PlainCommand{"/player/playback/stepForward"},
    "StepBack":         PlainCommand{"/player/playback/stepBack"},
    // Legacy
    "NextLetter":       PlainCommand{"/player/navigation/nextLetter"},
    "PrevLetter":       PlainCommand{"/player/navigation/previousLetter"},
    //
}
