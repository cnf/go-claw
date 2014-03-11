package denon

var AVRX2000 = map[string]Commander{
    "PowerOn":    PlainCommand{"PWON"},
    "VolumeUp":   PlainCommand{"MVUP"},
    "VolumeDown": PlainCommand{"MVDOWN"},
    "Volume":     VolumeCommand{"MV%02d", 0, 98},
    "Example":    RangeCommand{"MV%02d", 0, 98},
}
