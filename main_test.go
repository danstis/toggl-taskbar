package main

import (
	"bytes"
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

func TestCreateIcon(t *testing.T) {
	tests := []struct {
		name      string
		x         int
		y         int
		hours     int
		threshold int
		expected  []byte
		expectErr bool
	}{
		{
			name:      "Below threshold",
			x:         16,
			y:         16,
			hours:     5,
			threshold: 8,
			expected:  []byte{0, 0, 1, 0, 1, 0, 16, 16, 0, 0, 1, 0, 32, 0, 245, 0, 0, 0, 22, 0, 0, 0, 137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 16, 0, 0, 0, 16, 8, 6, 0, 0, 0, 31, 243, 255, 97, 0, 0, 0, 188, 73, 68, 65, 84, 120, 156, 164, 147, 177, 14, 65, 65, 16, 69, 239, 200, 171, 208, 74, 72, 132, 168, 68, 175, 80, 40, 125, 132, 74, 231, 67, 104, 253, 145, 127, 208, 82, 168, 52, 52, 18, 137, 130, 80, 56, 34, 217, 200, 203, 102, 101, 87, 246, 54, 47, 111, 102, 114, 230, 238, 230, 110, 69, 153, 202, 6, 20, 161, 34, 80, 147, 212, 8, 180, 78, 102, 246, 136, 82, 129, 41, 97, 141, 147, 28, 148, 180, 146, 116, 45, 253, 31, 162, 219, 61, 7, 237, 216, 108, 236, 18, 7, 64, 43, 7, 176, 150, 116, 4, 182, 192, 36, 230, 230, 43, 160, 0, 234, 64, 23, 88, 184, 227, 220, 128, 78, 50, 196, 3, 110, 28, 100, 238, 247, 82, 131, 116, 119, 223, 103, 234, 198, 33, 208, 7, 170, 192, 12, 120, 1, 23, 32, 20, 174, 32, 96, 233, 5, 104, 7, 140, 66, 179, 246, 3, 240, 169, 247, 36, 53, 37, 157, 37, 237, 205, 140, 164, 237, 255, 42, 251, 53, 190, 3, 0, 0, 255, 255, 188, 35, 145, 62, 186, 57, 16, 198, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130},
			expectErr: false,
		},
		{
			name:      "Above threshold",
			x:         16,
			y:         16,
			hours:     10,
			threshold: 8,
			expected:  []byte{0, 0, 1, 0, 1, 0, 16, 16, 0, 0, 1, 0, 32, 0, 201, 1, 0, 0, 22, 0, 0, 0, 137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 16, 0, 0, 0, 16, 8, 2, 0, 0, 0, 144, 145, 104, 54, 0, 0, 1, 144, 73, 68, 65, 84, 120, 156, 98, 153, 199, 64, 26, 96, 34, 81, 61, 170, 6, 54, 1, 1, 49, 75, 75, 20, 105, 22, 22, 25, 79, 79, 133, 144, 16, 54, 1, 1, 136, 8, 11, 3, 3, 3, 151, 164, 164, 110, 121, 185, 164, 131, 131, 160, 174, 238, 219, 243, 231, 55, 153, 152, 64, 229, 184, 184, 60, 246, 236, 17, 49, 53, 253, 247, 235, 215, 239, 47, 95, 118, 56, 57, 189, 191, 122, 21, 100, 3, 151, 140, 12, 183, 140, 204, 203, 35, 71, 24, 152, 80, 44, 212, 171, 172, 20, 181, 180, 220, 237, 237, 189, 74, 94, 158, 137, 149, 213, 118, 193, 2, 168, 147, 222, 156, 62, 189, 47, 36, 228, 222, 138, 21, 104, 206, 85, 75, 74, 250, 241, 242, 229, 211, 93, 187, 126, 188, 121, 243, 120, 203, 22, 97, 19, 19, 65, 93, 93, 156, 158, 230, 146, 150, 230, 148, 146, 250, 116, 247, 46, 3, 3, 3, 143, 156, 156, 132, 189, 61, 3, 3, 131, 136, 137, 9, 78, 13, 28, 162, 162, 12, 12, 12, 191, 63, 127, 230, 16, 21, 181, 93, 176, 224, 193, 218, 181, 12, 12, 12, 156, 226, 226, 4, 130, 149, 141, 143, 207, 126, 201, 146, 227, 57, 57, 223, 158, 60, 97, 96, 96, 248, 255, 239, 31, 78, 13, 223, 158, 62, 101, 96, 96, 16, 50, 48, 56, 91, 83, 243, 225, 218, 53, 72, 176, 126, 127, 254, 28, 167, 134, 31, 175, 95, 127, 186, 125, 251, 239, 247, 239, 111, 206, 156, 97, 96, 96, 16, 208, 210, 98, 96, 96, 120, 121, 244, 40, 72, 3, 43, 15, 143, 144, 129, 1, 159, 138, 10, 3, 3, 3, 51, 39, 167, 144, 129, 1, 191, 186, 58, 3, 3, 195, 141, 105, 211, 216, 132, 132, 212, 146, 146, 132, 244, 244, 100, 60, 61, 159, 238, 220, 249, 249, 222, 61, 134, 121, 12, 12, 219, 29, 28, 254, 163, 130, 119, 23, 46, 204, 99, 96, 152, 207, 204, 124, 119, 233, 82, 136, 200, 251, 203, 151, 87, 72, 73, 205, 99, 96, 96, 36, 152, 248, 120, 149, 149, 217, 248, 249, 223, 93, 184, 240, 255, 223, 63, 104, 210, 192, 15, 62, 131, 163, 2, 14, 72, 78, 173, 128, 0, 0, 0, 255, 255, 51, 169, 145, 163, 21, 109, 87, 240, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130},
			expectErr: false,
		},
		{
			name:      "Equal to threshold",
			x:         16,
			y:         16,
			hours:     8,
			threshold: 8,
			expected:  []byte{0, 0, 1, 0, 1, 0, 16, 16, 0, 0, 1, 0, 32, 0, 112, 1, 0, 0, 22, 0, 0, 0, 137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 16, 0, 0, 0, 16, 8, 2, 0, 0, 0, 144, 145, 104, 54, 0, 0, 1, 55, 73, 68, 65, 84, 120, 156, 98, 153, 199, 64, 26, 96, 34, 81, 61, 3, 11, 22, 49, 70, 70, 49, 11, 11, 62, 85, 213, 111, 207, 159, 191, 216, 191, 255, 223, 159, 63, 120, 53, 48, 50, 58, 175, 95, 47, 231, 239, 255, 243, 237, 91, 118, 97, 225, 119, 231, 207, 111, 179, 179, 251, 253, 229, 11, 78, 39, 73, 216, 218, 202, 249, 251, 95, 237, 235, 91, 38, 34, 114, 169, 189, 93, 200, 208, 80, 214, 199, 7, 159, 31, 216, 133, 132, 224, 236, 47, 15, 31, 50, 48, 48, 252, 250, 244, 9, 159, 147, 94, 28, 62, 252, 243, 205, 27, 173, 188, 188, 47, 15, 30, 168, 38, 37, 189, 57, 125, 250, 217, 174, 93, 40, 78, 198, 12, 86, 197, 176, 48, 135, 149, 43, 65, 172, 127, 255, 182, 90, 91, 191, 58, 113, 2, 159, 147, 132, 13, 13, 173, 103, 207, 126, 186, 115, 231, 237, 185, 115, 25, 152, 152, 220, 119, 239, 22, 53, 55, 71, 86, 192, 236, 143, 170, 193, 110, 225, 66, 1, 45, 173, 61, 62, 62, 183, 231, 207, 255, 251, 253, 187, 140, 183, 55, 175, 146, 210, 157, 69, 139, 112, 218, 192, 175, 169, 201, 192, 192, 240, 245, 201, 19, 6, 6, 134, 171, 125, 125, 12, 255, 255, 115, 73, 74, 226, 115, 210, 135, 171, 87, 25, 24, 24, 52, 210, 211, 153, 57, 56, 148, 99, 98, 24, 24, 25, 159, 238, 220, 137, 207, 211, 124, 42, 42, 206, 235, 215, 11, 232, 232, 128, 56, 255, 255, 223, 95, 177, 226, 72, 114, 242, 159, 239, 223, 241, 133, 18, 3, 35, 163, 128, 134, 6, 155, 128, 192, 231, 123, 247, 190, 191, 124, 137, 46, 73, 243, 212, 10, 8, 0, 0, 255, 255, 83, 176, 111, 51, 96, 121, 253, 143, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createIcon(tt.x, tt.y, tt.hours, tt.threshold)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, but got: %v", tt.expectErr, err)
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Expected icon bytes: %v, but got: %v", tt.expected, result)
			}
		})
	}
}
