package denon

var AVRX2000 = map[string]Commander{
    "poweron":    PlainCommand{"ZMON"},
    "poweroff":   PlainCommand{"PWSTANDBY"},
    "volumeup":   PlainCommand{"MVUP"},
    "volumedown": PlainCommand{"MVDOWN"},
    "volume":     VolumeCommand{"MV%02d", 0, 98},
    "muteon":     PlainCommand{"MUON"},
    "muteoff":    PlainCommand{"MUOFF"},
    //"MuteToggle": MU?
    "example":    RangeCommand{"MV%02d", 0, 98},
    // Select Input
    "input1":     PlainCommand{"SISAT/CBL"},
    "input2":     PlainCommand{"SIDVD"},
    "input3":     PlainCommand{"SIBD"},
    "input4":     PlainCommand{"SIGAME"},
    "input5":     PlainCommand{"SISMPLAY"},
    "input6":     PlainCommand{"SICD"},
    "input7":     PlainCommand{"SIAUX1"},
    // Zone 2
    "z2poweron":  PlainCommand{"Z2ON"},
    "z2poweroff": PlainCommand{"Z2OFF"},
    // SD Mode?
    // Select video
    // MS?

    // NAVIGATION
    "moveUp":       PlainCommand{"MNCUP"},
    "moveDown":     PlainCommand{"MNCDOWN"},
    "moveLeft":     PlainCommand{"MNCLT"},
    "moveRight":    PlainCommand{"MNCRT"},
    "select":       PlainCommand{"MNENT"},
    "back":         PlainCommand{"MNRTN"},
    "option":       PlainCommand{"MNOPT"},
    "info":         PlainCommand{"MNINF"},
    "setupon":      PlainCommand{"MNMEN ON"},
    "setupoff":     PlainCommand{"MNMEN OFF"},
}
