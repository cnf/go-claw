package plex

var pht = map[string]commander{
    // Navigation
    "Home":             plainCommand{"/player/navigation/home"},
    "Music":            plainCommand{"/player/navigation/music"},
    "MoveUp":           plainCommand{"/player/navigation/moveUp"},
    "MoveDown":         plainCommand{"/player/navigation/moveDown"},
    "MoveLeft":         plainCommand{"/player/navigation/moveLeft"},
    "MoveRight":        plainCommand{"/player/navigation/moveRight"},
    "Select":           plainCommand{"/player/navigation/select"},
    "Back":             plainCommand{"/player/navigation/back"},
    // Player
    "Play":             plainCommand{"/player/playback/play"},
    "Pause":            plainCommand{"/player/playback/pause"},
    "Stop":             plainCommand{"/player/playback/stop"},
    "SkipNext":         plainCommand{"/player/playback/skipNext"},
    "SkipPrevious":     plainCommand{"/player/playback/skipPrevious"},
    "StepForward":      plainCommand{"/player/playback/stepForward"},
    "StepBack":         plainCommand{"/player/playback/stepBack"},
/*
* `/player/playback/setParameters?volume=[0, 100]&shuffle=0/1&repeat=0/1/2`
* `/player/playback/setStreams?audioStreamID=X&subtitleStreamID=Y&videoStreamID=Z`
* `/player/playback/seekTo?offset=XXX` - Offset is measured in milliseconds.
* `/player/playback/skipTo?key=X` - Playback skips to item with matching key.
* `/player/playback/playMedia` now accepts key, offset, machineIdentifier,
*/
    // Legacy
    "LNextLetter":       plainCommand{"/player/navigation/nextLetter"},
    "LPrevLetter":       plainCommand{"/player/navigation/previousLetter"},
    "LOSD":              plainCommand{"/player/navigation/toggleOSD"},
    //
}
