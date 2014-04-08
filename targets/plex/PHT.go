package plex

var pht = map[string]commander{
    // Navigation
    "home":             plainCommand{"/player/navigation/home"},
    "music":            plainCommand{"/player/navigation/music"},
    "moveup":           plainCommand{"/player/navigation/moveUp"},
    "movedown":         plainCommand{"/player/navigation/moveDown"},
    "moveleft":         plainCommand{"/player/navigation/moveLeft"},
    "moveright":        plainCommand{"/player/navigation/moveRight"},
    "select":           plainCommand{"/player/navigation/select"},
    "back":             plainCommand{"/player/navigation/back"},
    // Player
    "play":             plainCommand{"/player/playback/play"},
    "pause":            plainCommand{"/player/playback/pause"},
    "stop":             plainCommand{"/player/playback/stop"},
    "skipnext":         plainCommand{"/player/playback/skipNext"},
    "skipprevious":     plainCommand{"/player/playback/skipPrevious"},
    "stepforward":      plainCommand{"/player/playback/stepForward"},
    "stepback":         plainCommand{"/player/playback/stepBack"},
/*
* `/player/playback/setParameters?volume=[0, 100]&shuffle=0/1&repeat=0/1/2`
* `/player/playback/setStreams?audioStreamID=X&subtitleStreamID=Y&videoStreamID=Z`
* `/player/playback/seekTo?offset=XXX` - Offset is measured in milliseconds.
* `/player/playback/skipTo?key=X` - Playback skips to item with matching key.
* `/player/playback/playMedia` now accepts key, offset, machineIdentifier,
*/
    // Legacy
    "lnextletter":       plainCommand{"/player/navigation/nextLetter"},
    "lprevletter":       plainCommand{"/player/navigation/previousLetter"},
    "losd":              plainCommand{"/player/navigation/toggleOSD"},
    //
}
