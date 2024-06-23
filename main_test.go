package main

import (
	"testing"
	"time"
)

func TestGetLastMonday(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "Monday",
			date:     time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
			expected: "2022-01-03",
		},
		{
			name:     "Tuesday",
			date:     time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC),
			expected: "2022-01-03",
		},
		{
			name:     "Sunday",
			date:     time.Date(2022, time.January, 9, 0, 0, 0, 0, time.UTC),
			expected: "2022-01-03",
		},
		{
			name:     "Leap year February 29",
			date:     time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: "2024-02-26",
		},
		{
			name:     "Week spanning two months",
			date:     time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			expected: "2022-01-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLastMonday(tt.date)
			if result != tt.expected {
				t.Errorf("Expected %s, but got %s", tt.expected, result)
			}
		})
	}
}
func TestGetHours(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected int
	}{
		{
			name:     "Zero duration",
			duration: 0,
			expected: 0,
		},
		{
			name:     "Less than an hour",
			duration: time.Minute * 30,
			expected: 0,
		},
		{
			name:     "Exactly one hour",
			duration: time.Hour,
			expected: 1,
		},
		{
			name:     "More than one hour",
			duration: time.Hour + time.Minute*30,
			expected: 1,
		},
		{
			name:     "Multiple hours",
			duration: time.Hour*3 + time.Minute*45,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getHours(tt.duration)
			if result != tt.expected {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
			}
		})
	}
}
func TestGetMinutes(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected int
	}{
		{
			name:     "Zero duration",
			duration: 0,
			expected: 0,
		},
		{
			name:     "Less than a minute rounds up",
			duration: time.Second * 30,
			expected: 1,
		},
		{
			name:     "Less than a minute rounds down",
			duration: time.Second * 20,
			expected: 0,
		},
		{
			name:     "Exactly one minute",
			duration: time.Minute,
			expected: 1,
		},
		{
			name:     "More than one minute",
			duration: time.Minute + time.Second*30,
			expected: 2,
		},
		{
			name:     "Multiple minutes",
			duration: time.Minute*3 + time.Second*45,
			expected: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMinutes(tt.duration)
			if result != tt.expected {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		t1       togglTime
		t2       togglTime
		expected togglTime
	}{
		{
			name:     "Adding zero time",
			t1:       togglTime{hours: 1, minutes: 30},
			t2:       togglTime{hours: 0, minutes: 0},
			expected: togglTime{hours: 1, minutes: 30},
		},
		{
			name:     "Adding minutes",
			t1:       togglTime{hours: 1, minutes: 30},
			t2:       togglTime{hours: 0, minutes: 45},
			expected: togglTime{hours: 2, minutes: 15},
		},
		{
			name:     "Adding hours",
			t1:       togglTime{hours: 1, minutes: 30},
			t2:       togglTime{hours: 2, minutes: 0},
			expected: togglTime{hours: 3, minutes: 30},
		},
		{
			name:     "Adding hours and minutes",
			t1:       togglTime{hours: 1, minutes: 30},
			t2:       togglTime{hours: 2, minutes: 45},
			expected: togglTime{hours: 4, minutes: 15},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := togglTime{}
			result.add(tt.t1)
			result.add(tt.t2)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
