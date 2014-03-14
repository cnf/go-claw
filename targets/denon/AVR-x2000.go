package denon

var AVRX2000 = map[string]Commander{
    "PowerOn":    PlainCommand{"PWON"},
    "PowerOff":   PlainCommand{"PWSTANDBY"},
    "VolumeUp":   PlainCommand{"MVUP"},
    "VolumeDown": PlainCommand{"MVDOWN"},
    "Volume":     VolumeCommand{"MV%02d", 0, 98},
    "MuteOn":     PlainCommand{"MUON"},
    "MuteOff":    PlainCommand{"MUOFF"},
    //"MuteToggle": MU?
    "Example":    RangeCommand{"MV%02d", 0, 98},
    // Select Input
    "Input1":     PlainCommand{"SISAT/CBL"},
    "Input2":     PlainCommand{"SIDVD"},
    "Input3":     PlainCommand{"SIBD"},
    "Input4":     PlainCommand{"SIGAME"},
    "Input5":     PlainCommand{"SISMPLAY"},
    "Input6":     PlainCommand{"SICD"},
    "Input7":     PlainCommand{"SIAUX1"},
    // Zone 2
    "Z2PowerOn":  PlainCommand{"Z2ON"},
    "Z2PowerOff": PlainCommand{"Z2OFF"},
    // SD Mode?
    // Select video
    // MS?

    // NAVIGATION
    "MoveUp":       PlainCommand{"MNCUP"},
    "MoveDown":     PlainCommand{"MNCDOWN"},
    "MoveLeft":     PlainCommand{"MNCLT"},
    "MoveRight":    PlainCommand{"MNCRT"},
    "Select":       PlainCommand{"MNENT"},
    "Back":         PlainCommand{"MNRTN"},
    "Option":       PlainCommand{"MNOPT"},
    "Info":         PlainCommand{"MNINF"},
    "SetupOn":      PlainCommand{"MNMEN ON"},
    "SetupOff":     PlainCommand{"MNMEN OFF"},
}
