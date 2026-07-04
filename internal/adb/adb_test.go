package adb

import (
	"reflect"
	"testing"
)

func TestParseDevices(t *testing.T) {
	out := "List of devices attached\n" +
		"ABC123\tdevice\n" +
		"DEF456\tunauthorized\n" +
		"GHI789\toffline\n" +
		"JKL012\tdevice product:x model:X200\n" +
		"\n"
	got := parseDevices(out)
	want := []string{"ABC123", "JKL012"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseDevices() = %#v, want %#v", got, want)
	}
}

func TestNormalizePackageList(t *testing.T) {
	out := "package:com.vivo.browser\r\n" +
		"package:com.android.notes\n" +
		"\n" +
		"package:com.vivo.browser\n" // duplicate
	got := normalizePackageList(out)
	want := []string{"com.android.notes", "com.vivo.browser"} // sorted, unique
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizePackageList() = %#v, want %#v", got, want)
	}
}

func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		out  string
		want Result
	}{
		{"not installed", "Failure [not installed for 0]", ResultSkipped},
		{"unknown package", "Error: Unknown package: com.x", ResultSkipped},
		{"does not exist", "package com.x does not exist", ResultSkipped},
		{"real failure", "Exception occurred while executing", ResultFailed},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(tc.out); got != tc.want {
				t.Fatalf("Classify(%q) = %v, want %v", tc.out, got, tc.want)
			}
		})
	}
}
