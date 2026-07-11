package main

import (
    "encoding/binary"
    "encoding/json"
    "io"
    "log"
    "os"
    "path/filepath"

    "golang.org/x/sys/unix"
)

// Define the IOCTL constant directly to bypass build tag issues
const EVIOCGRAB = 0x40044590

type Config struct {
    Controller string `json:"controller"`
    Evtest     string `json:"evtest"`
}

type InputEvent struct {
    Sec   int64
    Usec  int64
    Type  uint16
    Code  uint16
    Value int32
}

func main() {
    exe, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }

    configPath := filepath.Join(filepath.Dir(exe), "..", "defaultConfig.json")

    file, err := os.Open(configPath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    var cfg Config

    if err := json.NewDecoder(file).Decode(&cfg); err != nil {
        log.Fatal(err)
    }

    dev, err := os.Open(cfg.Evtest)
    if err != nil {
        log.Fatal(err)
    }
    defer dev.Close()

    fd := int(dev.Fd())

    // Changed unix.EVIOCGRAB to our local constant EVIOCGRAB
    if err := unix.IoctlSetInt(fd, EVIOCGRAB, 1); err != nil {
        log.Fatal(err)
    }
    defer unix.IoctlSetInt(fd, EVIOCGRAB, 0)

    for {
        var ev InputEvent

        if err := binary.Read(dev, binary.LittleEndian, &ev); err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }

        if err := binary.Write(os.Stdout, binary.LittleEndian, &ev); err != nil {
            log.Fatal(err)
        }
    }
}