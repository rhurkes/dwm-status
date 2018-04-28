/* A simple program to update the dwm status bar
Requires:
	 Material Design Icons font: https://material.io/icons/
	 acpi to display battery details
	 amixer to display volume details */

package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	useUTC                 = true
	wirelessInterface      = "wlp1s0"
	interval               = 1000 * time.Millisecond
	separator              = "      "
	label                  = "%s %s"
	powerSupply            = "/sys/class/power_supply/"
	batteryDischargingIcon = "\uE1A5"
	batteryChargingIcon    = "\uE1A3"
	batteryUnknownIcon     = "\uE1A6"
	wifiConnectedIcon      = "\uE1BA"
	mutedIcon              = "\uE04F"
	lowVolumeIcon          = "\uE04E"
	mediumVolumeIcon       = "\uE04D"
	highVolumeIcon         = "\uE050"
)

var (
	ssidRegex             = regexp.MustCompile(`ESSID:"(.+)"`)
	amixerRegex           = regexp.MustCompile(`Front Left: .+\[(.+)\] \[(.+)\]`)
	acpiRegex             = regexp.MustCompile(`Battery 0: (.+), (\d{1,3}%), (\d{2}:\d{2})?`)
	batteryTime           = ""
	lastChargeState       = ""
	nextBatteryTimeUpdate = time.Now().Add(-1 * time.Minute)
)

func getTime() string {
	now := time.Now()
	if useUTC {
		return now.UTC().Format("15:04Z")
	}

	return now.Local().Format("03:04")
}

func getNetwork() string {
	network := ""

	stdout, err := exec.Command("iwconfig", wirelessInterface).Output()
	if err != nil {
		return network
	}

	ssidMatch := ssidRegex.FindStringSubmatch(string(stdout[:]))
	if len(ssidMatch) == 2 {
		network = fmt.Sprintf(label, wifiConnectedIcon, ssidMatch[1])
	}

	return network
}

func getPower() string {
	acpiResult, err := exec.Command("acpi").Output()
	if err != nil {
		return batteryUnknownIcon
	}

	acpi := strings.TrimSpace(string(acpiResult[:]))
	if acpi == "Battery 0: Full, 100%" {
		return fmt.Sprintf(label, batteryDischargingIcon, "100%")
	}

	match := acpiRegex.FindStringSubmatch(string(acpiResult[:]))
	if len(match) != 3 && len(match) != 4 {
		return batteryUnknownIcon
	}

	icon := batteryDischargingIcon
	if match[1] == "Charging" {
		icon = batteryChargingIcon
	}

	if lastChargeState != match[1] || time.Now().After(nextBatteryTimeUpdate) {
		nextBatteryTimeUpdate = time.Now().Add(time.Minute)
		lastChargeState = match[1]
		if len(match) == 4 {
			batteryTime = match[3]
		} else {
			batteryTime = ""
		}
	}

	return fmt.Sprintf(label+" %s", icon, match[2], batteryTime)
}

func getVolume() string {
	volume := "-"
	stdout, err := exec.Command("amixer", "get", "Master").Output()
	if err != nil {
		return volume
	}

	match := amixerRegex.FindStringSubmatch(string(stdout[:]))
	if len(match) != 3 {
		return volume
	}

	icon := lowVolumeIcon
	volumeLevel, err := strconv.Atoi(strings.Replace(match[1], "%", "", 1))

	if strings.ToLower(match[2]) == "off" {
		icon = mutedIcon
	} else if volumeLevel > 66 {
		icon = highVolumeIcon
	} else if volumeLevel > 33 {
		icon = mediumVolumeIcon
	}

	volume = fmt.Sprintf(label, icon, strconv.Itoa(volumeLevel))

	return volume
}

func aggregateValues() string {
	var values []string
	values = append(values, getNetwork())
	values = append(values, getVolume())
	values = append(values, getPower())
	values = append(values, getTime())

	return strings.Join(values, separator)
}

func updateStatusBar(text string) {
	paddedText := fmt.Sprintf(" %s ", text)
	exec.Command("xsetroot", "-name", paddedText).Run()
}

func main() {
	for {
		updateStatusBar(aggregateValues())
		time.Sleep(interval)
	}
}
