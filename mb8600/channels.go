package mb8600

import (
	"fmt"
	"strconv"
	"strings"
)

type DownstreamChannel struct {
	Channel           int
	ChannelID         int
	LockStatus        string
	Modulation        string
	Frequency         float64
	Power             float64
	SignalToNoise     float64
	CorrectedErrors   uint64
	UncorrectedErrors uint64
}

type UpstreamChannel struct {
	Channel     int
	ChannelID   int
	LockStatus  string
	ChannelType string
	SymbolRate  float64
	Frequency   float64
	Power       float64
}

func NewDownstreamChannelsFromResponse(response string) ([]*DownstreamChannel, error) {
	var channels []*DownstreamChannel

	for _, line := range strings.Split(response, "|+|") {
		if channel, err := NewDownstreamChannelFromLine(line); err != nil {
			return nil, err
		} else {
			channels = append(channels, channel)
		}
	}

	return channels, nil
}

func NewDownstreamChannelFromLine(line string) (*DownstreamChannel, error) {
	parts := strings.Split(line, "^")
	if len(parts) != 9 {
		return nil, fmt.Errorf("invalid number of parts in downstream channel line: %d", len(parts))
	}

	channel, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	lockStatus := parts[1]
	modulation := parts[2]
	channelId, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	frequency, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, err
	}

	power, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return nil, err
	}

	snr, err := strconv.ParseFloat(parts[6], 64)
	if err != nil {
		return nil, err
	}

	corrected, err := strconv.ParseUint(parts[7], 10, 64)
	if err != nil {
		return nil, err
	}

	uncorrected, err := strconv.ParseUint(parts[8], 10, 64)
	if err != nil {
		return nil, err
	}

	return &DownstreamChannel{
		Channel:           channel,
		ChannelID:         channelId,
		LockStatus:        lockStatus,
		Modulation:        modulation,
		Frequency:         frequency,
		Power:             power,
		SignalToNoise:     snr,
		CorrectedErrors:   corrected,
		UncorrectedErrors: uncorrected,
	}, nil
}

func NewUpstreamChannelsFromResponse(response string) ([]*UpstreamChannel, error) {
	var channels []*UpstreamChannel

	for _, line := range strings.Split(response, "|+|") {
		if channel, err := NewUpstreamChannelFromLine(line); err != nil {
			return nil, err
		} else {
			channels = append(channels, channel)
		}
	}

	return channels, nil
}

func NewUpstreamChannelFromLine(line string) (*UpstreamChannel, error) {
	parts := strings.Split(line, "^")
	if len(parts) != 7 {
		return nil, fmt.Errorf("invalid number of parts in upstream channel line: %d", len(parts))
	}

	channel, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	lockStatus := parts[1]
	channelType := parts[2]
	channelId, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	symbolRate, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, err
	}

	frequency, err := strconv.ParseFloat(parts[5], 64)
	if err != nil {
		return nil, err
	}

	power, err := strconv.ParseFloat(parts[6], 64)
	if err != nil {
		return nil, err
	}

	return &UpstreamChannel{
		Channel:     channel,
		ChannelID:   channelId,
		LockStatus:  lockStatus,
		ChannelType: channelType,
		SymbolRate:  symbolRate,
		Frequency:   frequency,
		Power:       power,
	}, nil
}
