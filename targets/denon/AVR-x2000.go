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
    // SD Mode?
    // Select video
    // MS?

    // NAVIGATION
    "CtrlUp":       PlainCommand{"MNCUP"},
    "CtrlDown":     PlainCommand{"MNCDOWN"},
    "CtrlLeft":     PlainCommand{"MNCLT"},
    "CtrlRight":    PlainCommand{"MNCRT"},
    "CtrlEnter":    PlainCommand{"MNENT"},
    "CtrlReturn":   PlainCommand{"MNRTN"},
    "CtrlOption":   PlainCommand{"MNOPT"},
    "CtrlInfo":     PlainCommand{"MNINF"},
    "CtrlSetupOn":  PlainCommand{"MNMEN ON"},
    "CtrlSetupOff": PlainCommand{"MNMEN OFF"},
}
